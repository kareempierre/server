package router

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
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

var middleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return key, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

// DB refers to the database
var (
	DB  *sql.DB
	key = []byte("secret")
)

// HandleRequests deals with all incoming requests
func HandleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	v2 := myRouter.PathPrefix("/v1").Subrouter()

	// authentication routes
	v2.HandleFunc("/auth", auth).Methods("GET")

	// user routes
	v2.HandleFunc("/users", users).Methods("GET")
	v2.HandleFunc("/users/{user}", user).Methods("GET")

	myRouter.HandleFunc("/all", returnAllArticles)
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	http.ListenAndServe(":12000", v2)
}

// auth is used to authenticate the user logging in
func auth(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["admin"] = true
	claims["name"] = "Ado Kukic"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, _ := token.SignedString(key)

	fmt.Println(tokenString)

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
