package main

import (
	"database/sql"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kunalkumar-1/Evently/internals/database"
	"github.com/kunalkumar-1/Evently/internals/env"
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func main() {
	db, err := sql.Open("sqlite3", "./data.db") //connects to db
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close() // close the database connection

	models := database.NewModels(db)

	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "secrk23set"),
		models:    *models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}
