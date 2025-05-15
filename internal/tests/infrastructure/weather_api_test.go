package infrastructure_test

import (
	"bytes"
	"fmt"
	"github.com/samuel-go-expert/weather-api/internal/domain"
	"github.com/samuel-go-expert/weather-api/internal/infrastructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type EnvMock struct {
	mock.Mock
}

func (e *EnvMock) Getenv(key string) string {
	args := e.Called(key)
	return args.String(0)
}

func TestGetWeather(t *testing.T) {
	httpMock := new(HttpClientMock)
	envMock := new(EnvMock)
	weatherApiClient := infrastructure.NewWeatherAPI(httpMock, envMock)

	mockResponse := httptest.NewRecorder()
	mockResponse.Code = http.StatusOK
	mockResponse.Body = bytes.NewBufferString(`{
		"current": {
			"temp_c": 15,
			"temp_f": 85
		},
		"location": {
			"name": "São Paulo",
			"region": "SP"
		}
	}`)

	envMock.On("Getenv", "WEATHER_API_KEY").Return("some_api_key")

	httpMock.On("MakeGet", mock.Anything).Return(mockResponse.Result(), nil)

	weatherResponse, err := weatherApiClient.GetWeather("São Paulo")

	assert.Nil(t, err)
	assert.Equal(t, "São Paulo", weatherResponse.City)
	assert.Equal(t, "SP", weatherResponse.State)
	assert.Equal(t, 15.0, weatherResponse.Celsius)
	assert.Equal(t, 85.0, weatherResponse.Fahrenheit)
	assert.Equal(t, 288.15, weatherResponse.Kelvin)
	httpMock.AssertExpectations(t)
}

func TestGetErrorWhenApiKeyWasNotFound(t *testing.T) {
	httpMock := new(HttpClientMock)
	envMock := new(EnvMock)
	weatherApiClient := infrastructure.NewWeatherAPI(httpMock, envMock)

	envMock.On("Getenv", "WEATHER_API_KEY").Return("")

	weatherResponse, err := weatherApiClient.GetWeather("São Paulo")

	assert.Error(t, err)
	assert.Equal(t, "API key not found", err.Error())
	assert.Equal(t, domain.Weather{}, weatherResponse)
	httpMock.AssertExpectations(t)
}

func TestGetErrorOnFailedRequest(t *testing.T) {

	// Set
	httpMock := new(HttpClientMock)
	envMock := new(EnvMock)
	weatherApiClient := infrastructure.NewWeatherAPI(httpMock, envMock)

	// Expectations
	envMock.On("Getenv", "WEATHER_API_KEY").Return("some_api_key")
	httpMock.On("MakeGet", mock.Anything).Return(httptest.NewRecorder().Result(), fmt.Errorf("some error"))

	// Actions
	weatherResponse, err := weatherApiClient.GetWeather("São Paulo")

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, "failed to fetch weather data for location São Paulo: some error", err.Error())
	assert.Equal(t, domain.Weather{}, weatherResponse)
	httpMock.AssertExpectations(t)
}

func TestGetErrorWhenNotReceivedAnOkHttpStatus(t *testing.T) {
	// Set
	httpMock := new(HttpClientMock)
	envMock := new(EnvMock)
	weatherApiClient := infrastructure.NewWeatherAPI(httpMock, envMock)

	mockResponse := httptest.NewRecorder()
	mockResponse.Code = http.StatusUnprocessableEntity
	mockResponse.Body = bytes.NewBufferString(`{}`)

	// Expectations
	envMock.On("Getenv", "WEATHER_API_KEY").Return("some_api_key")
	httpMock.On("MakeGet", mock.Anything).Return(mockResponse.Result(), nil)

	// Actions
	weatherResponse, err := weatherApiClient.GetWeather("São Paulo")

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, "unexpected status code 422 when fetching weather data for location São Paulo", err.Error())
	assert.Equal(t, domain.Weather{}, weatherResponse)
	httpMock.AssertExpectations(t)
}

func TestGetErrorWhenReceivedAnInvalidWeatherResponse(t *testing.T) {
	// Set
	httpMock := new(HttpClientMock)
	envMock := new(EnvMock)
	weatherApiClient := infrastructure.NewWeatherAPI(httpMock, envMock)

	mockResponse := httptest.NewRecorder()
	mockResponse.Code = http.StatusOK
	mockResponse.Body = bytes.NewBufferString(`{"current": "comma after the value is invalid",}`)

	// Expectations
	envMock.On("Getenv", "WEATHER_API_KEY").Return("some_api_key")
	httpMock.On("MakeGet", mock.Anything).Return(mockResponse.Result(), nil)

	// Actions
	weatherResponse, err := weatherApiClient.GetWeather("São Paulo")

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, "failed to decode weather data for location São Paulo: invalid character '}' looking for beginning of object key string", err.Error())
	assert.Equal(t, domain.Weather{}, weatherResponse)
	httpMock.AssertExpectations(t)
}
