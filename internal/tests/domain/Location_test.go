package domain_test

import (
	"github.com/samuel-go-expert/weather-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidZipCode(t *testing.T) {
	zipCode := "79745000"
	result := domain.IsValidZipCode(zipCode)

	assert.True(t, result, "Expected zip code %s to be valid", zipCode)
}

func TestIsInvalidZipCode(t *testing.T) {
	zipCode := "797450000"
	result := domain.IsValidZipCode(zipCode)

	assert.False(t, result, "Expected zip code %s to be valid", zipCode)
}
