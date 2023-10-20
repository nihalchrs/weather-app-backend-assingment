package auth_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"

	auth "weather-app/auth/handlers"
	"weather-app/models"
	repositories "weather-app/repositories"
)

type MockedRepository struct {
	mock.Mock
}
type UserRepository interface {
	GetUserByUsername(username string) (*models.User, error)
}

var UserRepo repositories.UserRepository = &repositories.DBUserRepository{}

func (m *MockedRepository) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	user, ok := args.Get(0).(*models.User)
	if !ok {
		return nil, args.Error(1)
	}

	return user, args.Error(1)
}

func (m *MockedRepository) GetUserByUsernameForLogin(username string) (*models.User, error) {
	args := m.Called(username)
	user, ok := args.Get(0).(*models.User)
	if !ok {
		return nil, args.Error(1)
	}

	return user, args.Error(1)
}
func (m *MockedRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockedRepository) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	user, ok := args.Get(0).(*models.User)
	if !ok {
		return nil, args.Error(1)
	}

	return user, args.Error(1)
}

func TestLoginUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/login", auth.LoginUser)

	mockRepo := new(MockedRepository)
	UserRepo = mockRepo

	t.Run("binding error - missing fields", func(t *testing.T) {
		loginPayload, _ := json.Marshal(map[string]string{})
		request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(loginPayload))
		response := httptest.NewRecorder()

		r.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("error fetching user", func(t *testing.T) {
		mockRepo.On("GetUserByUsernameForLogin", "sadsadsa").Return(nil, errors.New("mock error"))

		loginPayload, _ := json.Marshal(map[string]string{
			"username": "sadsadsa",
			"password": "$2a$10$sssss2XwT5JjxkRmVCS37JnvuMOnSKKKUR/Se39mVpyIUk.PCs6VFU2F5C",
		})
		request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(loginPayload))
		response := httptest.NewRecorder()

		r.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

	t.Run("successful login", func(t *testing.T) {
		mockUser := &models.User{
			ID:       44,
			Username: "nihal",
			Email:    "nihal@gmail.com",
			Password: "$2a$10$notaDhRlmWZZ6OqKBMXBaOLty1TKK/qpvdA8nZB6Fh10tAwzMquEu",
		}

		mockRepo.On("GetUserByUsernameForLogin", "validuser").Return(mockUser, nil)

		loginPayload, _ := json.Marshal(map[string]string{
			"username": "validuser",
			"password": "validpassword",
		})
		request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(loginPayload))
		response := httptest.NewRecorder()

		r.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("invalid password", func(t *testing.T) {
		mockUser := &models.User{
			ID:       44,
			Username: "nihal",
			Email:    "nihal@gmail.com",
			Password: "$2a$10$notaDhRlmWZZ6OqKBMXBaOLty1TKK/qpvdA8nZB6Fh10tAwzMquEu",
		}
		mockRepo.On("GetUserByUsernameForLogin", "validuser").Return(mockUser, nil)

		loginPayload, _ := json.Marshal(map[string]string{
			"username": "validuser",
			"password": "wrongpassword",
		})
		request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(loginPayload))
		response := httptest.NewRecorder()

		r.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Code)
	})

}

func TestRegisterUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/register", auth.RegisterUser)

	mockRepo := new(MockedRepository)
	UserRepo = mockRepo
	t.Run("missing fields during registration", func(t *testing.T) {
		registerPayload, _ := json.Marshal(map[string]string{
			"username": "onlyusername",
		})
		request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(registerPayload))
		response := httptest.NewRecorder()

		r.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("username already exists", func(t *testing.T) {
		mockRepo.On("GetUserByUsername", "nihal").Return(&models.User{}, nil)

		registerPayload, _ := json.Marshal(map[string]string{
			"username": "nihal",
			"password": "password",
			"email":    "nihal@test.com",
		})
		request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(registerPayload))
		response := httptest.NewRecorder()

		r.ServeHTTP(response, request)

		assert.Equal(t, http.StatusConflict, response.Code)
	})

	t.Run("email already exists", func(t *testing.T) {
		mockRepo.On("GetUserByUsername", "newuser").Return(nil, errors.New("not found"))
		mockRepo.On("GetUserByEmail", "nihal@test.com").Return(&models.User{}, nil)

		registerPayload, _ := json.Marshal(map[string]string{
			"username": "newuser",
			"password": "password",
			"email":    "nihal@test.com",
		})
		request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(registerPayload))
		response := httptest.NewRecorder()

		r.ServeHTTP(response, request)

		assert.Equal(t, http.StatusConflict, response.Code)
	})
	t.Run("invalid email format", func(t *testing.T) {
		registerPayload, _ := json.Marshal(map[string]string{
			"username": "username",
			"password": "password",
			"email":    "invalidemail",
		})
		request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(registerPayload))
		response := httptest.NewRecorder()

		r.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
	t.Run("password too short", func(t *testing.T) {
		registerPayload, _ := json.Marshal(map[string]string{
			"username": "username",
			"password": "short",
			"email":    "user@test.com",
		})
		request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(registerPayload))
		response := httptest.NewRecorder()

		r.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("successful registration", func(t *testing.T) {
		mockRepo.On("GetUserByUsername", "newuser").Return(nil, errors.New("not found"))
		mockRepo.On("GetUserByEmail", "new@test.com").Return(nil, errors.New("not found"))
		mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)

		registerPayload, _ := json.Marshal(map[string]string{
			"username": "newuser",
			"password": "password",
			"email":    "new@test.com",
			"dob":      "2023-10-03",
		})
		request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(registerPayload))
		response := httptest.NewRecorder()

		r.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})
}
