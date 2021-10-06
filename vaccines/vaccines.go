package vaccines

import (
	"encoding/json"
	"io"
	"net/http"
)

type Vaccine int

const (
	Pfizer Vaccine = iota
	Moderna
	JJ
)

const apiUrl = "https://api.us.castlighthealth.com/vaccine-finder/v1/provider-locations/search"

var vaccineToGuid = map[Vaccine]string{
	Pfizer:  "a84fb9ed-deb4-461c-b785-e17c782ef88b",
	Moderna: "779bfe52-0dd8-4023-a183-457eb100fccc",
	JJ:      "784db609-dc1f-45a5-bad6-8db02e79d44f",
}
var httpClient = &http.Client{}

type VaccineProvider struct {
	Guid, Name                                      string
	Address1, Address2, City, State, Zipcode, Phone string
	Distance, Lat, Long                             float64
	AcceptsWalkIns, AppointmentsAvailable, InStock  bool
}

type ApiResponse struct {
	Providers []VaccineProvider
}

func (v Vaccine) Guid() string {
	return vaccineToGuid[v]
}

type Api struct {
	Vaccine Vaccine
}

func (a *Api) Request() ([]VaccineProvider, error) {
	request, err := http.NewRequest("GET", apiUrl, nil)

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
	var apiResponse ApiResponse

	body, err := io.ReadAll(r)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(body, &apiResponse)

	return apiResponse.Providers, nil
}
