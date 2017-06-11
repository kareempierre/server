package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/server/router"
)

func main() {

	db, err := sql.Open("postgres", "dbname=Bishop port=27108 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	router.db = db
	router.HandleRequests()
}
