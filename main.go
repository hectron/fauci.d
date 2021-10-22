package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hectron/fauci.d/mapbox"
	"github.com/hectron/fauci.d/vaccines"
	"github.com/slack-go/slack"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	mapboxClient   mapbox.Client
	vaccinesClient vaccines.Client
	slackClient    *slack.Client
	lambdaInvoked  bool
)

func init() {
	mapboxClient = mapbox.Client{
		ApiToken: os.Getenv("MAPBOX_API_TOKEN"),
		ApiUrl:   os.Getenv("MAPBOX_API_URL"),
	}
	slackClient = slack.New(os.Getenv("SLACK_API_TOKEN"))
	vaccinesClient = vaccines.Client{ApiUrl: os.Getenv("VACCINE_API_URL")}
	lambdaInvoked = os.Getenv("LAMBDA") == "true"
}

func main() {
	lambda.Start(SimpleHandler)
}

func SimpleHandler(ctx context.Context) {
	fmt.Println("This is a line")
}
