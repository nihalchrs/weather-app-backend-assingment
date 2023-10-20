package weather

import (
	"github.com/gin-gonic/gin"

	"weather-app/middleware"
)

func RunWeatherServer(r *gin.Engine) {
	r.Use(middleware.ConfigureCORS())
	protected := r.Group("/")
	protected.Use(middleware.TokenAuthMiddleware())
	defineRoutes(protected)
}

func defineRoutes(r *gin.RouterGroup) {
	r.GET("/weather/:city", getWeatherByCity)
	r.GET("/weather-history", getWeatherHistory)
	r.PUT("/weather/:id", updateWeatherData)
	r.DELETE("/weather/:id", deleteWeatherData)
	r.DELETE("/weather/bulk", bulkDeleteWeatherData)
}
