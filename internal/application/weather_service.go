package application

import (
	"github.com/samuel-go-expert/weather-api/internal/domain"
	"github.com/samuel-go-expert/weather-api/internal/infrastructure"
)

type WeatherServiceInterface interface {
	GetWeatherByCity(city string) (domain.Weather, error)
}

type WeatherService struct {
	weatherApi infrastructure.WeatherApi
}

func NewWeatherService(weatherApi infrastructure.WeatherApi) *WeatherService {
	return &WeatherService{
		weatherApi: weatherApi,
	}
}

func (s *WeatherService) GetWeatherByCity(city string) (domain.Weather, error) {
	weather, err := s.weatherApi.GetWeather(city)

	if err != nil {
		return domain.Weather{}, err
	}

	return weather, nil

}
