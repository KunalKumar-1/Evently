package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/events", app.getAllEvent)                         // get all events
		v1.GET("/events/:id", app.getEvent)                        // get event by id
		v1.GET("/events/:id/attendees/", app.getAttendeesForEvent) // add attendee to event
		v1.GET("/attendees/:id/events", app.getEventsByAttendee)   // get attendees for event
		v1.POST("/auth/register", app.registerUser)                // register user
		v1.POST("/auth/login", app.login)
	}

	auth := v1.Group("/")
	auth.Use(app.AuthMiddleware())
	{
		auth.POST("/events", app.createEvent)                                     // create event
		auth.PUT("/events/:id", app.updateEvent)                                  // update event by id
		auth.DELETE("/events/:id", app.deleteEvent)                               // delete the event by id
		auth.POST("/events/:id/attendees/:userId", app.addAttendeeToEvent)        // get attendee to event
		auth.DELETE("/events/:id/attendees/:userId", app.deleteAttendeeFromEvent) // delete attendee from event
	}

	return r
}
