package main

import (
	"database/sql"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
	"github.com/server/api"
)

const (
	// privKeyPath is the path to the private key
	privKeyPath = "keys/app.rsa"
	// pubkeyPath is the path to the public key
	pubKeyPath = "keys/app.rsa.pub"
)

func main() {

	// Initialize keys
	initKeys()

	// initialize database
	db, err := sql.Open("postgres", "dbname=Bishop port=27018 sslmode=disable")

	// Check for error on database initialization
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Error: Could not establish a connection with the database")
	}

	// close database when server stops
	defer db.Close()

	// Send database to api
	api.DB = db

	// Initialize API
	api.API()
}

func initKeys() {
	var err error

	// SignKey is the private key
	api.SignKey, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}

	// VerifyKey is the public key
	api.VerifyKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal("Error reading public key")
		return
	}

	//return SignKey, VerifyKey

}
