package main

import (
	"database/sql"
	"log"

	_ "github.com/kunalkumar-1/Evently/docs"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kunalkumar-1/Evently/internals/database"
	"github.com/kunalkumar-1/Evently/internals/env"
	_ "github.com/mattn/go-sqlite3"
)

// @title           Evently API
// @version         1.0
// @description     Event Management REST API built with Go and Gin.

// @contact.name   API Support
// @contact.url    https://github.com/kunalkumar-1/Evently
// @contact.email  kunaldevspro@gmail.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1
// @schemes   http

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
