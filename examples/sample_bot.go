package main

import (
	"fmt"
	"os"

	"github.com/hectron/fauci.d/mapbox"
	"github.com/hectron/fauci.d/slack"
	"github.com/hectron/fauci.d/vaccines"
	slackGo "github.com/slack-go/slack"
)

func exampleMain() {
	mapboxApiToken := os.Getenv("MAPBOX_API_TOKEN")
	slackApiToken := os.Getenv("SLACK_API_TOKEN")
	mapboxApiUrl := "https://api.mapbox.com/geocoding/v5/mapbox.places"
	vaccineApiUrl := "https://api.us.castlighthealth.com/vaccine-finder/v1/provider-locations/search"
	postalCode := "60640"

	mapboxClient := mapbox.Client{ApiToken: mapboxApiToken, ApiUrl: mapboxApiUrl}
	coordinates, err := mapboxClient.GeocodePostalCode(postalCode)

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
	blocks := slack.BuildBlocksForProviders(postalCode, "moderna", providers)
	slackApi := slackGo.New(slackApiToken)
	slackApi.PostMessage(channelId, slackGo.MsgOptionBlocks(blocks...))

	fmt.Printf("Succesful message sent to channel %s", channelId)
}
