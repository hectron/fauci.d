package vaccines

import (
	"encoding/json"
	"io"
	"net/http"
)

type Api struct {
	Vaccine Vaccine
}

type VaccineProvider struct {
	Guid, Name                                      string
	Address1, Address2, City, State, Zipcode, Phone string
	Distance, Lat, Long                             float64
	AcceptsWalkIns, AppointmentsAvailable, InStock  bool
}

type apiResponse struct {
	Providers []VaccineProvider
}

const (
	url = "https://api.us.castlighthealth.com/vaccine-finder/v1/provider-locations/search"
)

func (a *Api) Request(httpClient *http.Client) ([]VaccineProvider, error) {
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	query := request.URL.Query()
	query.Set("medicationGuids", a.Vaccine.Guid())
	query.Set("long", "-87.7025")
	query.Set("lat", "41.9215")
	query.Set("appointments", "true")
	query.Set("radius", "5")

	request.URL.RawQuery = query.Encode()

	request.Header.Add("Accept-Language", "en-US,en;q=0.9")
	request.Header.Add("Accept", "application/json, text/plain, */*")
	request.Header.Add("User-Agent", "Mozilla/5.0")

	response, err := httpClient.Do(request)
	defer response.Body.Close()

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
