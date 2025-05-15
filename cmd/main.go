package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/samuel-go-expert/weather-api/internal/api"
	"github.com/samuel-go-expert/weather-api/internal/application"
	"github.com/samuel-go-expert/weather-api/internal/infrastructure"
	"os"
)

func main() {

	loadEnvs()

	router := gin.Default()
	router.GET("/zip-code/:zipCode/weather", getController().GetWeatherByZipCode)

	err := router.Run(":8080")

	if err != nil {
		fmt.Printf("error starting server: %v", err)
		return
	}
}

func getController() (WeatherController *api.WeatherController) {

	httpClient := infrastructure.NewHttpClient()
	env := infrastructure.NewEnv()

	weatherApi := infrastructure.NewWeatherAPI(httpClient, env)
	weatherService := application.NewWeatherService(weatherApi)

	addressApi := infrastructure.NewCepApi(httpClient)
	addressService := application.NewAddressService(addressApi)

	return api.NewWeatherController(weatherService, addressService)
}

func loadEnvs() bool {
	if os.Getenv("LOAD_ENV") != "true" {
		return false
	}

	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: .env file not found, proceeding with default environment variables\n")
		return false
	}

	fmt.Println("Loaded .env file")

	return true
}
