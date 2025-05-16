package application

import (
	"github.com/samuel-go-expert/weather-api/internal/domain"
)

type AddressApi interface {
	GetAddress(cep string) (domain.Address, error)
}

type AddressServiceInterface interface {
	GetAddressByZipCode(zipCode string) (domain.Address, error)
}

type AddressService struct {
	addressApi AddressApi
}

func NewAddressService(addressApi AddressApi) *AddressService {
	return &AddressService{
		addressApi: addressApi,
	}
}

func (s *AddressService) GetAddressByZipCode(zipCode string) (domain.Address, error) {

	if !domain.IsValidZipCode(zipCode) {
		return domain.Address{}, &domain.InvalidZipCodeError{ZipCode: zipCode}
	}

	location, err := s.addressApi.GetAddress(zipCode)

	if err != nil {
		return domain.Address{}, err
	}

	return domain.Address{
		City:  location.City,
		State: location.State,
	}, nil

}
