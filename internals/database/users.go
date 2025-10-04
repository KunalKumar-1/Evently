package database

import (
	"context"
	"database/sql"
	"time"
)

type UserModel struct {
	Db *sql.DB
}

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (u *UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `INSERT INTO users(name, email, password) VALUES($1, $2, $3) RETURNING id`

	return u.Db.QueryRowContext(ctx, query, user.Name, user.Email, user.Password).Scan(&user.Id)

}
