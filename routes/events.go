package routes

import (
	"net/http"

	"go-rest-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, events)
}

func getEvent(context *gin.Context) {
	event, err := models.GetEvent(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, event)
}

func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := context.GetString("userId")
	event.UserID, err = primitive.ObjectIDFromHex(data)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := event.Save(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"status": "Event created successfully", "event": event})
}

func updateEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err = models.GetEvent(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if event.UserID.Hex() != context.GetString("userId") {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to update this event"})
		return
	}

	if err := event.Update(context.Param("id")); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "Event updated successfully", "event": event})
}

func deleteEvent(context *gin.Context) {

	event, err := models.GetEvent(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if event.UserID.Hex() != context.GetString("userId") {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to delete this event"})
		return
	}

	err = models.DeleteEvent(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"status": "Event deleted successfully"})
}
