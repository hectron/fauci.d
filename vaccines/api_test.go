package vaccines

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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
			query := r.URL.Query()
			headers := r.Header

			requiredParams := []string{
				"lat",
				"long",
				"appointments",
				"radius",
				"medicationGuids",
			}

			requiredHeaders := []string{
				"Accept",
				"Accept-Language",
				"User-Agent",
			}

			missingParam := false
			missingHeader := false

			for _, param := range requiredParams {
				if query.Get(param) == "" {
					fmt.Println("missing " + param)
					missingParam = true
					break
				}
			}

			if missingParam {
				fmt.Println("Missing params")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			for _, header := range requiredHeaders {
				if headers.Get(header) == "" {
					missingHeader = true
					break
				}
			}

			if missingHeader {
				fmt.Println("Missing headers")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			validGuids := []string{
				Pfizer.Guid(),
				Moderna.Guid(),
				JJ.Guid(),
			}

			hasGuid := false

			for _, guid := range validGuids {
				if strings.Contains(query.Get("medicationGuids"), guid) {
					hasGuid = true
					break
				}
			}

			if !hasGuid {
				fmt.Println("Missing guid")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

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

		expectedProviders := []VaccineProvider{
			{
				Guid:                  "5572f886-7105-4267-84e1-cfa8aa052608",
				Name:                  "MICHIGAN AVENUE PRIMARY CARE / IMMEDIATE CARE",
				Address1:              "180 N Michigan Ave #1605",
				Address2:              "",
				City:                  "Chicago",
				State:                 "IL",
				Zipcode:               "60601",
				Phone:                 "",
				Distance:              0.08,
				Lat:                   41.885511,
				Long:                  -87.624889,
				AcceptsWalkIns:        false,
				AppointmentsAvailable: false,
				InStock:               true,
			}, {
				Guid:                  "46bf1df3-6215-44b9-b1ed-b9863b6825a6",
				Name:                  "CVS Pharmacy, Inc. #04781",
				Address1:              "205 N Michigan Ave",
				Address2:              "",
				City:                  "Chicago",
				State:                 "IL",
				Zipcode:               "60601",
				Phone:                 "(312) 938-4091",
				Distance:              0.09,
				Lat:                   41.886098,
				Long:                  -87.624033,
				AcceptsWalkIns:        true,
				AppointmentsAvailable: true,
				InStock:               true,
			}, {
				Guid:                  "ed3f9095-3156-4f8c-86bc-46f0894497bc",
				Name:                  "Walgreens Co. #9438",
				Address1:              "30 N Michigan Ave LBBY 1",
				Address2:              "",
				City:                  "Chicago",
				State:                 "IL",
				Zipcode:               "60602",
				Phone:                 "312-332-3540",
				Distance:              0.14,
				Lat:                   41.883004,
				Long:                  -87.624772,
				AcceptsWalkIns:        true,
				AppointmentsAvailable: false,
				InStock:               true,
			}, {
				Guid:                  "aadf16d6-5525-46a6-9f8c-4ada508107bb",
				Name:                  "Walgreens Co. #15196",
				Address1:              "Flair Tower, 151 N State St FL 1ST",
				Address2:              "",
				City:                  "Chicago",
				State:                 "IL",
				Zipcode:               "60601",
				Phone:                 "312-863-4249",
				Distance:              0.21,
				Lat:                   41.884799,
				Long:                  -87.627623,
				AcceptsWalkIns:        true,
				AppointmentsAvailable: false,
				InStock:               true,
			}, {
				Guid:                  "6e0393cd-8040-4f47-945e-9d55b6a4ea98",
				Name:                  "CVS Pharmacy, Inc. #08910",
				Address1:              "205 N Columbus Dr",
				Address2:              "",
				City:                  "Chicago",
				State:                 "IL",
				Zipcode:               "60611",
				Phone:                 "(312) 861-0315",
				Distance:              0.22,
				Lat:                   41.886222,
				Long:                  -87.619717,
				AcceptsWalkIns:        true,
				AppointmentsAvailable: true,
				InStock:               true,
			},
		}

		for _, want := range expectedProviders {
			var got VaccineProvider
			emptyProvider := VaccineProvider{}

			for _, p := range providers {
				if p.Guid == want.Guid {
					got = p
					break
				}
			}

			if got != emptyProvider {
				if !reflect.DeepEqual(got, want) {
					t.Errorf("got %v, want %v", got, want)
				}
			} else {
				t.Errorf("want %v, but didn't get anything", want)
			}
		}
	})
}
