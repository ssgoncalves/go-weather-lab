package application_test

import (
	"fmt"
	"github.com/samuel-go-expert/weather-api/internal/application"
	"github.com/samuel-go-expert/weather-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type WeatherApiMock struct {
	mock.Mock
}

func (m *WeatherApiMock) GetWeather(city string) (domain.Weather, error) {
	args := m.Called(city)
	return args.Get(0).(domain.Weather), args.Error(1)
}

func TestGetWeatherByCity(t *testing.T) {
	weatherApiMock := new(WeatherApiMock)
	weatherService := application.NewWeatherService(weatherApiMock)

	city := "São Carlos"
	expectedWeather := domain.Weather{
		City:       "São Carlos",
		State:      "SP",
		Celsius:    25,
		Fahrenheit: 77,
		Kelvin:     298.15,
	}

	weatherApiMock.On("GetWeather", city).Return(expectedWeather, nil)

	weather, err := weatherService.GetWeatherByCity(city)

	assert.NoError(t, err)
	assert.Equal(t, "São Carlos", weather.City)
	assert.Equal(t, "SP", weather.State)
	assert.Equal(t, 25.0, weather.Celsius)
	assert.Equal(t, 77.0, weather.Fahrenheit)
	assert.Equal(t, 298.15, weather.Kelvin)
	weatherApiMock.AssertExpectations(t)
}

func TestGetErrorForFailedApiGetWeatherByCityCall(t *testing.T) {
	weatherApiMock := new(WeatherApiMock)
	weatherService := application.NewWeatherService(weatherApiMock)

	city := "São Carlos"
	expectedError := fmt.Errorf("something went wrong")

	weatherApiMock.On("GetWeather", city).Return(domain.Weather{}, expectedError)

	weather, err := weatherService.GetWeatherByCity(city)

	assert.Equal(t, domain.Weather{}, weather)
	assert.Error(t, err)
	assert.Equal(t, "something went wrong", err.Error())

	weatherApiMock.AssertExpectations(t)
}
