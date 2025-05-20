package api

import (
	"github.com/gin-gonic/gin"
	"github.com/samuel-go-expert/weather-api/internal/application"
	"github.com/samuel-go-expert/weather-api/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type ZipCodeController struct {
	addressService application.AddressServiceInterface
	tracer         trace.Tracer
}

func NewZipCodeController(as application.AddressServiceInterface, tracer trace.Tracer) *ZipCodeController {
	return &ZipCodeController{
		addressService: as,
		tracer:         tracer,
	}
}

func (controller *ZipCodeController) GetZipCodeInfo(c *gin.Context) {

	carrier := propagation.HeaderCarrier(c.Request.Header)
	ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), carrier)

	ctx, span := controller.tracer.Start(ctx, "GetZipCodeData")
	defer span.End()

	zipCode := c.Param("zipCode")

	if !domain.IsValidZipCode(zipCode) {
		c.String(422, "invalid zipcode")
		return
	}

	ctx2, span2 := controller.tracer.Start(ctx, "GetZipCodeData - External Call")
	address, errAddress := controller.addressService.GetAddressByZipCode(zipCode, ctx2)
	span2.End()

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

	c.JSON(200, address)

	// Propagar o contexto de trace para a requisição de volta
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(c.Request.Header))
}
