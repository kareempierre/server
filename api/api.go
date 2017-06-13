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

// DB refers to the database
var (
	DB        *sql.DB
	VerifyKey []byte
	SignKey   []byte
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

	// Spin up api
	http.ListenAndServe(":12000", v2)
}

// authHandler is used to authenticate the user logging in
func authHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

}

// users returns all users for a specific organization
func usersHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Endpoint Hit: user")
}

// user requests information based on the user being viewed
func viewUserHandler(res http.ResponseWriter, req *http.Request) {

}

func authMiddleware(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	token, err := request.ParseFromRequest(req, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return VerifyKey, nil
		})
	if err != nil {
		log.Fatal("Error: Unauthorized access")
		return
	}
	if token.Valid {
		next(res, req)
	}
	if !token.Valid {
		res.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(res, "Token is not Valid")
	}

}
