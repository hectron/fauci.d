package mapbox

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeocodePostalCode(t *testing.T) {
	t.Run("It returns empty coordinates when the status code is not 200", func(t *testing.T) {
		mock_http_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))

		defer mock_http_server.Close()

		client := Client{ApiUrl: mock_http_server.URL, ApiToken: "mocktoken"}
		coordinates, err := client.GeocodePostalCode("60601")

		if coordinates.Latitude != 0 || coordinates.Longitude != 0 {
			t.Errorf(
				"Expected coordinates to be empty. Got: latitude - %f, longitude - %f",
				coordinates.Latitude,
				coordinates.Longitude,
			)
		}

		if err == nil {
			t.Errorf("Expected %s, got no error", err)
		}
	})

	t.Run("sending params", func(t *testing.T) {
		expectedCountry := "us"
		expectedTypes := "postcode"
		expectedAccessToken := "totally-valid-token"

		mockResponse, err := ioutil.ReadFile("fixtures/mapbox_response.json")

		if err != nil {
			t.Errorf("Could not load fixture for test")
		}

		serverMux := http.NewServeMux()
		serverMux.HandleFunc("/60601.json", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("access_token") != expectedAccessToken {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if r.URL.Query().Get("country") != expectedCountry || r.URL.Query().Get("types") != expectedTypes {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			fmt.Fprintf(w, string(mockResponse))
		})

		mock_http_server := httptest.NewServer(serverMux)
		defer mock_http_server.Close()

		client := Client{ApiUrl: mock_http_server.URL, ApiToken: expectedAccessToken}
		got, err := client.GeocodePostalCode("60601")

		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		if got.Latitude == 0 || got.Longitude == 0 {
			t.Errorf(
				"Expected coordinates to be empty. Got: latitude - %f, longitude - %f",
				got.Latitude,
				got.Longitude,
			)
		}
	})
}
