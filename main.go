package main

import (
	"fmt"
	"os"

	"github.com/hectron/fauci.d/mapbox"
	"github.com/hectron/fauci.d/vaccines"
)

func main() {
	mapboxToken := os.Getenv("MAPBOX_API_TOKEN")
	mapboxClient := mapbox.Client{ApiToken: mapboxToken, ApiUrl: "https://api.mapbox.com/geocoding/v5/mapbox.places"}
	coordinates, err := mapboxClient.GeocodePostalCode("60640")

	if err != nil {
		fmt.Println(err)
		return
	}

	vaccineClient := vaccines.Client{ApiUrl: "https://api.us.castlighthealth.com/vaccine-finder/v1/provider-locations/search"}

	req := vaccines.ApiRequest{
		Vaccine: vaccines.Moderna,
		Lat:     coordinates.Latitude,
		Long:    coordinates.Longitude,
	}

	providers, err := vaccineClient.FindVaccines(req)

	if err == nil {
		fmt.Println(providers)
	} else {
		fmt.Println("Could not load response")
		fmt.Println(err)
	}
}
