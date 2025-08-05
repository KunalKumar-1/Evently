package database

import (
	"database/sql"
	"time"
)

type EventModel struct {
	Db *sql.DB
}

type event struct {
	Id          int       `json:"id"`
	OwnerId     int       `json:"ownerId" binding:"required"`
	Name        string    `json:"name" binding:"required, min=3, max=50"`
	Description string    `json:"description" binding:"required, min=3, max=500"`
	Date        time.Time `json:"date" binding:"required" datetime:"2003-02-12"`
	Location    string    `json:"location" binding:"required, min=3"`
}
