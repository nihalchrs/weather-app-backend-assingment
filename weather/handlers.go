package weather

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"weather-app/helper"
	"weather-app/models"
	repositories "weather-app/repositories"
)

func getWeatherByCity(c *gin.Context) {
	city := c.Param("city")
	userID, _ := c.Get("userID")

	weatherData, err := helper.GetWeatherDataForUser(city, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "data": weatherData})
}

func getWeatherHistory(c *gin.Context) {
	userID, _ := c.Get("userID")

	weatherHistory, err := repositories.GetWeatherHistoryByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "data": weatherHistory})
}

func updateWeatherData(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": "UserID is required"})
		return
	}
	var weather models.Weather
	weather.UserID = userID.(int)
	if err := c.BindJSON(&weather); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": "Invalid request data", "details": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": "Invalid ID format"})
		return
	}

	belongsToUser, err := repositories.CheckWeatherBelongsToUser(userID.(int), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}

	if !belongsToUser {
		c.JSON(http.StatusForbidden, gin.H{"status": false, "error": "Permission denied"})
		return
	}

	err = repositories.UpdateWeatherData(userID.(int), id, weather.City)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Updated successfully"})
}

func deleteWeatherData(c *gin.Context) {
	userID, _ := c.Get("userID")
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": "Invalid ID format"})
		return
	}

	belongsToUser, err := repositories.CheckWeatherBelongsToUser(userID.(int), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}

	if !belongsToUser {
		c.JSON(http.StatusForbidden, gin.H{"status": false, "error": "You are not allowed to delete this weather record"})
		return
	}

	err = repositories.DeleteWeatherData(userID.(int), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Deleted successfully"})
}

func bulkDeleteWeatherData(c *gin.Context) {
	var ids []int
	if err := c.BindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": "Invalid request data"})
		return
	}

	err := repositories.BulkDeleteWeatherData(ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "message": fmt.Sprintf("Deleted %d records successfully", len(ids))})
}
