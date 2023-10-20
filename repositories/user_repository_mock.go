package repositories

import (
	"database/sql"
	"errors"
	"time"

	"weather-app/database"
	"weather-app/models"
)

type UserRepository interface {
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(user *models.User) error
	GetUserByUsernameForLogin(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

type DBUserRepository struct{}

func (repo *DBUserRepository) CreateUser(user *models.User) error {
	conn := database.GetDB()

	currentTime := time.Now()

	user.CreatedAt = currentTime
	user.UpdatedAt = currentTime

	_, err := conn.Exec("INSERT INTO users(username, email, password, dob, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		user.Username, user.Email, user.Password, user.Dob, user.CreatedAt, user.UpdatedAt)
	return err
}

func (repo *DBUserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	conn := database.GetDB()
	if conn == nil {
		return nil, errors.New("database connection is not initialized")
	}
	err := conn.QueryRow("SELECT id, username, email, password, dob FROM users WHERE email = ?", email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Dob)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *DBUserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	conn := database.GetDB()
	if conn == nil {
		return nil, errors.New("database connection is not initialized")
	}
	err := conn.QueryRow("SELECT id, username, email, password, dob FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Dob)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *DBUserRepository) GetUserByUsernameForLogin(username string) (*models.User, error) {
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
