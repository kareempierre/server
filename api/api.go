// will add authboss to this file

package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"encoding/json"

	"github.com/codegangsta/negroni"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/mux"

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

var (
	// DB is the postregres database
	DB *sql.DB
	// VerifyKey is the public key
	VerifyKey *rsa.PublicKey
	// SignKey is the private key
	SignKey *rsa.PrivateKey
	// SignedKey is the full key with signing
)

// API deals with all incoming requests
func API() {

	// Route version
	v2 := mux.NewRouter().StrictSlash(true)

	// Public endpoints
	v2.HandleFunc("/auth", authHandler).Methods("POST")
	v2.HandleFunc("/auth/user/create", createUser).Methods("POST")

	// Protected endpoints
	protectedUserBaseRoute := mux.NewRouter().PathPrefix("/v2").Subrouter().StrictSlash(true)
	protectedUserBaseRoute.HandleFunc("/", usersHandler).Methods("GET")
	protectedUserBaseRoute.HandleFunc("/{user}", viewUserHandler).Methods("GET")

	v2.PathPrefix("/v2").Handler(negroni.New(
		negroni.HandlerFunc(authMiddleware),
		negroni.Wrap(protectedUserBaseRoute),
	))

	// Protected user routes
	// protectedUserRoute := protectedUserBaseRoute.PathPrefix("/users").Subrouter()

	// Protected gallery, threads, posts for admin

	// Spin up api
	n := negroni.New()
	n.UseHandler(v2)
	fmt.Println("Server is running")
	http.ListenAndServe(":12000", n)
}

func createUser(res http.ResponseWriter, req *http.Request) {
	var user UserConstruct

	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "Some fields were not filled correctly")
	}

	//TODO: WILL NEED TO CHECK TO MAKE SURE THAT THE EMAIL DOES NOT ALREADY EXIST

	//Hash the password and check to see if the email address is already in use
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error: failed to convert the password")
	}
	stmt, err := DB.Prepare(`INSERT INTO users(firstname, lastname, email, admin, creator, organization, password, register)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8);`)
	defer stmt.Close()
	if err != nil {
		log.Fatal("Error: Error on prepare")
	}

	_, ExecErr := stmt.Exec(user.FirstName, user.LastName, user.Email, false, false, user.Organization, string(hashedPassword), time.Now())
	if ExecErr != nil {
		log.Fatal("Error: on exec")
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
		log.Fatal("Error: Failed to sign key")
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
		log.Fatal("Error: error on creating json string")
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(userInfo)

}

// authHandler is used to authenticate the user logging in
func authHandler(res http.ResponseWriter, req *http.Request) {
	var user UserCredentials
	var dbUser UserConstruct

	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		res.WriteHeader(http.StatusForbidden)
		fmt.Fprint(res, "Error in request")
	}
	// DB query needs to be made here
	row := DB.QueryRow(`SELECT firstname, lastname, email, organization, creator, admin, password FROM users WHERE email=$1`, user.Email)
	rowErr := row.Scan(&dbUser.FirstName, &dbUser.LastName, &dbUser.Email, &dbUser.Organization,
		&dbUser.Creator, &dbUser.Admin, &dbUser.Password)

	if rowErr != nil || rowErr == sql.ErrNoRows {
		log.Fatal("Error: Fetching Rows from db ")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		res.WriteHeader(http.StatusForbidden)
		fmt.Fprint(res, "Invalid password")
	}

	token := jwt.New(jwt.SigningMethodRS256)

	// set some claims
	// claims := make(jwt.MapClaims)
	// claims["username"] = dbUser.Email
	// claims["password"] = dbUser.Password
	// claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	// token.Claims = claims

	SignedKey, err := token.SignedString(SignKey)

	if err != nil {
		log.Fatal("Error: Failed to sign key")
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
		log.Fatal("Error: error on creating json string")
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(userInfo)

}

// usersHandler returns all users for a specific organization
func usersHandler(res http.ResponseWriter, req *http.Request) {
	example, err := json.Marshal(struct {
		Example string `json:"Example"`
	}{
		Example: "This is just a test: the Endpoint for users worked",
	})
	if err != nil {
		log.Fatal("Error: error on creating json string")
	}
	res.Header().Set("Content-Type", "application/json")
	fmt.Println("Endpoint Hit: user")

	res.Write(example)
}

// viewUserHandler requests information based on the user being viewed in the admin section
func viewUserHandler(res http.ResponseWriter, req *http.Request) {

}

// authMiddleware authenticates the user logging in
func authMiddleware(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	// token is returned if the Authrization in the header matches with the users token
	token, ParseErr := request.ParseFromRequest(req, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return VerifyKey, nil
		})

	if ParseErr != nil {
		log.Fatal(ParseErr.Error())
		return
	}

	// Check to see if the token is valid
	if token.Valid {
		next(res, req)
	}

	if !token.Valid {
		res.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(res, "Token is not Valid")
	}

}
