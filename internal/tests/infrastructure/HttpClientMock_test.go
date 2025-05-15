package infrastructure_test

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type HttpClientMock struct {
	mock.Mock
}

func (m *HttpClientMock) MakeGet(url string) (*http.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}
