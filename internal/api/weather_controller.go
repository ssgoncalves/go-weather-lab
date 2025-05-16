package api

import (
	"github.com/gin-gonic/gin"
	"github.com/samuel-go-expert/weather-api/internal/application"
	"github.com/samuel-go-expert/weather-api/internal/domain"
)

type WeatherController struct {
	weatherService application.WeatherServiceInterface
}

func NewWeatherController(ws application.WeatherServiceInterface) *WeatherController {
	return &WeatherController{
		weatherService: ws,
	}
}

func (controller *WeatherController) GetWeatherByZipCode(c *gin.Context) {

	weather, errWeather := controller.weatherService.GetWeatherByCity(c.Param("zipCode"))

	if errWeather == nil {
		c.JSON(200, weather)
		return
	}

	switch e := errWeather.(type) {
	case *domain.InvalidZipCodeError:
		c.String(422, e.Error())
	case *domain.AddressNotFoundError:
		c.String(404, e.Error())
	default:
		c.String(500, e.Error())
	}

}
