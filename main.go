package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

const (
	weatherApiUrl     = "https://api.weatherapi.com/v1/current.json"
	weatherApiKeyName = "WEATHER_API_KEY"
	cepApiUrl         = "https://viacep.com.br/ws/"
)

type ApiResponse struct {
	City       string  `json:"city"`
	State      string  `json:"state"`
	Celsius    float64 `json:"temp_c"`
	Fahrenheit float64 `json:"temp_f"`
	Kelvin     float64 `json:"temp_k"`
}

type LocationResponse struct {
	City  string `json:"localidade"`
	State string `json:"Estado"`
}

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

func main() {

	zipCode := "79041050"

	err := loadEnvs()

	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	response, err := getWeatherFromZipCode(zipCode)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Cidade:", response.City)
	fmt.Println("Estado:", response.State)
	fmt.Printf("Temperature in Celsius: %.2f\n", response.Celsius)
	fmt.Printf("Temperature in Fahrenheit: %.2f\n", response.Fahrenheit)
	fmt.Printf("Temperature in Kelvin: %.2f\n", response.Kelvin)

}

func loadEnvs() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}
	return nil
}

func getWeatherFromZipCode(zipCode string) (ApiResponse, error) {

	locationResponse, err := getLocation(zipCode)

	if err != nil {
		return ApiResponse{}, err
	}

	weatherResponse, err := getWeather(locationResponse.City)

	if err != nil {
		return ApiResponse{}, err
	}

	return ApiResponse{
		City:       weatherResponse.Location.City,
		State:      weatherResponse.Location.State,
		Celsius:    weatherResponse.Current.TempC,
		Fahrenheit: weatherResponse.Current.TempF,
		Kelvin:     weatherResponse.Current.TempC + 273.15,
	}, nil
}

func getLocation(zipCode string) (LocationResponse, error) {

	response, err := http.Get(buildLocationApiUrl(zipCode))

	if err != nil {
		return LocationResponse{}, fmt.Errorf("failed to fetch location data for zip code %s: %v", zipCode, err)
	}

	defer response.Body.Close()

	var locationResponse LocationResponse

	err = json.NewDecoder(response.Body).Decode(&locationResponse)

	if err != nil {
		return LocationResponse{}, fmt.Errorf("failed to decode location data for zip code %s: %v", zipCode, err)
	}

	return locationResponse, nil
}

func buildLocationApiUrl(zipCode string) string {
	return fmt.Sprintf("%s%s/json/", cepApiUrl, zipCode)
}

func getWeather(location string) (WeatherResponse, error) {

	apiKey := os.Getenv(weatherApiKeyName)

	if apiKey == "" {
		return WeatherResponse{}, fmt.Errorf("API key not found")
	}

	response, err := http.Get(buildWeatherApiUrl(apiKey, location))

	if err != nil {
		return WeatherResponse{}, fmt.Errorf("failed to fetch weather data for location %s: %v", location, err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return WeatherResponse{}, fmt.Errorf("unexpected status code %d when fetching weather data for location %s", response.StatusCode, location)
	}

	var weatherResponse WeatherResponse

	if err := json.NewDecoder(response.Body).Decode(&weatherResponse); err != nil {
		return WeatherResponse{}, fmt.Errorf("failed to decode weather data for location %s: %v", location, err)
	}

	return weatherResponse, nil
}

func buildWeatherApiUrl(apiKey, location string) string {
	return fmt.Sprintf("%s?key=%s&q=%s&lang=pt", weatherApiUrl, apiKey, url.QueryEscape(location))
}
