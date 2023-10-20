package weather_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"weather-app/models"
	"weather-app/repositories"
)

func TestNewMockWeatherRepository(t *testing.T) {
	mockRepo := new(repositories.MockWeatherRepository)

	t.Run("GetWeatherDataForUser", func(t *testing.T) {
		expectedWeather := &models.Weather{City: "TestCity", Temp: 20.5}
		mockRepo.On("GetWeatherDataForUser", "TestCity", 1).Return(expectedWeather, nil)

		result, err := mockRepo.GetWeatherDataForUser("TestCity", 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedWeather, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetWeatherHistoryByUserID", func(t *testing.T) {
		expectedWeatherHistory := []models.Weather{{City: "TestCity", Temp: 20.5}}
		mockRepo.On("GetWeatherHistoryByUserID", 1).Return(expectedWeatherHistory, nil)

		result, err := mockRepo.GetWeatherHistoryByUserID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedWeatherHistory, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("CheckWeatherBelongsToUser", func(t *testing.T) {
		mockRepo.On("CheckWeatherBelongsToUser", 1, 1001).Return(true, nil)

		isOwner, err := mockRepo.CheckWeatherBelongsToUser(1, 1001)

		assert.NoError(t, err)
		assert.True(t, isOwner)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateWeatherData", func(t *testing.T) {
		mockRepo.On("UpdateWeatherData", 1, 1001, "UpdatedCity").Return(nil)

		err := mockRepo.UpdateWeatherData(1, 1001, "UpdatedCity")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteWeatherData", func(t *testing.T) {
		mockRepo.On("DeleteWeatherData", 1, 1001).Return(nil)

		err := mockRepo.DeleteWeatherData(1, 1001)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("BulkDeleteWeatherData", func(t *testing.T) {
		ids := []int{1001, 1002, 1003}
		mockRepo.On("BulkDeleteWeatherData", ids).Return(nil)

		err := mockRepo.BulkDeleteWeatherData(ids)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorInGetWeatherDataForUser", func(t *testing.T) {
		mockRepo.On("GetWeatherDataForUser", "TestCity", 1).Return(nil, errors.New("fetch error"))

		_, err := mockRepo.GetWeatherDataForUser("TestCity", 1)

		assert.Error(t, err)
		assert.Equal(t, "fetch error", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorInCheckWeatherBelongsToUser", func(t *testing.T) {
		mockRepo.On("CheckWeatherBelongsToUser", 1, 1002).Return(false, errors.New("check error"))

		_, err := mockRepo.CheckWeatherBelongsToUser(1, 1002)

		assert.Error(t, err)
		assert.Equal(t, "check error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
