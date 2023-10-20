package auth

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"weather-app/helper"
)

func LoginUser(c *gin.Context) {
	var loginUser struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	user, err := helper.ValidateUserCredentials(loginUser.Username, loginUser.Password)
	if err != nil {
		switch err.Error() {
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "User not found."})
			return
		case "invalid password":
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Invalid login credentials."})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Internal server error."})
			return
		}
	}

	tokenString, err := helper.GenerateJWTToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Error generating token."})
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &helper.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return helper.GetJWTKey(), nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Error parsing token."})
		return
	}

	if claims, ok := token.Claims.(*helper.MyCustomClaims); ok && token.Valid {
		c.JSON(http.StatusOK, gin.H{
			"status":     true,
			"message":    "Login successful.",
			"token":      tokenString,
			"expires_in": claims.ExpiresAt,
			"user":       user,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Error verifying token claims."})
	}
}
