package main

import (
	"fmt"
	"net/http"

	"github.com/hectron/fauci.d/vaccines"
)

func main() {
	//mapboxToken := os.Getenv("MAPBOX_API_TOKEN")
	//mapboxApi := mapbox.Api{Token: mapboxToken}
	//httpClient := &http.Client{}
	//coordinates, err := mapboxApi.GeocodePostalCode("60640", httpClient)

	//if err == nil {
	//	fmt.Println(coordinates)
	//} else {
	//	fmt.Println(err)
	//}

	// vaccineApi := vaccines.Api{Vaccine: vaccines.Moderna}
	// client := &http.Client{}
	//
	// providers, err := vaccineApi.Request(client)
	//
	// if err == nil {
	// 	fmt.Println(providers)
	// } else {
	// 	fmt.Println("Could not load response")
	// 	fmt.Println(err)
	// }

	vaccineApi := vaccines.Api{}
	client := &http.Client{}

	apiRequest := vaccines.ApiRequest{
		Vaccine: vaccines.Moderna,
		Lat:     41.97,
		Long:    -87.66,
	}

	providers, err := vaccineApi.Request(apiRequest, client)

	if err == nil {
		fmt.Println(providers)
	} else {
		fmt.Println("Could not load response")
		fmt.Println(err)
	}
}
