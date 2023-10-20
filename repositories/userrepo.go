package repositories

import (
	"database/sql"
	"errors"
	"time"

	"weather-app/database"
	"weather-app/models"
)

func CreateUser(user *models.User) error {
	db := database.GetDB()
	currentTime := time.Now()

	user.CreatedAt = currentTime
	user.UpdatedAt = currentTime

	_, err := db.Exec("INSERT INTO users(username, email, password, dob, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		user.Username, user.Email, user.Password, user.Dob, user.CreatedAt, user.UpdatedAt)
	return err
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database is not initialized")
	}
	err := db.QueryRow("SELECT id, username, email, password, dob FROM users WHERE email = ?", email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Dob)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database is not initialized")
	}
	err := db.QueryRow("SELECT id, username, email, password, dob FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Dob)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func GetUserByUsernameForLogin(username string) (*models.User, error) {
	var user models.User
	conn := database.GetDB()
	if conn == nil {
		return nil, errors.New("database connection is not initialized")
	}
	err := conn.QueryRow("SELECT id, username, email, password, dob FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Dob)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
