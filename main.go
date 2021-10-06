package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/hectron/fauci.d/mapbox"
	"github.com/hectron/fauci.d/vaccines"
)

func main() {
	mapboxToken := os.Getenv("MAPBOX_API_TOKEN")
	mapboxApi := mapbox.Api{Token: mapboxToken}
	httpClient := &http.Client{}
	coordinates, err := mapboxApi.GeocodePostalCode("60640", httpClient)

	if err != nil {
		fmt.Println(err)
		return
	}

	vaccineApi := vaccines.Api{}
	client := &http.Client{}

	apiRequest := vaccines.ApiRequest{
		Vaccine: vaccines.Moderna,
		Lat:     coordinates.Latitude,
		Long:    coordinates.Longitude,
	}

	providers, err := vaccineApi.Request(apiRequest, client)

	if err == nil {
		fmt.Println(providers)
	} else {
		fmt.Println("Could not load response")
		fmt.Println(err)
	}
}
