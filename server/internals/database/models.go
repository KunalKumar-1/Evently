package database

import "database/sql"

type Models struct {
	Users     UserModel
	Events    EventModel
	Attendees AttendeeModel
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Users:     UserModel{Db: db},
		Events:    EventModel{Db: db},
		Attendees: AttendeeModel{Db: db},
	}
}
