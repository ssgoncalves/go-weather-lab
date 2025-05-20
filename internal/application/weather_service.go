package application

import (
	"context"
	"github.com/samuel-go-expert/weather-api/internal/domain"
)

type WeatherApi interface {
	GetWeather(city string, c context.Context) (domain.Weather, error)
}

type WeatherServiceInterface interface {
	GetWeatherByCity(city string, c context.Context) (domain.Weather, error)
}

type WeatherService struct {
	weatherApi     WeatherApi
	zipCodeService ZipCodeServiceInterface
}

func NewWeatherService(weatherApi WeatherApi, zipCodeService ZipCodeServiceInterface) *WeatherService {
	return &WeatherService{
		weatherApi:     weatherApi,
		zipCodeService: zipCodeService,
	}
}

func (s *WeatherService) GetWeatherByCity(zipCode string, c context.Context) (domain.Weather, error) {

	address, zipCodeErro := s.zipCodeService.GetAddressByZipCode(zipCode, c)

	if zipCodeErro != nil {
		return domain.Weather{}, zipCodeErro
	}

	weather, err := s.weatherApi.GetWeather(address.City, c)

	if err != nil {
		return domain.Weather{}, err
	}

	return weather, nil

}
