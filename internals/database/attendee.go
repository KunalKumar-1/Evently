package database

import "database/sql"

type AttendeeModel struct {
	Db *sql.DB
}

type Attendee struct {
	Id      int `json:"id"`
	EventId int `json:"eventid"`
	UserId  int `json:"userid"`
}
