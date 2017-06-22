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
	//TODO: WILL NEED TO ADD THE CRYPTO LIBRARY TO HANDLE THE RSA KEYS. NEED TO LEARN
	// HOW TO DO THIS
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

	privBytes, err := ioutil.ReadFile(privKeyPath)
	checkErr(err)

	api.SignKey, err = jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	checkErr(err)

	// block, _ := pem.Decode(privKey)
	// if block == nil {
	// 	log.Fatal("Error: Failed to decode RSA")
	// }
	// if block.Type != "RSA PRIVATE KEY" {
	// 	log.Fatal("Error: Type failed on RSA pem conversion")
	// }
	// // SignKey is the private key
	// api.SignKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	// if err != nil {
	// 	log.Fatal("Error: occured on Private key parse from pem")
	// }

	// VerifyKey is the public key
	pubBytes, err := ioutil.ReadFile(pubKeyPath)
	checkErr(err)

	api.VerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	checkErr(err)

	// pubBlock, _ := pem.Decode(pubKey)
	// if pubBlock == nil {
	// 	log.Fatal("Error: Failed to decode RSA")
	// }

	// pub, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	// if err != nil {
	// 	log.Fatal("Error: occured on Public key parse from pem")
	// }
	// api.VerifyKey = pub.(*rsa.PublicKey)
	//return SignKey, VerifyKey

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
	return
}
