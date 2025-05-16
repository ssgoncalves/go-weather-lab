package application

import (
	"github.com/samuel-go-expert/weather-api/internal/domain"
)

type ZipCodeServiceInterface interface {
	GetAddressByZipCode(zipCode string) (domain.Address, error)
}

type ZipCodeApiInterface interface {
	GetZipCodeInfo(zipCode string) (domain.Address, error)
}

type ZipCodeService struct {
	addressApi ZipCodeApiInterface
}

func NewZipCodeService(zipCodeApi ZipCodeApiInterface) *ZipCodeService {
	return &ZipCodeService{
		addressApi: zipCodeApi,
	}
}

func (s *ZipCodeService) GetAddressByZipCode(zipCode string) (domain.Address, error) {

	return s.addressApi.GetZipCodeInfo(zipCode)

}
