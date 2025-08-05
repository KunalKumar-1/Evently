package main

import (
	"database/sql"
	"log"
)

func main() {
	db, err := sql.Open("sqlite3", "./data.db") //connects to db
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close() // close the database connection
}