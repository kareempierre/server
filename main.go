package main

import (
	"database/sql"
	"io/ioutil"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
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
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}

	// Ping test the database
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping the database")
	}

	// close database when server stops
	defer db.Close()

	// Send database to api
	api.DB = db

	// Initialize API
	api.API()
}

func initKeys() {

	// privBytes is the private RSA file
	privBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal("Failed on private key")
	}

	// SignKey parses the private file
	api.SignKey, err = jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		log.Fatal("Failed on private key")
	}
	// pubBytes is the public RSA file
	pubBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal("Failed on public key")
	}

	// VerifyKey parses the public file
	api.VerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		log.Fatal("Failed on public key")
	}
}
