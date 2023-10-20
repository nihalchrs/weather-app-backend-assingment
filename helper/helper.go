package helper

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/resty.v1"

	"weather-app/models"
	"weather-app/repositories"
)

type MyCustomClaims struct {
	UserId   int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

type WeatherResponse struct {
	Main struct {
		Humidity int     `json:"humidity"`
		Temp     float64 `json:"temp"`
		Pressure int     `json:"pressure"`
	} `json:"main"`
	City string `json:"name"`
}

var UserRepo repositories.UserRepository = &repositories.DBUserRepository{}

func HandleValidationError(err error) ([]string, bool) {
	var errMessages []string
	if validationErrors, ok := err.(*validator.ValidationErrors); ok {
		for _, e := range *validationErrors {
			field := e.Field()
			switch field {
			case "Email":
				errMessages = append(errMessages, fmt.Sprintf("Email is invalid (%s)", e.ActualTag()))
			case "Username":
				errMessages = append(errMessages, fmt.Sprintf("Username is invalid (%s)", e.ActualTag()))
			case "Password":
				if e.ActualTag() == "min" {
					errMessages = append(errMessages, "Password must be at least 6 characters long")
				} else {
					errMessages = append(errMessages, fmt.Sprintf("Password is invalid (%s)", e.ActualTag()))
				}
			case "Dob":
				if e.ActualTag() == "datetime" {
					errMessages = append(errMessages, "Date of Birth must be in the format YYYY-MM-DD")
				} else if e.ActualTag() == "lt" {
					errMessages = append(errMessages, "Date of Birth must be a date in the past")
				} else {
					errMessages = append(errMessages, fmt.Sprintf("Date of Birth is invalid (%s)", e.ActualTag()))
				}
			}
		}
		return errMessages, true
	}
	return nil, false
}

func EmailExists(email string) bool {
	user, err := repositories.GetUserByEmail(email)
	if err != nil && err != sql.ErrNoRows {
		return true
	}
	return user != nil
}

func UsernameExists(username string) bool {
	user, err := UserRepo.GetUserByUsername(username)
	if err != nil && err != sql.ErrNoRows {
		return true
	}
	return user != nil
}

func GetJWTKey() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func ValidateUserCredentials(username, password string) (*models.User, error) {
	user, err := UserRepo.GetUserByUsernameForLogin(username)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}
	return user, nil
}

func GenerateJWTToken(user models.User) (string, error) {
	claims := &MyCustomClaims{
		UserId:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(GetJWTKey())
}

func GetWeatherDataForUser(city string, userID int) (interface{}, error) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	baseURL := os.Getenv("OPENWEATHER_BASE_URL")
	client := resty.New()
	url := baseURL + "?q=" + city + "&APPID=" + apiKey + "&units=metric"

	resp, err := client.R().Get(url)
	fmt.Println(url)
	if err != nil {
		log.Println("Request failed:", err)
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		log.Printf("OpenWeatherMap API returned status code: %d", resp.StatusCode())
		return nil, errors.New("API request returned an error")
	}

	var weatherData WeatherResponse
	if err := json.Unmarshal(resp.Body(), &weatherData); err != nil {
		log.Println("Failed to unmarshal response:", err)
		return nil, err
	}

	if err := repositories.InsertWeatherHistoryForUser(userID, weatherData.City, weatherData.Main.Temp, weatherData.Main.Pressure, weatherData.Main.Humidity); err != nil {
		log.Println("Failed to insert weather data into database:", err)
		return WeatherResponse{}, err
	}

	return weatherData, nil
}
