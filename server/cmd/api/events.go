package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kunalkumar-1/Evently/internals/database"
)

// create event
func (app *application) createEvent(c *gin.Context) {

	var event database.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := app.models.Events.Insert(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create event",
		})
		return
	}

	c.JSON(http.StatusCreated, event)
}

func (app *application) getAllEvent(c *gin.Context) {
	events, err := app.models.Events.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retireve events",
		})
		return
	}
	c.JSON(http.StatusOK, events)
}

// get events
func (app *application) getEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event Id",
		})
		return
	}

	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Event not found",
		})
		return
	}

	if event == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retireve event",
		})
		return
	}
	c.JSON(http.StatusCreated, event)
}

// update event
func (app *application) updateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event Id",
		})
		return
	}

	existingEvent, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retireve event",
		})
		return
	}

	if existingEvent == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erorr": "Event not found",
		})
		return
	}

	updatedEvent := &database.Event{}

	if err := c.ShouldBind(updatedEvent); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
	}

	updatedEvent.Id = id

	if err := app.models.Events.Update(updatedEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update event",
		})
		return
	}
}

func (app *application) deleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event Id",
		})
		return
	}

	if err := app.models.Events.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete event",
		})
	}

	c.JSON(http.StatusNoContent, nil)
}
