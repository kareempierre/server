// will add authboss to this file

package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// DB refers to the database
var (
	DB *sql.DB
)

// HandleRequests deals with all incoming requests
func HandleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	v2 := myRouter.PathPrefix("/v1").Subrouter()

	// authentication routes
	v2.HandleFunc("/auth", authHandler).Methods("GET")

	// user routes
	v2.HandleFunc("/users", usersHandler).Methods("GET")
	v2.HandleFunc("/users/{user}", viewUserHandler).Methods("GET")
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
