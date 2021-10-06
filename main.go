package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/hectron/fauci.d/mapbox"
)

func main() {
	mapboxToken := os.Getenv("MAPBOX_API_TOKEN")
	mapboxApi := mapbox.Api{Token: mapboxToken}
	httpClient := &http.Client{}
	coordinates, err := mapboxApi.GeocodePostalCode("60640", httpClient)

	if err == nil {
		fmt.Println(coordinates)
	} else {
		fmt.Println(err)
	}

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
}
