// will add authboss to this file

package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

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
	Password string `json:"password"`
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

// authHandler is used to authenticate the user logging in
func authHandler(res http.ResponseWriter, req *http.Request) {
	var user, loggedInUser UserCredentials
	//var email, password string
	// Set Content header type on response
	//res.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		res.WriteHeader(http.StatusForbidden)
		fmt.Fprint(res, "Error in request")
	}

	// DB query needs to be made here
	rowErr := DB.QueryRow(`SELECT email, password  FROM users WHERE email=$1`, user.Email).Scan(&loggedInUser.Email, &loggedInUser.Password)
	if rowErr != nil || rowErr == sql.ErrNoRows {
		log.Fatal("Error: Fetching Rows from db")
		return
	}

	fmt.Println(loggedInUser)
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
