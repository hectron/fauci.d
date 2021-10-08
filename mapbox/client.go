package mapbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type Client struct {
	ApiUrl, ApiToken string
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

func (c *Client) GeocodePostalCode(postalCode string) (Coordinates, error) {
	var (
		coordinates Coordinates
		httpClient  *http.Client
		httpRequest *http.Request
		err         error
	)

	httpClient = &http.Client{}
	httpRequest, err = c.buildHttpRequest(postalCode)

	if err != nil {
		return coordinates, err
	}

	response, err := httpClient.Do(httpRequest)

	if err != nil {
		return coordinates, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return coordinates, errors.New(fmt.Sprintf("Mapbox API response: %s", response.Status))
	}

	return c.parseCoordinates(response.Body)
}

func (c *Client) buildHttpRequest(postalCode string) (*http.Request, error) {
	apiUrl, err := url.Parse(c.ApiUrl)

	if err != nil {
		return nil, err
	}

	apiUrl.Path = path.Join(apiUrl.Path, fmt.Sprintf("%s.json", postalCode))

	request, err := http.NewRequest(http.MethodGet, apiUrl.String(), nil)

	if err != nil {
		return request, err
	}

	q := request.URL.Query()
	q.Set("country", "us")
	q.Set("types", "postcode")
	q.Set("access_token", c.ApiToken)

	request.URL.RawQuery = q.Encode()

	return request, nil
}

func (c *Client) parseCoordinates(r io.Reader) (Coordinates, error) {
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
