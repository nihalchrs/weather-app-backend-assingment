package repositories

import (
	"github.com/stretchr/testify/mock"

	"weather-app/models"
)

type WeatherRepository interface {
	GetWeatherDataForUser(city string, userID int) (*models.Weather, error)
	GetWeatherHistoryByUserID(userID int) ([]models.Weather, error)
	CheckWeatherBelongsToUser(userID, weatherID int) (bool, error)
	UpdateWeatherData(userID, weatherID int, city string) error
	DeleteWeatherData(userID, weatherID int) error
	BulkDeleteWeatherData(ids []int) error
}

type MockWeatherRepository struct {
	mock.Mock
}

func (m *MockWeatherRepository) GetWeatherDataForUser(city string, userID int) (*models.Weather, error) {
	args := m.Called(city, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Weather), args.Error(1)
}

func (m *MockWeatherRepository) GetWeatherHistoryByUserID(userID int) ([]models.Weather, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Weather), args.Error(1)
}

func (m *MockWeatherRepository) CheckWeatherBelongsToUser(userID, weatherID int) (bool, error) {
	args := m.Called(userID, weatherID)
	return args.Bool(0), args.Error(1)
}

func (m *MockWeatherRepository) UpdateWeatherData(userID, weatherID int, city string) error {
	args := m.Called(userID, weatherID, city)
	return args.Error(0)
}

func (m *MockWeatherRepository) DeleteWeatherData(userID, weatherID int) error {
	args := m.Called(userID, weatherID)
	return args.Error(0)
}

func (m *MockWeatherRepository) BulkDeleteWeatherData(ids []int) error {
	args := m.Called(ids)
	return args.Error(0)
}
