package main

import (
	"os"

	"github.com/hectron/fauci.d/mapbox"
	"github.com/hectron/fauci.d/vaccines"
	"github.com/slack-go/slack"
)

var (
	mapboxClient   *mapbox.Client
	slackClient    *slack.Client
	vaccinesClient *vaccines.Client
	lambdaInvoked  bool
)

func init() {
	mapboxClient = mapbox.Client{
		ApiToken: os.Getenv("MAPBOX_API_TOKEN"),
		ApiUrl:   os.Getenv("MAPBOX_API_URL"),
	}
	slackClient = slack.New(os.Getenv("SLACK_API_TOKEN"))
	vaccinesClient = vaccines.Client{ApiUrl: os.Getenv("VACCINE_API_URL")}
	lambdaInvoked = os.Getenv("LAMBDA") != ""
}

func main() {
	if lambdaInvoked {
		// handle the appropriate request
	} else {
		// spin up http server to listen to requests
		// create servemux for handling routes
		// delegate accordingly
	}
}
