package database

import "database/sql"

type UserModel struct {
	Db *sql.DB
}

type user struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}
