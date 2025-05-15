package application_test

import (
	"fmt"
	"github.com/samuel-go-expert/weather-api/internal/application"
	"github.com/samuel-go-expert/weather-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type AddressApiMock struct {
	mock.Mock
}

func (m *AddressApiMock) GetAddress(zipCode string) (domain.Address, error) {
	args := m.Called(zipCode)
	return args.Get(0).(domain.Address), args.Error(1)
}

func TestGetAddressByZipCode(t *testing.T) {
	addressApiMock := new(AddressApiMock)
	addressService := application.NewAddressService(addressApiMock)

	zipCode := "13503100"
	expectedLocation := domain.Address{
		City:  "São Carlos",
		State: "SP",
	}

	addressApiMock.On("GetAddress", zipCode).Return(expectedLocation, nil)

	location, err := addressService.GetAddressByZipCode(zipCode)

	assert.NoError(t, err)
	assert.Equal(t, expectedLocation, location)
	addressApiMock.AssertExpectations(t)
}

func TestGetErrorForInvalidZipCode(t *testing.T) {
	addressApiMock := new(AddressApiMock)
	addressService := application.NewAddressService(addressApiMock)

	zipCode := "135031000"

	address, err := addressService.GetAddressByZipCode(zipCode)

	assert.Equal(t, domain.Address{}, address)
	assert.Error(t, err)
	assert.IsType(t, &domain.InvalidZipCodeError{}, err)
	assert.Equal(t, "invalid zipcode", err.Error())
}

func TestGetErrorForFailedGetAddressByZipCodeApiCall(t *testing.T) {
	addressApiMock := new(AddressApiMock)
	addressService := application.NewAddressService(addressApiMock)

	zipCode := "13503100"
	expectedError := fmt.Errorf("something went wrong")

	addressApiMock.On("GetAddress", zipCode).Return(domain.Address{}, expectedError)

	address, err := addressService.GetAddressByZipCode(zipCode)

	assert.Equal(t, domain.Address{}, address)
	assert.Error(t, err)
	assert.Equal(t, "something went wrong", err.Error())
}
