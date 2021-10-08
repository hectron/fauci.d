package vaccines

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApiRequest(t *testing.T) {
	t.Run("A bad request returns an error and no providers", func(t *testing.T) {
		mock_http_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))

		client := Client{ApiUrl: mock_http_server.URL}
		request := ApiRequest{Moderna, 41.8848, -87.6235}
		providers, err := client.FindVaccines(request)

		if len(providers) > 0 {
			t.Errorf("Expected no providers. Received: %d", len(providers))
		}

		want := "Vaccines API returned with status: 400 Bad Request"
		got := err.Error()

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("A successful request returns providers", func(t *testing.T) {
		mockResponse, err := ioutil.ReadFile("fixtures/api_response.json")

		if err != nil {
			t.Errorf("Could not load fixture for test")
		}

		mock_http_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, string(mockResponse))
		}))

		client := Client{ApiUrl: mock_http_server.URL}
		request := ApiRequest{Moderna, 41.8848, -87.6235}
		providers, err := client.FindVaccines(request)

		if err != nil {
			t.Errorf("Did not expect error. Got: %s", err)
		}

		if len(providers) == 0 {
			t.Errorf("Expected providers, but did not receive any")
		}
	})
}
