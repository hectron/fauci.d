package mapbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
)

var (
	url string
)

func init() {
	url = "https://api.mapbox.com/geocoding/v5/mapbox.places"
}

type Api struct {
	Token string
}

type Coordinates struct {
	Latitude, Longitude float64
}

type apiResponse struct {
	Features []struct {
		PlaceType []string  `json:"place_type"`
		Center    []float64 `json:"center"`
	}
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (a *Api) GeocodePostalCode(postalCode string, client HTTPClient) (Coordinates, error) {
	var coordinates Coordinates

	request, err := a.buildRequest(postalCode)

	if err != nil {
		return coordinates, err
	}

	response, err := client.Do(request)

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return coordinates, errors.New(fmt.Sprintf("Mapbox API response: %s", response.Status))
	}

	return a.parseCoordinates(response.Body)
}

func (a *Api) buildRequest(postalCode string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return request, err
	}

	endpoint := fmt.Sprintf("%s.json", postalCode)
	request.URL.Path = path.Join(request.URL.Path, endpoint)

	q := request.URL.Query()
	q.Set("country", "us")
	q.Set("types", "postcode")
	q.Set("access_token", a.Token)

	request.URL.RawQuery = q.Encode()

	return request, nil
}

func (a *Api) parseCoordinates(r io.Reader) (Coordinates, error) {
	var (
		resp        apiResponse
		coordinates Coordinates
	)

	body, err := io.ReadAll(r)

	if err != nil {
		return coordinates, err
	}

	json.Unmarshal(body, &resp)

	for _, feature := range resp.Features {
		for _, placeType := range feature.PlaceType {
			if placeType == "postcode" {
				coordinates.Latitude = feature.Center[1]
				coordinates.Longitude = feature.Center[0]

				break
			}
		}

		if coordinates.Latitude != 0 && coordinates.Longitude != 0 {
			break
		}
	}

	return coordinates, nil
}
