package main

import (
	"github.com/gin-gonic/gin"

	"weather-app/config"
	"weather-app/database"
)

func main() {
	config.LoadEnv()
	r := gin.Default()
	config.InitializeRoutes(r)
	config.InitializeDatabase()
	defer database.GetDB().Close()
	r.Run(":8080")
}
