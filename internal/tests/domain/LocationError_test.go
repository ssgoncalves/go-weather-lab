package domain_test

import (
	"github.com/samuel-go-expert/weather-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvalidZipCodeError(t *testing.T) {
	zipCode := "79745000"
	err := &domain.InvalidZipCodeError{ZipCode: zipCode}

	assert.Equal(t, "invalid zipcode", err.Error())
}

func TestAddressNotFoundError(t *testing.T) {
	zipCode := "79745000"
	err := &domain.AddressNotFoundError{ZipCode: zipCode}

	assert.Equal(t, "can not find zipcode", err.Error())
}
