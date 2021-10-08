package vaccines

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// url = "https://api.us.castlighthealth.com/vaccine-finder/v1/provider-locations/search"

type VaccineProvider struct {
	Guid, Name                             string
	Address1, Address2, City, State, Phone string
	Zipcode                                string `json:"zip"`
	Distance, Lat, Long                    float64
	AcceptsWalkIns                         bool                    `json:"accepts_walk_ins"`
	AppointmentsAvailable                  appointmentAvailability `json:"appointments_available"`
	InStock                                bool                    `json:"in_stock"`
}

func (v *VaccineProvider) Website() string {
	return fmt.Sprintf("https://www.vaccines.gov/provider/?id=%s", v.Guid)
}

type appointmentAvailability bool

func (a *appointmentAvailability) UnmarshalJSON(data []byte) error {
	var s string

	err := json.Unmarshal(data, &s)

	if err != nil {
		return err
	}

	*a = s == "TRUE"

	return nil
}

type ApiRequest struct {
	Vaccine   Vaccine
	Lat, Long float64
}

type apiResponse struct {
	Providers []VaccineProvider
}

type Client struct {
	ApiUrl string
}

func (c *Client) FindVaccines(request ApiRequest) ([]VaccineProvider, error) {
	var (
		httpClient *http.Client
		httpReq    *http.Request
		response   *http.Response
		err        error
	)
	httpClient = &http.Client{}
	httpReq, err = http.NewRequest("GET", c.ApiUrl, nil)

	if err != nil {
		return nil, err
	}

	setQueryString(httpReq, &request)
	setRequestHeaders(httpReq)

	response, err = httpClient.Do(httpReq)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Vaccines API returned with status: %s", response.Status))
	}

	defer response.Body.Close()

	return c.parseProviders(response.Body)
}

func (c *Client) parseProviders(r io.Reader) ([]VaccineProvider, error) {
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
