package infrastructure_test

import (
	"bytes"
	"fmt"
	"github.com/samuel-go-expert/weather-api/internal/domain"
	"github.com/samuel-go-expert/weather-api/internal/infrastructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAddress(t *testing.T) {
	mockHttp := new(HttpClientMock)
	cepApiClient := infrastructure.NewCepApi(mockHttp)

	mockResponse := httptest.NewRecorder()
	mockResponse.Code = http.StatusOK
	mockResponse.Body = bytes.NewBufferString(`{
		"localidade": "São Paulo",
		"Estado": "SP",
		"erro": ""
	}`)

	mockHttp.On("MakeGet", mock.Anything).Return(mockResponse.Result(), nil)

	address, err := cepApiClient.GetAddress("01333-003")

	assert.Nil(t, err)
	assert.Equal(t, "São Paulo", address.City)
	assert.Equal(t, "SP", address.State)
	mockHttp.AssertExpectations(t)
}

func TestGetAddressErrorOnMakeRequest(t *testing.T) {
	mockHttp := new(HttpClientMock)
	cepApiClient := infrastructure.NewCepApi(mockHttp)

	mockResponse := httptest.NewRecorder()
	mockResponse.Code = http.StatusOK
	mockResponse.Body = bytes.NewBufferString(`{
		"localidade": "São Paulo",
		"Estado": "SP",
		"erro": ""
	}`)

	mockHttp.On("MakeGet", mock.Anything).Return(mockResponse.Result(), fmt.Errorf("some error"))

	address, err := cepApiClient.GetAddress("01333-003")

	assert.Error(t, err)
	assert.Equal(t, domain.Address{}, address)
	mockHttp.AssertExpectations(t)
}

func TestGetAddressErrorForHttpStatusNotOk(t *testing.T) {
	mockHttp := new(HttpClientMock)
	cepApiClient := infrastructure.NewCepApi(mockHttp)

	mockResponse := httptest.NewRecorder()
	mockResponse.Code = http.StatusUnprocessableEntity
	mockResponse.Body = bytes.NewBufferString(`{}`)

	mockHttp.On("MakeGet", mock.Anything).Return(mockResponse.Result(), nil)

	address, err := cepApiClient.GetAddress("01333-003")

	assert.Error(t, err)
	assert.Equal(t, domain.Address{}, address)
	assert.Equal(t, http.StatusUnprocessableEntity, mockResponse.Code)
	assert.Equal(t, "failed request statusCode: 422", err.Error())
	mockHttp.AssertExpectations(t)
}

func TestGetAddressErrorOnDecodeBody(t *testing.T) {
	mockHttp := new(HttpClientMock)
	cepApiClient := infrastructure.NewCepApi(mockHttp)

	mockResponse := httptest.NewRecorder()
	mockResponse.Code = http.StatusOK
	mockResponse.Body = bytes.NewBufferString(`Not json`)

	mockHttp.On("MakeGet", mock.Anything).Return(mockResponse.Result(), nil)

	address, err := cepApiClient.GetAddress("01333-003")

	assert.Error(t, err)
	assert.Equal(t, domain.Address{}, address)
	assert.Equal(t, "failed to decode response", err.Error())
	mockHttp.AssertExpectations(t)
}

func TestGetAddressErrorReceiveError(t *testing.T) {
	mockHttp := new(HttpClientMock)
	cepApiClient := infrastructure.NewCepApi(mockHttp)

	mockResponse := httptest.NewRecorder()
	mockResponse.Code = http.StatusOK
	mockResponse.Body = bytes.NewBufferString(`{"erro": "some error"}`)

	mockHttp.On("MakeGet", mock.Anything).Return(mockResponse.Result(), nil)

	address, err := cepApiClient.GetAddress("01333-003")

	assert.Error(t, err)
	assert.Equal(t, domain.Address{}, address)
	assert.Equal(t, "can not find zipcode", err.Error())
	mockHttp.AssertExpectations(t)
}
