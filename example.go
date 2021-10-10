package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hectron/fauci.d/mapbox"
	"github.com/hectron/fauci.d/vaccines"
	"github.com/slack-go/slack"
)

func FormatForSlackUsingBlocks(providers []vaccines.VaccineProvider) []slack.Block {
	blocks := []slack.Block{}
	divSection := slack.NewDividerBlock()

	// header section
	headerText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Found %d providers near you!", len(providers)), false, false)
	blocks = append(blocks, slack.NewSectionBlock(headerText, nil, nil))

	for idx, provider := range providers {
		if idx > 10 {
			break
		}

		text := slack.NewTextBlockObject("mrkdwn", ProviderAsString(provider), false, false)
		section := slack.NewSectionBlock(text, nil, nil)
		blocks = append(blocks, section, divSection)
	}

	return blocks
}

func ProviderAsString(provider vaccines.VaccineProvider) string {
	return fmt.Sprintf(
		"<%s|*%s*> located at %s, %s, %s %s (about %s miles away). Phone Number: %s",
		provider.Website(),
		provider.Name,
		provider.Address1,
		provider.City,
		provider.State,
		provider.Zipcode,
		strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", provider.Distance), "0"), "."),
		provider.Phone,
	)
}

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
