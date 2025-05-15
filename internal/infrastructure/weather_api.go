package infrastructure

import (
	"encoding/json"
	"fmt"
	"github.com/samuel-go-expert/weather-api/internal/domain"
	"log"
	"net/http"
	"net/url"
)

const (
	weatherApiUrl     = "https://api.weatherapi.com/v1/current.json"
	weatherApiKeyName = "WEATHER_API_KEY"
)

type WeatherResponse struct {
	Location struct {
		City  string `json:"name"`
		State string `json:"region"`
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}

type WeatherAPIClient struct {
	httpClient HttpClientInterface
	env        EnvInterface
}

type WeatherApi interface {
	GetWeather(city string) (domain.Weather, error)
}

func NewWeatherAPI(h HttpClientInterface, e EnvInterface) WeatherApi {
	return &WeatherAPIClient{
		httpClient: h,
		env:        e,
	}
}

func (c *WeatherAPIClient) GetWeather(city string) (domain.Weather, error) {

	apiKey := c.env.Getenv(weatherApiKeyName)

	if apiKey == "" {
		return domain.Weather{}, fmt.Errorf("API key not found")
	}

	response, err := c.httpClient.MakeGet(buildWeatherApiUrl(apiKey, city))

	if err != nil {
		return domain.Weather{}, fmt.Errorf("failed to fetch weather data for location %s: %v", city, err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return domain.Weather{}, fmt.Errorf("unexpected status code %d when fetching weather data for location %s", response.StatusCode, city)
	}

	var weatherResponse WeatherResponse

	if err := json.NewDecoder(response.Body).Decode(&weatherResponse); err != nil {
		log.Print("Error decoding weather response:", err)
		return domain.Weather{}, fmt.Errorf("failed to decode weather data for location %s: %v", city, err)
	}

	return domain.Weather{
		City:       weatherResponse.Location.City,
		State:      weatherResponse.Location.State,
		Celsius:    weatherResponse.Current.TempC,
		Fahrenheit: weatherResponse.Current.TempF,
		Kelvin:     weatherResponse.Current.TempC + 273.15,
	}, nil
}

func buildWeatherApiUrl(apiKey, location string) string {
	return fmt.Sprintf("%s?key=%s&q=%s&lang=pt", weatherApiUrl, apiKey, url.QueryEscape(location))
}
