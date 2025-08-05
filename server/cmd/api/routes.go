package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.POST("/events", app.createEvent)       // create event
		v1.GET("/events", app.getAllEvent)        // get all events
		v1.GET("/events/:id", app.getEvent)       // get event by id
		v1.PUT("/events/:id", app.updateEvent)    // update event by id
		v1.DELETE("/events/:id", app.deleteEvent) //delete the event by id
	}

	return r
}
