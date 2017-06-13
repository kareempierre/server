package main

import (
	"database/sql"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
	"github.com/server/api"
)

const (
	privKeyPath = "keys/app.rsa"
	pubKeyPath  = "keys/app.rsa.pub"
)

var (
	// VerifyKey is the path to the public key
	VerifyKey []byte
	// SignKey is the path to the private key
	SignKey []byte
)

func main() {
	// Initialize keys
	initKeys()

	db, err := sql.Open("postgres", "dbname=Bishop port=27108 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	api.DB = db

	api.Api()
}

func initKeys() {
	var err error

	// SignKey is the private key
	SignKey, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}

	// VerifyKey is the public key
	VerifyKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal("Error reading public key")
		return
	}

}
