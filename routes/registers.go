package routes

import (
	"go-rest-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerForEvent(context *gin.Context) {
	userId := context.MustGet("userId").(string)
	event, err := models.GetEvent(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	err = event.Register(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"status": "Registered successfully"})
}

func unregisterFromEvent(context *gin.Context) {
	userId := context.MustGet("userId").(string)
	event, err := models.GetEvent(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	err = event.Unregister(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "Unregistered successfully"})
}
