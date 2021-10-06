package vaccines

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Api struct{}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type VaccineProvider struct {
	Guid, Name                                      string
	Address1, Address2, City, State, Zipcode, Phone string
	Distance, Lat, Long                             float64
	AcceptsWalkIns, AppointmentsAvailable, InStock  bool
}

type ApiRequest struct {
	Vaccine   Vaccine
	Lat, Long float64
}

type apiResponse struct {
	Providers []VaccineProvider
}

var (
	url string
)

func init() {
	url = "https://api.us.castlighthealth.com/vaccine-finder/v1/provider-locations/search"
}

func (a *Api) Request(request ApiRequest, httpClient HTTPClient) ([]VaccineProvider, error) {
	httpReq, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	setQueryString(httpReq, &request)
	setRequestHeaders(httpReq)

	response, err := httpClient.Do(httpReq)
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Vaccines API returned with status: %s", response.Status))
	}

	return a.parseProviders(response.Body)
}

func (a *Api) parseProviders(r io.Reader) ([]VaccineProvider, error) {
	var resp apiResponse

	body, err := io.ReadAll(r)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(body, &resp)

	return resp.Providers, nil
}

func setQueryString(request *http.Request, apiRequest *ApiRequest) {
	query := request.URL.Query()

	query.Set("medicationGuids", apiRequest.Vaccine.Guid())
	query.Set("long", fmt.Sprintf("%f", apiRequest.Long))
	query.Set("lat", fmt.Sprintf("%f", apiRequest.Lat))
	query.Set("appointments", "true")
	query.Set("radius", "5")

	request.URL.RawQuery = query.Encode()
}

func setRequestHeaders(request *http.Request) {
	request.Header.Add("Accept-Language", "en-US,en;q=0.9")
	request.Header.Add("Accept", "application/json, text/plain, */*")
	request.Header.Add("User-Agent", "Mozilla/5.0")
}
