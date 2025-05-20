package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/samuel-go-expert/weather-api/internal/api"
	"github.com/samuel-go-expert/weather-api/internal/application"
	"github.com/samuel-go-expert/weather-api/internal/infrastructure"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func main() {

	loadEnvs()

	grpcCon, err := getGrpConnection("otel-collector:4317")
	defer grpcCon.Close()

	if err != nil {
		fmt.Printf("error creating gRPC connection: %v", err)
		return
	}

	ctx := context.Background()
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(grpcCon))

	if err != nil {
		fmt.Printf("failed to create trace exporter: %w", err)
		return
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	tracer := otel.Tracer(os.Getenv("TITLE"))

	router := gin.Default()
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/zip-code/:zipCode/weather", getWeatherController(tracer).GetWeatherByZipCode)
	router.GET("/zip-code/:zipCode", getZipCodeController(tracer).GetZipCodeInfo)

	err = router.Run(":8080")

	if err != nil {
		fmt.Printf("error starting server: %v", err)
		return
	}

}

func getZipCodeController(trace trace.Tracer) (ZipCodeController *api.ZipCodeController) {

	httpClient := infrastructure.NewHttpClient()

	addressApi := infrastructure.NewViaCepApi(httpClient)
	addressService := application.NewAddressService(addressApi)

	return api.NewZipCodeController(addressService, trace)

}

func getWeatherController(trace trace.Tracer) (WeatherController *api.WeatherController) {

	httpClient := infrastructure.NewHttpClient()
	env := infrastructure.NewEnv()

	zipCodeApi := infrastructure.NewZipApi(httpClient)
	zipCodeService := application.NewZipCodeService(zipCodeApi)

	weatherApi := infrastructure.NewWeatherAPI(httpClient, env)
	weatherService := application.NewWeatherService(weatherApi, zipCodeService)
	return api.NewWeatherController(weatherService, trace)
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

func getGrpConnection(collectorURL string) (*grpc.ClientConn, error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	return grpc.NewClient(
		collectorURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

}
