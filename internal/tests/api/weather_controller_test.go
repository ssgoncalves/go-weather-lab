package api_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/samuel-go-expert/weather-api/internal/api"
	"github.com/samuel-go-expert/weather-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type weatherServiceMock struct {
	mock.Mock
}

func (m *weatherServiceMock) GetWeatherByCity(city string) (domain.Weather, error) {
	args := m.Called(city)
	return args.Get(0).(domain.Weather), args.Error(1)
}

type addressServiceMock struct {
	mock.Mock
}

func (m *addressServiceMock) GetAddressByZipCode(zipCode string) (domain.Address, error) {
	args := m.Called(zipCode)
	return args.Get(0).(domain.Address), args.Error(1)
}

func TestGetWeatherByZipCode(t *testing.T) {

	addressServiceMock := new(addressServiceMock)
	weatherServiceMock := new(weatherServiceMock)
	api.NewWeatherController(weatherServiceMock, addressServiceMock)

	addressServiceMock.On("GetAddressByZipCode", "13503100").Return(domain.Address{
		City:  "São Carlos",
		State: "SP",
	}, nil)

	weatherServiceMock.On("GetWeatherByCity", "São Carlos").Return(domain.Weather{
		City:       "São Carlos",
		State:      "SP",
		Celsius:    25.0,
		Fahrenheit: 77.0,
		Kelvin:     298.15,
	}, nil)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/weather/:zipCode", api.NewWeatherController(weatherServiceMock, addressServiceMock).GetWeatherByZipCode)

	req, _ := http.NewRequest(http.MethodGet, "/weather/13503100", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), `{"city":"São Carlos","state":"SP","temp_c":25,"temp_f":77,"temp_k":298.15}`)

	addressServiceMock.AssertExpectations(t)
	weatherServiceMock.AssertExpectations(t)
}

func TestGetErrorForInvalidZipCode(t *testing.T) {
	addressServiceMock := new(addressServiceMock)
	weatherServiceMock := new(weatherServiceMock)
	api.NewWeatherController(weatherServiceMock, addressServiceMock)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/weather/:zipCode", api.NewWeatherController(weatherServiceMock, addressServiceMock).GetWeatherByZipCode)

	req, _ := http.NewRequest(http.MethodGet, "/weather/1350310000", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, `invalid zipcode`, w.Body.String())
}

func TestGetErrorForAddressNotFound(t *testing.T) {

	addressServiceMock := new(addressServiceMock)
	weatherServiceMock := new(weatherServiceMock)
	api.NewWeatherController(weatherServiceMock, addressServiceMock)

	addressServiceMock.On("GetAddressByZipCode", "13503100").Return(domain.Address{}, &domain.AddressNotFoundError{ZipCode: "13503100"})

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/weather/:zipCode", api.NewWeatherController(weatherServiceMock, addressServiceMock).GetWeatherByZipCode)

	req, _ := http.NewRequest(http.MethodGet, "/weather/13503100", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, `can not find zipcode`, w.Body.String())

	addressServiceMock.AssertExpectations(t)
}

func TestGetErrorForInternalServerErro(t *testing.T) {

	addressServiceMock := new(addressServiceMock)
	weatherServiceMock := new(weatherServiceMock)
	api.NewWeatherController(weatherServiceMock, addressServiceMock)

	addressServiceMock.On("GetAddressByZipCode", "13503100").Return(domain.Address{}, fmt.Errorf("internal server error"))

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/weather/:zipCode", api.NewWeatherController(weatherServiceMock, addressServiceMock).GetWeatherByZipCode)

	req, _ := http.NewRequest(http.MethodGet, "/weather/13503100", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, `internal server error`, w.Body.String())

	addressServiceMock.AssertExpectations(t)
}
