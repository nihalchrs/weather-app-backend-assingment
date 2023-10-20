package config

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	auth "weather-app/auth/handlers"
	"weather-app/database"
	"weather-app/middleware"
	"weather-app/weather"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetDBConfig() string {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbPassword == "" {
		return fmt.Sprintf("%s@tcp(%s:%s)/%s", dbUser, dbHost, dbPort, dbName)
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
}

func InitializeRoutes(r *gin.Engine) {
	r.Use(middleware.ConfigureCORS())
	r.POST("/register", auth.RegisterUser)
	r.POST("/login", auth.LoginUser)
	weather.RunWeatherServer(r)
}

func InitializeDatabase() {
	dataSourceName := GetDBConfig()
	database.InitDB(dataSourceName)
}
