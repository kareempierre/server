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
	// the use for the database
	_ "github.com/lib/pq"
)

// UserCredentials requires an Email and Password
type UserCredentials struct {
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

// UserConstruct is what you get back from the database
type UserConstruct struct {
	FirstName    string
	LastName     string
	Email        string
	Password     []byte
	Organization string
	Admin        bool
	Creator      bool
}

var (
	// DB is the postregres database
	DB *sql.DB
	// VerifyKey is the public key
	VerifyKey []byte
	// SignKey is the private key
	SignKey []byte
)

// API deals with all incoming requests
func API() {

	// Route version
	myRouter := mux.NewRouter().StrictSlash(true)
	v2 := myRouter.PathPrefix("/v2").Subrouter()

	// Public endpoints
	v2.HandleFunc("/auth", authHandler).Methods("POST")
	v2.HandleFunc("/auth/create/user", createUser).Methods("POST")

	// Protected endpoints
	protectedUserBaseRoute := mux.NewRouter()
	v2.PathPrefix("/users").Handler(negroni.New(
		negroni.HandlerFunc(authMiddleware),
		negroni.Wrap(protectedUserBaseRoute),
	))

	// Protected user routes
	protectedUserRoute := protectedUserBaseRoute.PathPrefix("/users").Subrouter()
	protectedUserRoute.HandleFunc("/users", usersHandler).Methods("GET")
	protectedUserRoute.HandleFunc("/users/{user}", viewUserHandler).Methods("GET")

	// Protected gallery, threads, posts for admin

	// Spin up api
	fmt.Println("Server is running")
	http.ListenAndServe(":12000", myRouter)
}

func createUser(res http.ResponseWriter, req *http.Request) {
	var user UserConstruct

	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "Some fields were not filled correctly")
	}

	// Hash the password and check to see if the email address is already in use
	hashedPassword, err := bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error: failed to convert the password")
	}

	rowErr := DB.QueryRow(
		`INSERT INTO users(firstname, lastname, email, admin, creator, organization, password, register) 
		VALUES($1,$2,$3,$4,$5,$6,$7)`,
		user.FirstName, user.LastName, user.Email, false, false, user.Organization, hashedPassword, time.Now(),
	)
	if rowErr != nil {
		log.Fatal("Error: error creating new user")
	}
	// save to postgres database
	// return a token

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
		log.Fatal("Error: Fetching Rows from db")
		return
	}

	err = bcrypt.CompareHashAndPassword(dbUser.Password, user.Password)
	if err != nil {
		res.WriteHeader(http.StatusForbidden)
		fmt.Fprint(res, "Invalid password")
	}

	//res.Header().Set("Content-Type", "application/json")
	// create token and send it back as json data using json.unmarshal
	// will have to read up on this
	fmt.Println(dbUser)
}

// usersHandler returns all users for a specific organization
func usersHandler(res http.ResponseWriter, req *http.Request) {

	fmt.Println("Endpoint Hit: user")
}

// viewUserHandler requests information based on the user being viewed in the admin section
func viewUserHandler(res http.ResponseWriter, req *http.Request) {

}

// authMiddleware authenticates the user logging in
func authMiddleware(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	// token is returned if the Authrization in the header matches with the users token
	token, err := request.ParseFromRequest(req, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return VerifyKey, nil
		})

	// check to see if there is an error
	if err != nil {
		log.Fatal("Error: Unauthorized access")
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
