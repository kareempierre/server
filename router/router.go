package router

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Article discribes the content of a message.
type Article struct {
	ID      int    `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

// Articles type is an array of Articles.
type Articles []Article

var (
	db *sql.DB
)

// HandleRequests deals with all incoming requests
func HandleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/v1/auth", auth)

	myRouter.HandleFunc("/v1/users", users)
	myRouter.HandleFunc("/v1/users/{user}", user)

	myRouter.HandleFunc("/all", returnAllArticles)
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	log.Fatal(http.ListenAndServe(":12000", myRouter))
}

// auth is used to authenticate the user logging in
func auth(res http.ResponseWriter, req *http.Request) {

}

// users returns all users for a specific organization
func users(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Welcome to the HomePage")
	fmt.Println("Endpoint Hit: homePage")
}

// user requests information based on the user being viewed
func user(res http.ResponseWriter, req *http.Request) {

}

func returnAllArticles(res http.ResponseWriter, req *http.Request) {
	articles := Articles{
		Article{Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		Article{Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	}

	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(res).Encode(articles)
}

func returnSingleArticle(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	key := vars["id"]

	fmt.Fprintf(res, "Key: "+key)
}
