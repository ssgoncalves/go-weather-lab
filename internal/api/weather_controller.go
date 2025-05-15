package api

import (
	"github.com/gin-gonic/gin"
	"github.com/samuel-go-expert/weather-api/internal/application"
	"github.com/samuel-go-expert/weather-api/internal/domain"
)

type WeatherController struct {
	weatherService application.WeatherServiceInterface
	addressService application.AddressServiceInterface
}

func NewWeatherController(ws application.WeatherServiceInterface, as application.AddressServiceInterface) *WeatherController {
	return &WeatherController{
		weatherService: ws,
		addressService: as,
	}
}

func (controller *WeatherController) GetWeatherByZipCode(c *gin.Context) {
	zipCode := c.Param("zipCode")

	if (domain.IsValidZipCode(zipCode)) == false {
		c.String(422, "invalid zipcode")
		return
	}

	address, errAddress := controller.addressService.GetAddressByZipCode(zipCode)

	if errAddress != nil {

		switch e := errAddress.(type) {
		case *domain.InvalidZipCodeError:
			c.String(422, e.Error())
		case *domain.AddressNotFoundError:
			c.String(404, e.Error())
		default:

			c.String(500, e.Error())
		}

		return
	}

	weather, errWeather := controller.weatherService.GetWeatherByCity(address.City)

	if errWeather != nil {
		c.String(502, errWeather.Error())
		return
	}

	c.JSON(200, weather)
}
