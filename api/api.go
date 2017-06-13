// will add authboss to this file

package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/mux"
)

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
	v2.HandleFunc("/auth", authHandler).Methods("GET")

	// Protected endpoints
	protectedUserBaseRoute := mux.NewRouter()
	v2.PathPrefix("/users").Handler(negroni.New(
		negroni.HandlerFunc(authMiddleware),
		negroni.Wrap(protectedUserBaseRoute),
	))

	// Protected user routes
	protectedUserRoute := protectedUserBaseRoute.PathPrefix("users").Subrouter()
	protectedUserRoute.HandleFunc("/users", usersHandler).Methods("GET")
	protectedUserRoute.HandleFunc("/users/{user}", viewUserHandler).Methods("GET")

	// Protected gallery, threads, posts for admin

	// Spin up api
	http.ListenAndServe(":12000", v2)
}

// authHandler is used to authenticate the user logging in
func authHandler(res http.ResponseWriter, req *http.Request) {

	// Set Content header type on response
	res.Header().Set("Content-Type", "application/json")

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
