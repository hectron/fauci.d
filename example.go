package main

import (
	"fmt"
	"os"

	"github.com/hectron/fauci.d/mapbox"
	"github.com/hectron/fauci.d/vaccines"
	"github.com/slack-go/slack"
)

func main() {
	mapboxApiToken := os.Getenv("MAPBOX_API_TOKEN")
	slackApiToken := os.Getenv("SLACK_API_TOKEN")
	mapboxApiUrl := "https://api.mapbox.com/geocoding/v5/mapbox.places"
	vaccineApiUrl := "https://api.us.castlighthealth.com/vaccine-finder/v1/provider-locations/search"

	mapboxClient := mapbox.Client{ApiToken: mapboxApiToken, ApiUrl: mapboxApiUrl}
	coordinates, err := mapboxClient.GeocodePostalCode("60640")

	if err != nil {
		fmt.Println(err)
		return
	}

	vaccineClient := vaccines.Client{ApiUrl: vaccineApiUrl}

	req := vaccines.ApiRequest{
		Vaccine: vaccines.Moderna,
		Lat:     coordinates.Latitude,
		Long:    coordinates.Longitude,
	}

	providers, err := vaccineClient.FindVaccines(req)

	if err != nil {
		fmt.Println("Could not load response")
		fmt.Println(err)
		return
	}

	channelId := "CUP3PES12"
	blocks := FormatForSlackUsingBlocks(providers)
	slackApi := slack.New(slackApiToken)
	slackApi.PostMessage(channelId, slack.MsgOptionBlocks(blocks...))

	fmt.Printf("Succesful message sent to channel %s", channelId)
}
