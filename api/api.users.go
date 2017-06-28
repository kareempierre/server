package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
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

// UsersList is a struct for the information returned on all users
type UsersList struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Admin     bool   `json:"admin"`
}

// ErrorResponse is the type and message of an error
type ErrorResponse struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
}

var (
	// DBUser is the current logged in user
	DBUser UserConstruct
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
	//var dbUser UserConstruct

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
	rowErr := row.Scan(&DBUser.FirstName, &DBUser.LastName, &DBUser.Email, &DBUser.Organization,
		&DBUser.Creator, &DBUser.Admin, &DBUser.Password)

	if rowErr != nil || rowErr == sql.ErrNoRows {
		errMessage, _ := json.Marshal(ErrorResponse{
			Type:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(errMessage)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(DBUser.Password), []byte(user.Password))
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
		FirstName: DBUser.FirstName,
		LastName:  DBUser.LastName,
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
	var usersList UsersList
	var usersArray []UsersList

	//var organization string
	rows, err := DB.Query(`SELECT firstname, lastname, email, admin FROM users WHERE organization=$1;`, DBUser.Organization)
	if DBUser.Organization == "all" && DBUser.Creator == true {
		rows, err = DB.Query(`SELECT firstname, lastname, email, admin FROM users;`)
	}

	if err != nil {
		errMessage, _ := json.Marshal(ErrorResponse{
			Type:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(errMessage)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&usersList.FirstName, &usersList.LastName, &usersList.Email, &usersList.Admin)
		if err != nil {
			errMessage, _ := json.Marshal(ErrorResponse{
				Type:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(errMessage)
		}
		usersArray = append(usersArray, usersList)
	}

	if DBUser.Admin == true {
		users, err := json.Marshal(usersArray)
		if err != nil {
			errMessage, _ := json.Marshal(ErrorResponse{
				Type:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(errMessage)
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write(users)
	} else {
		unAuthUser, err := json.Marshal(struct {
			Type  int  `json:"type"`
			Admin bool `json:"admin"`
		}{
			Type:  http.StatusUnauthorized,
			Admin: false,
		})
		if err != nil {
			errMessage, _ := json.Marshal(ErrorResponse{
				Type:    http.StatusUnauthorized,
				Message: err.Error(),
			})
			res.WriteHeader(http.StatusUnauthorized)
			res.Write(errMessage)
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write(unAuthUser)
	}
}

// ViewUserHandler requests information based on the user being viewed in the admin section
func ViewUserHandler(res http.ResponseWriter, req *http.Request) {
	//var email string

	//row := DB.QueryRow(`SELECT firstname, lastname, email`)
}
