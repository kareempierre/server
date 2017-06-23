package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"encoding/json"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"

	"crypto/rsa"
	// the use for the database
	_ "github.com/lib/pq"
)

// UserCredentials requires an Email and Password
type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserConstruct is what you get back from the database
type UserConstruct struct {
	FirstName    string
	LastName     string
	Email        string
	Password     string
	Organization string
	Admin        bool
	Creator      bool
	register     time.Time
}

// ErrorResponse is the type and message of an error
type ErrorResponse struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
}

var (
	// DB is the postregres database
	DB *sql.DB
	// VerifyKey is the public key
	VerifyKey *rsa.PublicKey
	// SignKey is the private key
	SignKey *rsa.PrivateKey
	// SignedKey is the full key with signing
)

// CreateUser creates a new user
func CreateUser(res http.ResponseWriter, req *http.Request) {
	var user UserConstruct

	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		errMessage, _ := json.Marshal(ErrorResponse{
			Type:    http.StatusBadRequest,
			Message: err.Error(),
		})
		res.WriteHeader(http.StatusBadRequest)
		res.Write(errMessage)
		return
	}

	//TODO: WILL NEED TO CHECK TO MAKE SURE THAT THE EMAIL DOES NOT ALREADY EXIST

	//Hash the password and check to see if the email address is already in use
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		errMessage, _ := json.Marshal(ErrorResponse{
			Type:    http.StatusBadRequest,
			Message: err.Error(),
		})
		res.WriteHeader(http.StatusBadRequest)
		res.Write(errMessage)
		return
	}
	stmt, err := DB.Prepare(`INSERT INTO users(firstname, lastname, email, admin, creator, organization, password, register)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8);`)
	defer stmt.Close()
	if err != nil {
		errMessage, _ := json.Marshal(ErrorResponse{
			Type:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(errMessage)
		return
	}

	_, ExecErr := stmt.Exec(user.FirstName, user.LastName, user.Email, false, false, user.Organization, string(hashedPassword), time.Now())
	if ExecErr != nil {
		errMessage, _ := json.Marshal(ErrorResponse{
			Type:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(errMessage)
		return
	}
	token := jwt.New(jwt.SigningMethodRS256)

	// set some claims
	// claims := make(jwt.MapClaims)
	// claims["username"] = user.Email
	// claims["password"] = user.Password
	// claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	// token.Claims = claims
	SignedKey, err := token.SignedString(SignKey)
	if err != nil {
		errMessage, _ := json.Marshal(ErrorResponse{
			Type:    http.StatusUnauthorized,
			Message: err.Error(),
		})
		res.WriteHeader(http.StatusUnauthorized)
		res.Write(errMessage)
		return
	}

	userInfo, err := json.Marshal(struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Token     string `json:"token"`
	}{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Token:     SignedKey,
	})
	if err != nil {
		errMessage, _ := json.Marshal(ErrorResponse{
			Type:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(errMessage)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(userInfo)

}

// AuthHandler is used to authenticate the user logging in
func AuthHandler(res http.ResponseWriter, req *http.Request) {
	var user UserCredentials
	var dbUser UserConstruct

	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		errMessage, _ := json.Marshal(ErrorResponse{
			Type:    http.StatusBadRequest,
			Message: err.Error(),
		})
		res.WriteHeader(http.StatusBadRequest)
		res.Write(errMessage)
		return
	}

	// DB query needs to be made here
	row := DB.QueryRow(`SELECT firstname, lastname, email, organization, creator, admin, password FROM users WHERE email=$1`, user.Email)
	rowErr := row.Scan(&dbUser.FirstName, &dbUser.LastName, &dbUser.Email, &dbUser.Organization,
		&dbUser.Creator, &dbUser.Admin, &dbUser.Password)

	if rowErr != nil || rowErr == sql.ErrNoRows {
		errMessage, _ := json.Marshal(ErrorResponse{
			Type:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(errMessage)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		errMessage, _ := json.Marshal(ErrorResponse{
			Type:    http.StatusUnauthorized,
			Message: err.Error(),
		})
		res.WriteHeader(http.StatusUnauthorized)
		res.Write(errMessage)
		return
	}

	token := jwt.New(jwt.SigningMethodRS256)

	SignedKey, err := token.SignedString(SignKey)
	if err != nil {
		errMessage, _ := json.Marshal(ErrorResponse{
			Type:    http.StatusBadRequest,
			Message: err.Error(),
		})
		res.WriteHeader(http.StatusBadRequest)
		res.Write(errMessage)
		return
	}

	userInfo, err := json.Marshal(struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Token     string `json:"token"`
	}{
		FirstName: dbUser.FirstName,
		LastName:  dbUser.LastName,
		Token:     SignedKey,
	})
	if err != nil {
		errMessage, _ := json.Marshal(ErrorResponse{
			Type:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(errMessage)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(userInfo)

}

// UsersHandler returns all users for a specific organization
func UsersHandler(res http.ResponseWriter, req *http.Request) {
	example, err := json.Marshal(struct {
		Example string `json:"Example"`
	}{
		Example: "This is just a test: the Endpoint for users worked",
	})
	if err != nil {
		errMessage, _ := json.Marshal(ErrorResponse{
			Type:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(errMessage)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	fmt.Println("Endpoint Hit: user")

	res.Write(example)
}

// ViewUserHandler requests information based on the user being viewed in the admin section
func ViewUserHandler(res http.ResponseWriter, req *http.Request) {

}

// AuthMiddleware authenticates the user logging in
func AuthMiddleware(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	// token is returned if the Authrization in the header matches with the users token
	token, err := request.ParseFromRequest(req, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return VerifyKey, nil
		})
	if err != nil {
		errMessage, _ := json.Marshal(ErrorResponse{
			Type:    http.StatusUnauthorized,
			Message: err.Error(),
		})
		res.WriteHeader(http.StatusUnauthorized)
		res.Write(errMessage)
		return
	}

	if token.Valid {
		next(res, req)
	}

	if !token.Valid {
		unauthorizedMessage, _ := json.Marshal(struct {
			Error string `json:"error"`
		}{
			Error: errors.New("Unauthorized Entry").Error(),
		})
		res.WriteHeader(http.StatusUnauthorized)
		res.Write(unauthorizedMessage)
		return
	}

}
