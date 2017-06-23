// will add authboss to this file

package api

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	// the use for the database
	_ "github.com/lib/pq"
)

// API deals with all incoming requests
func API() {

	// Route version
	v2 := mux.NewRouter().StrictSlash(true)

	// Public endpoints
	v2.HandleFunc("/v2/auth", AuthHandler).Methods("POST")
	v2.HandleFunc("/v2/auth/user/create", CreateUser).Methods("POST")

	// Protected endpoints
	protectedUserBaseRoute := mux.NewRouter().PathPrefix("/v2").Subrouter().StrictSlash(true)
	protectedUserBaseRoute.HandleFunc("/users", UsersHandler).Methods("GET")
	protectedUserBaseRoute.HandleFunc("/{user}", ViewUserHandler).Methods("GET")

	v2.PathPrefix("/v2").Handler(negroni.New(
		negroni.HandlerFunc(AuthMiddleware),
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
