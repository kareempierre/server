// will add authboss to this file

package api

import (
	"crypto/rsa"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	// the use for the database
	_ "github.com/lib/pq"
)

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
	v2.HandleFunc("/v2/auth", AuthHandler).Methods("POST")
	v2.HandleFunc("/v2/auth/user/create", CreateUser).Methods("POST")

	// Protected endpoints
	protectedSubRoute := mux.NewRouter().PathPrefix("/v2").Subrouter().StrictSlash(true)
	// Users
	protectedSubRoute.HandleFunc("/users", UsersHandler).Methods("GET")
	protectedSubRoute.HandleFunc("/users/{user}", ViewUserHandler).Methods("GET")
	// Blog
	protectedSubRoute.HandleFunc("/blog", BlogHandler).Methods("GET")
	protectedSubRoute.HandleFunc("/blog/{blog}", ViewBlogHandler).Methods("GET")
	// Gallery
	protectedSubRoute.HandleFunc("/gallery", ViewGalleryHandler).Methods("GET")

	v2.PathPrefix("/v2").Handler(negroni.New(
		negroni.HandlerFunc(AuthMiddleware),
		negroni.Wrap(protectedSubRoute),
	))

	// Protected gallery, threads, posts for admin

	// Spin up api
	n := negroni.New()
	n.UseHandler(v2)
	fmt.Println("Server is running")
	http.ListenAndServe(":12000", n)
}
