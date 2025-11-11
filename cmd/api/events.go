package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kunalkumar-1/Evently/internals/database"
)

// CreateEvent creates a new event
// @Summary      Create a new event
// @Description  Create a new event owned by the authenticated user
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        event  body      database.Event  true  "Event Data"
// @Success      201    {object}  database.Event
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /events [post]
// @Security     BearerAuth
func (app *application) createEvent(c *gin.Context) {

	var event database.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		fmt.Println("Bind error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := app.GetUserFromContext(c)
	event.OwnerId = user.Id

	// Insert into db
	err := app.models.Events.Insert(&event)

	if err != nil {
		fmt.Println("DB insert error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create event",
		})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GetAllEvent returns all events
// @Summary      Get all events
// @Description  Retrieve a list of all available events
// @Tags         Events
// @Accept       json
// @Produce      json
// @Success      200  {array}  database.Event
// @Failure      500  {object}  map[string]string  "Failed to retrieve events"
// @Router       /events [get]

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

// GetEvent retrieves an event by ID
// @Summary      Get event by ID
// @Description  Retrieve details of a specific event
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Success      200  {object}  database.Event
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /events/{id} [get]
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

// UpdateEvent updates an existing event
// @Summary      Update event
// @Description  Update an existing event by ID (only by owner)
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id     path      int             true  "Event ID"
// @Param        event  body      database.Event  true  "Updated Event Data"
// @Success      200    {object}  database.Event
// @Failure      400    {object}  map[string]string
// @Failure      403    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /events/{id} [put]
// @Security     BearerAuth
func (app *application) updateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id")) // get event id from url
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event Id",
		})
		return
	}

	user := app.GetUserFromContext(c) // get user from context
	existingEvent, err := app.models.Events.Get(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retireve event " + err.Error(),
		})
		return
	}

	if existingEvent == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erorr": "Event not found",
		})
		return
	}

	if existingEvent.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{
			"erorr": "You are not authorized to update this event",
		})
		return
	}

	updatedEvent := &database.Event{}

	fmt.Println("Existing Event:", existingEvent)
	fmt.Println("UpadtedEvent:", updatedEvent)

	if err := c.ShouldBindJSON(updatedEvent); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	updatedEvent.Id = id

	if err := app.models.Events.Update(updatedEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update event",
		})
		return
	}

	// Return the updated event
	c.JSON(http.StatusOK, updatedEvent)
}

// DeleteEvent deletes an event by ID
// @Summary      Delete event
// @Description  Delete an event by ID (only by owner)
// @Tags         Events
// @Param        id   path      int  true  "Event ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /events/{id} [delete]
// @Security     BearerAuth
func (app *application) deleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event Id",
		})
		return
	}

	user := app.GetUserFromContext(c) // get user from context
	existingEvent, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retireve event",
		})
	}
	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"erorr": "Event not found",
		})
		return
	}

	if existingEvent.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{
			"erorr": "You are not authorized to delete this event",
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

// AddAttendeeToEvent adds a user to an event
// @Summary      Add attendee to event
// @Description  Adds a user as an attendee to a specific event
// @Tags         Attendees
// @Param        id      path      int  true  "Event ID"
// @Param        userId  path      int  true  "User ID"
// @Success      201     {object}  database.Attendee
// @Failure      400     {object}  map[string]string
// @Failure      403     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Router       /events/{id}/attendees/{userId} [post]
// @Security     BearerAuth
func (app *application) addAttendeeToEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event Id",
		})
		return
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user Id",
		})
		return
	}

	event, err := app.models.Events.Get(eventId) //get event by id
	if err != nil {                              // if error in getting event
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve event",
		})
		return
	}
	if event == nil { // if event not found
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Event not found",
		})
		return
	}

	userToAdd, err := app.models.Users.Get(userId) //get user by id
	if err != nil {                                // if error in getting user
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user",
		})
		return
	}
	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	user := app.GetUserFromContext(c)

	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{
			"erorr": "You are not authorized to add attendees to this event",
		})
		return
	}

	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(event.Id, userToAdd.Id) //get user by id
	if err != nil {                                                                             // if error in getting user
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve attendee",
		})
		return
	}
	if existingAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Attendee already Exists",
		})
		return
	}

	attendee := &database.Attendee{
		EventId: eventId,
		UserId:  userToAdd.Id,
	}

	_, err = app.models.Attendees.Insert(attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to add attendee",
		})
		return
	}

	c.JSON(http.StatusCreated, attendee)
}

// GetAttendeesForEvent retrieves attendees for a given event
// @Summary      Get attendees for event
// @Description  Get a list of all attendees for a specific event
// @Tags         Attendees
// @Param        id   path      int  true  "Event ID"
// @Success      200  {array}  database.Attendee
// @Failure      400  {object}  map[string]string
// @Router       /events/{id}/attendees [get]
func (app *application) getAttendeesForEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event Id",
		})
		return
	}

	events, err := app.models.Attendees.GetAttendeeByEvent(id) //get attendees for event
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve Attendees for events",
		})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (app *application) deleteAttendeeFromEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id")) //event id
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event Id",
		})
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid userId Id",
		})
	}

	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retireve event",
		})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"erorr": "Event not found",
		})
		return
	}

	user := app.GetUserFromContext(c)
	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to delete attendees from this event",
		})
		return
	}

	err = app.models.Attendees.Delete(userId, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete attendee",
		})
	}

	c.JSON(http.StatusNoContent, nil)
}

func (app *application) getEventsByAttendee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid attendee Id",
		})
		return
	}
	user, err := app.models.Users.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get events",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
