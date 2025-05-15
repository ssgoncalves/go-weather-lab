package infrastructure

import "net/http"

type HttpClientInterface interface {
	MakeGet(url string) (*http.Response, error)
}

type HttpClient struct {
	client *http.Client
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		client: &http.Client{},
	}
}

func (c *HttpClient) MakeGet(url string) (*http.Response, error) {
	return c.client.Get(url)
}
