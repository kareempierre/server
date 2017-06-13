// will add authboss to this file

package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// DB refers to the database
var (
	DB *sql.DB
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
	next(res, req)
}
