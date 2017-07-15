package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"

	"log"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"github.com/server/api"
	"gopkg.in/yaml.v2"
)

const (
	// privKeyPath is the path to the private key
	privKeyPath = "./keys/app.rsa"
	// pubkeyPath is the path to the public key
	pubKeyPath = "./keys/app.rsa.pub"
)

// DBConfig configuration struct
type DBConfig struct {
	user     string `yaml:"user"`
	password string `yaml:"password"`
	dbname   string `yaml:"dbname"`
	host     string `yaml:"host"`
	port     string `yaml:"port"`
	sslmode  string `yaml:"sslmode"`
}

func main() {
	var dbConfig DBConfig
	// Initialize keys
	initKeys()
	filePtr := flag.String("f", "db.config.dev.yaml",
		"Path to the config file. Default: db.config.dev.yaml")

	dbConfig.getConfig(*filePtr)
	// initialize database
	db, err := sql.Open("postgres",
		"user="+dbConfig.user+
			" dbname="+dbConfig.dbname+
			" port="+dbConfig.port+
			" sslmode="+dbConfig.sslmode+
			" host="+dbConfig.host+
			" password="+dbConfig.password)

	if err != nil {
		fmt.Println("Failed to connect to the database")
	}

	// Ping test the database
	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to ping the database")
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
		fmt.Println("failed to read path of key")
	}

	// SignKey parses the private file
	api.SignKey, err = jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		fmt.Println("failed to parse private key")
	}
	// pubBytes is the public RSA file
	pubBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		fmt.Println("failed to read pub key")
	}

	// VerifyKey parses the public file
	api.VerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		fmt.Println("failed to parse key")
	}
}

func (c *DBConfig) getConfig(file string) *DBConfig {
	// Read from a yaml file
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("Failed to read db configuration file")
	}

	// Unmarshal yaml file
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatal("Failed to Unmarshal DB configuration file")
	}
	return c
}
