package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type EventModel struct {
	Db *sql.DB
}

type Event struct {
	Id          int    `json:"id"`
	OwnerId     int    `json:"ownerId" binding:"required"`
	Name        string `json:"name" binding:"required,min=3,max=50"`
	Description string `json:"description" binding:"required,min=3,max=500"`
	Date        string `json:"date" binding:"required"`
	Location    string `json:"location" binding:"required,min=3"`
}

func (e *EventModel) Insert(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	parsedDate, err := time.Parse("2006-01-02", event.Date)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	query := `INSERT INTO events(owner_id, name, description, date, location) VALUES($1, $2, $3, $4, $5) RETURNING id`

	return e.Db.QueryRowContext(ctx, query, event.OwnerId, event.Name, event.Description, parsedDate, event.Location).Scan(&event.Id)
}

func (e *EventModel) GetAll() ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `SELECT * FROM events`

	rows, err := e.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []*Event{}

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (e *EventModel) Get(id int) (*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `SELECT * FROM events WHERE id = $1`
	var event Event

	err := e.Db.QueryRowContext(ctx, query, id).Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (e *EventModel) Update(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	parsedDate, err := time.Parse("2006-01-02", event.Date)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	query := `UPDATE events SET name = $2, description = $3, date = $4, location = $5 WHERE id = $1`
	_, err = e.Db.ExecContext(ctx, query, event.Id, event.Name, event.Description, parsedDate, event.Location)
	if err != nil {
		return err
	}

	return nil
}

func (e *EventModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `DELETE FROM events WHERE id = $1`
	_, err := e.Db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
