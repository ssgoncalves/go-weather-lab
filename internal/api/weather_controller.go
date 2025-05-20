package api

import (
	"github.com/gin-gonic/gin"
	"github.com/samuel-go-expert/weather-api/internal/application"
	"github.com/samuel-go-expert/weather-api/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type WeatherController struct {
	weatherService application.WeatherServiceInterface
	tracer         trace.Tracer
}

func NewWeatherController(ws application.WeatherServiceInterface, tracer trace.Tracer) *WeatherController {
	return &WeatherController{
		weatherService: ws,
		tracer:         tracer,
	}
}

func (controller *WeatherController) GetWeatherByZipCode(c *gin.Context) {
	// Extração do contexto de trace da requisição
	carrier := propagation.HeaderCarrier(c.Request.Header)
	ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), carrier)

	// Criar o novo span para a operação
	ctx, span := controller.tracer.Start(ctx, "GetWeatherByZipCode")
	defer span.End()

	// Passar o contexto de trace para o serviço
	weather, errWeather := controller.weatherService.GetWeatherByCity(c.Param("zipCode"), ctx)

	if errWeather == nil {
		c.JSON(200, weather)
		span.AddEvent("Weather retrieved")
		return
	}

	// Handle errors as needed
	switch e := errWeather.(type) {
	case *domain.InvalidZipCodeError:
		c.String(422, e.Error())
	case *domain.AddressNotFoundError:
		c.String(404, e.Error())
	default:
		c.String(500, e.Error())
	}

	// Propagar o contexto para a requisição
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(c.Request.Header))
}
