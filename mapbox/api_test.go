package mapbox

import (
	"fmt"
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

	t.Run("It sends over the params", func(t *testing.T) {
		expectedCountry := "us"
		expectedTypes := "postcode"
		expectedAccessToken := "totally-valid-token"

		serverMux := http.NewServeMux()
		serverMux.HandleFunc("/60601.json", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("access_token") != expectedAccessToken {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if r.URL.Query().Get("country") != expectedCountry || r.Header.Get("types") != expectedTypes {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			fmt.Fprintf(w, `{ "some": "json" }`)
		})

		mock_http_server := httptest.NewServer(serverMux)
		defer mock_http_server.Close()

		// test the bad case
		client := Client{ApiUrl: mock_http_server.URL, ApiToken: "allegedly-valid-token"}

		_, err := client.GeocodePostalCode("60601")

		want := "Mapbox API response: 401 Unauthorized"
		got := err.Error()

		if got != want {
			t.Errorf("Expected %s, got %s", want, got)
		}
	})
}
