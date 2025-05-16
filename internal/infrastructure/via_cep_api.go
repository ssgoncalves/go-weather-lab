package infrastructure

import (
	"encoding/json"
	"fmt"
	"github.com/samuel-go-expert/weather-api/internal/application"
	"github.com/samuel-go-expert/weather-api/internal/domain"
)

type LocationResponse struct {
	City  string `json:"localidade"`
	State string `json:"Estado"`
	Error string `json:"erro"`
}

type CepApiClient struct {
	httpClient HttpClientInterface
}

func NewViaCepApi(h HttpClientInterface) application.AddressApi {
	return &CepApiClient{
		httpClient: h,
	}
}

func (c *CepApiClient) GetAddress(zipCode string) (domain.Address, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", zipCode)

	response, err := c.httpClient.MakeGet(url)

	if err != nil {
		return domain.Address{}, err
	}

	defer response.Body.Close()

	var viaCepLocation LocationResponse

	if response.StatusCode != 200 {
		return domain.Address{}, fmt.Errorf("failed request statusCode: %d", response.StatusCode)
	}

	err = json.NewDecoder(response.Body).Decode(&viaCepLocation)

	if err != nil {
		return domain.Address{}, fmt.Errorf("failed to decode response")
	}

	if viaCepLocation.Error != "" {
		return domain.Address{}, &domain.AddressNotFoundError{ZipCode: zipCode}
	}

	return domain.Address{
		City:  viaCepLocation.City,
		State: viaCepLocation.State,
	}, nil
}
