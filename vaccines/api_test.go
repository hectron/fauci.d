package vaccines

import (
	"net/http"
	"testing"
)

type MockClient struct {
}

func (m *MockClient) Do(request *http.Request) (*http.Response, error) {
	return &http.Response{}, nil
}

func TestApiRequest(t *testing.T) {
	// api := Api{Moderna}

	// mockClient := &MockClient{}
	// HttpClient = mockClient
	// _, err := api.Request()

	// if err != nil {
	// 	t.Errorf("Received an error making HTTP request")
	// }
}
