package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
	"github.com/server/api"
)

const (
	privKeyPath = "keys/app.rsa"
	pubKeyPath  = "keys/app.rsa.pub"
)

func main() {
	// Initialize keys
	initKeys()
	fmt.Println(VerifyKey, SignKey)
	db, err := sql.Open("postgres", "dbname=Bishop port=27108 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	api.DB = db

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
