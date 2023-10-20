package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/resty.v1"

	"weather-app/database"
	"weather-app/models"
)

type WeatherResponse struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure float64 `json:"pressure"`
		Humidity float64 `json:"humidity"`
	} `json:"main"`
	City string `json:"name"`
}

func InsertWeatherHistoryForUser(userID int, city string, temp float64, pressure int, humidity int) error {
	db := database.GetDB()
	query := "INSERT INTO weather_history (user_id, city, temp, pressure, humidity) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(query, userID, city, temp, pressure, humidity)
	return err
}

func GetWeatherHistoryByUserID(userID int) ([]models.Weather, error) {
	db := database.GetDB()
	var weathers []models.Weather

	query := "SELECT id, city, temp, pressure, humidity, created_at, updated_at FROM weather_history WHERE user_id = ?"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var weather models.Weather
		var createdAt, updatedAt string
		err = rows.Scan(
			&weather.ID,
			&weather.City,
			&weather.Temp,
			&weather.Pressure,
			&weather.Humidity,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		weather.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		weather.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

		weathers = append(weathers, weather)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return weathers, nil
}
func DeleteWeatherData(userID int, id int) error {
	db := database.GetDB()
	query := "DELETE FROM weather_history WHERE id=? AND user_id=?"

	result, err := db.Exec(query, id, userID)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return errors.New("ID not found for the given user")
	}

	return nil
}

func UpdateWeatherData(userID int, id int, city string) error {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	baseURL := os.Getenv("OPENWEATHER_BASE_URL")
	client := resty.New()
	url := baseURL + "?q=" + city + "&APPID=" + apiKey + "&units=metric"
	resp, err := client.R().Get(url)
	if err != nil {
		log.Println("Request to OpenWeatherMap API failed:", err)
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		log.Printf("OpenWeatherMap API returned status code: %d", resp.StatusCode())
		return errors.New("API request returned an error")
	}

	var weatherData WeatherResponse
	if err := json.Unmarshal(resp.Body(), &weatherData); err != nil {
		log.Println("Failed to unmarshal response from OpenWeatherMap API:", err)
		return err
	}

	db := database.GetDB()
	query := "UPDATE weather_history SET city=?, temp=?, pressure=?, humidity=? WHERE id=? AND user_id=?"

	result, err := db.Exec(query, city, weatherData.Main.Temp, weatherData.Main.Pressure, weatherData.Main.Humidity, id, userID)
	if err != nil {
		log.Println("Error executing SQL update:", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Println("No rows updated.")
		return errors.New("no rows updated")
	}
	return nil
}

func BulkDeleteWeatherData(ids []int) error {
	if len(ids) == 0 {
		return errors.New("no id's provided")
	}

	db := database.GetDB()

	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")

	queryArgs := make([]interface{}, len(ids))
	for i, id := range ids {
		queryArgs[i] = id
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM weather_history WHERE id IN (%s)", placeholders)
	var count int
	if err := db.QueryRow(countQuery, queryArgs...).Scan(&count); err != nil {
		return err
	}

	if count != len(ids) {
		return errors.New("some id's were not found in the database")
	}

	deleteQuery := fmt.Sprintf("DELETE FROM weather_history WHERE id IN (%s)", placeholders)
	_, err := db.Exec(deleteQuery, queryArgs...)
	if err != nil {
		return err
	}

	return nil
}

func CheckWeatherBelongsToUser(userID int, weatherID int) (bool, error) {
	db := database.GetDB()
	var ownerID int
	err := db.QueryRow("SELECT user_id FROM weather_history WHERE id = ?", weatherID).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("weather record not found")
		}
		return false, err
	}

	if userID == ownerID {
		return true, nil
	}

	return false, nil
}
