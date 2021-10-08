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

	slackApi := slack.New(slackApiToken)

	attachment := slack.Attachment{
		Pretext: "This is the pretext of an attechment",
		Text:    fmt.Sprintf("Found %d appointments for %s in %s", len(providers), "moderna", "60640"),
	}

	channelId, timestamp, err := slackApi.PostMessage(
		"CUP3PES12",
		slack.MsgOptionText("Optional text", false),
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	fmt.Printf("Succesful message sent to channel %s at %s", channelId, timestamp)
}
