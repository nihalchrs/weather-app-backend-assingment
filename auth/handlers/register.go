package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"weather-app/helper"
	"weather-app/models"
	repositories "weather-app/repositories"
)

func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {

		errMessages, isValidationError := helper.HandleValidationError(err)
		if !isValidationError {
			fmt.Println("Raw error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid request payload."})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": errMessages})
		return
	}

	if helper.EmailExists(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Email already exists."})
		return
	}

	if helper.UsernameExists(user.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Username already exists."})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Error hashing password."})
		return
	}
	user.Password = string(hashedPassword)

	if err := repositories.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Error registering the user."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "User registered successfully."})
}
