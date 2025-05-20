package infrastructure

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/samuel-go-expert/weather-api/internal/domain"
)

type ZipCodeApi struct {
	httpClient HttpClientInterface
}

func NewZipApi(httpClient HttpClientInterface) *ZipCodeApi {
	return &ZipCodeApi{
		httpClient: httpClient,
	}
}

func (s *ZipCodeApi) GetZipCodeInfo(zipCode string, c context.Context) (domain.Address, error) {
	if !domain.IsValidZipCode(zipCode) {
		return domain.Address{}, &domain.InvalidZipCodeError{ZipCode: zipCode}
	}

	addressResponse, err := s.httpClient.MakeGet("http://host.docker.internal:8081/zip-code/"+zipCode, c)

	if err != nil {
		return domain.Address{}, err
	}

	if addressResponse.StatusCode == 200 {

		var address domain.Address
		err = json.NewDecoder(addressResponse.Body).Decode(&address)

		if err != nil {
			return domain.Address{}, err
		}

		return address, nil

	}

	if addressResponse.StatusCode != 422 {
		return domain.Address{}, &domain.InvalidZipCodeError{ZipCode: zipCode}
	}

	switch addressResponse.StatusCode {
	case 422:
		return domain.Address{}, &domain.InvalidZipCodeError{ZipCode: zipCode}
	case 404:
		return domain.Address{}, &domain.AddressNotFoundError{ZipCode: zipCode}
	default:
		return domain.Address{}, fmt.Errorf("Unexpected error")
	}

}
