package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	r.GET("/swagger/*any", func(c *gin.Context){
		if c.Request.RequestURI == "/swagger/" {
		   c.Redirect(302, "/swagger/index.html")
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:8080/swagger/doc.json"))(c)
	})

	return r
}
