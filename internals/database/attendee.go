package database

import (
	"context"
	"database/sql"
	"time"
)

type AttendeeModel struct {
	Db *sql.DB
}

type Attendee struct {
	Id      int `json:"id"`
	EventId int `json:"eventid"`
	UserId  int `json:"userid"`
}

// Insert a new attendee into the database
func (a *AttendeeModel) Insert(attendee *Attendee) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO attendees (event_id, user_id) VALUES($1, $2) RETURNING id`
	err := a.Db.QueryRowContext(ctx, query, attendee.EventId, attendee.UserId).Scan(&attendee.Id)

	if err != nil {
		return nil, err
	}
	return attendee, nil
}

// GetByEventAndAttendee retrieves an attendee by event ID and user ID
func (a *AttendeeModel) GetByEventAndAttendee(eventId, userId int) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, event_id, user_id FROM attendees WHERE event_id = $1 AND user_id = $2`

	var attendee Attendee
	err := a.Db.QueryRowContext(ctx, query, eventId, userId).Scan(&attendee.Id, &attendee.EventId, &attendee.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}
	return &attendee, nil
}

func (a *AttendeeModel) GetAttendeeByEvent(eventId int) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT u.id, u.name, u.email
	FROM users u
	JOIN attendees a ON u.id = a.user_id
	WHERE a.event_id = $1
	`
	rows, err := a.Db.QueryContext(ctx, query, eventId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*User

	for rows.Next() { // looping through all the rows
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}

		users = append(users, &user) //	appending the user to the users slice
	}

	return users, nil
}

func (a *AttendeeModel) Delete(userId, eventId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM attendees WHERE user_id = $1 AND event_id = $2`
	_, err := a.Db.ExecContext(ctx, query, userId, eventId)
	if err != nil {
		return err
	}
	return nil
}

func (a *AttendeeModel) GetEventsByAttendee(AttendeeId int) ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			SELECT e.id, e.owner_id, e.name, e.description, e.date, e.location
			FROM events e
			JOIN attendees a ON e.id = a.event_id
			WHERE a.user_id = $1
		`
	row, err := a.Db.QueryContext(ctx, query, AttendeeId)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	var events []*Event

	for row.Next() {
		var event Event
		err := row.Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}
