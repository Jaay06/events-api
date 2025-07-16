package main

import (
	"go-rest-api/internal/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateEvent creates a new event
//
// @Summary Creates a new event
// @Description Creates a new event
// @Tags events
// @Accept json
// @Produce json
// @Param event body database.Event true "Event"
// @Success 201	{object} database.Event
// @Router /api/v1/events [post]
func (app *application) createEvent(c *gin.Context) {
	var event database.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := app.models.Events.Insert(&event)

	user := app.GetUserFromContext(c)

	event.OwnerId = user.Id

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create event",
		})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GetEvents returns all events
//
// @Summary Returns all events
// @Description Returns all events
// @Tags events
// @Accept json
// @Produce json
// @Success 200	{object} []database.Event
// @Router /api/v1/events [get]
func (app *application) getAllEvents(c *gin.Context) {
	events, err := app.models.Events.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrive events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetEvent returns single event
//
// @Summary Returns a single event
// @Description Returns a single event
// @Tags events
// @Accept json
// @Produce json
// @Success 200 {object} database.Event
// @Router /api/v1/events/{id} [get]
// @Param id path int true "Event ID"
func (app *application) getEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := app.models.Events.Get(id)

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrive event"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// UpdateEvent returns a single updated event
//
// @Summary Returns the updated event
// @Description Returns the updated event
// @Tags events
// @Accept json
// @Produce json
// @Success 200 {object} database.Event
// @Param id path int true "Event ID"
// @Router /api/v1/events/{id} [put]
func (app *application) updateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	user := app.GetUserFromContext(c)
	existingEvent, err := app.models.Events.Get(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrive event"})
		return
	}

	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if existingEvent.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this event"})
		return
	}

	updateEvent := &database.Event{}

	if err := c.ShouldBindJSON(updateEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateEvent.Id = id

	if err := app.models.Events.Update(updateEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	c.JSON(http.StatusOK, updateEvent)
}

// DeleteEvent deletes an event
//
//	@Summary 		deletes an event
//	@Description 	deletes an event
//	@Tags 			events
//	@Accept 		json
//	@Produces 		json
//	@Success 		204
//	@Param 			id path 	int true 	"Event ID"
//	@Router 		/api/v1/events/{id} 	[delete]
//	@Security		BearerAuth
func (app *application) deleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Event id"})
		return
	}

	user := app.GetUserFromContext(c)

	existingEvent, err := app.models.Events.Get(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive event"})
		return
	}

	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if existingEvent.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this event"})
		return
	}

	if err := app.models.Events.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// CreateAttendeesForEvent creates a new event
//
//	@Summary 		adds an attendee to an event
//	@Description	this endpoint is used to add a new attendee to an event using the "Event ID" in the param and attatches it to the user using the user id
//	@Tags 			attendees
//	@Accept 		json
//	@Produce 		json
//	@Param 			id path		int true	"Event ID"
//	@Param 			userId path	int true	"User ID"
//	@Success 		201	{object} database.Attendee
//	@Router 		/api/v1/events/{id}/attendees/{userId} [post]
//	@Security		BearerAuth
func (app *application) addAttendeeToEvent(c *gin.Context) {

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event Id"})
		return
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	event, err := app.models.Events.Get(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrive event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	userToAdd, err := app.models.Users.Get(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive user"})
		return
	}
	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user := app.GetUserFromContext(c)

	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to add an attendee"})
		return
	}

	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(event.Id, userToAdd.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive attendee"})
		return
	}
	if existingAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Attendee already exists"})
		return
	}

	attendee := database.Attendee{
		EventId: event.Id,
		UserId:  userToAdd.Id,
	}

	_, err = app.models.Attendees.Insert(&attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add attendee"})
		return
	}

	c.JSON(http.StatusCreated, attendee)
}

// GetAllAttendeesForEvent get all attendees for an event
//
//	@Summary 		get all attendees for an event
//	@Description	get all attendees for an event
//	@Tags 			attendees
//	@Accept 		json
//	@Produce 		json
//	@Param 			id path		int true	"Event ID"
//	@Success 		200	{object} []database.Attendee
//	@Router 		/api/v1/events/{id}/attendees/ [get]
//	@Security		BearerAuth
func (app *application) getAttendeesForEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event id"})
		return
	}

	users, err := app.models.Attendees.GetAttendeesByEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendee"})
		return
	}

	c.JSON(http.StatusOK, users)

}

// GetEventsByAttendee get attendee for an event
//
//	@Summary 		returns an attendee from an event
//	@Description	returns an attendee from an event
//	@Tags 			attendees
//	@Accept 		json
//	@Produce 		json
//	@Param 			id path		int true	"Event ID"
//	@Param			userId path	int true	"User ID"
//	@Success 		200 {object} database.Attendee
//	@Router 		/api/v1/events/{id}/attendees/{userId} [get]
//	@Security		BearerAuth
func (app *application) getEventsByAttendee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendee id"})
		return
	}

	events, err := app.models.Attendees.GetEventsByAttendee(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get event"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// DeleteAttendeeForEvent get all attendees for an event
//
//	@Summary 		delete an attendees from an event
//	@Description	delete an attendees from an event
//	@Tags 			attendees
//	@Accept 		json
//	@Produce 		json
//	@Param 			id path		int true	"Event ID"
//	@Param			userId path	int true	"User ID"
//	@Success 		204
//	@Router 		/api/v1/events/{id}/attendees/{userId} [delete]
//	@Security		BearerAuth
func (app *application) deleteAttendeeFromEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event id"})
		return
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user. id"})
		return
	}

	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	user := app.GetUserFromContext(c)
	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete an attendee from an event"})
		return
	}

	err = app.models.Attendees.Delete(userId, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendee"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
