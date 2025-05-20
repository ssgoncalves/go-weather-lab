package application

import (
	"context"
	"github.com/samuel-go-expert/weather-api/internal/domain"
)

type ZipCodeServiceInterface interface {
	GetAddressByZipCode(zipCode string, c context.Context) (domain.Address, error)
}

type ZipCodeApiInterface interface {
	GetZipCodeInfo(zipCode string, c context.Context) (domain.Address, error)
}

type ZipCodeService struct {
	addressApi ZipCodeApiInterface
}

func NewZipCodeService(zipCodeApi ZipCodeApiInterface) *ZipCodeService {
	return &ZipCodeService{
		addressApi: zipCodeApi,
	}
}

func (s *ZipCodeService) GetAddressByZipCode(zipCode string, c context.Context) (domain.Address, error) {

	return s.addressApi.GetZipCodeInfo(zipCode, c)

}
