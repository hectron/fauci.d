package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/hectron/fauci.d/mapbox"
	"github.com/hectron/fauci.d/vaccines"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"

	"github.com/aws/aws-lambda-go/events"
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

func SimpleHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("This is a line")
	fmt.Println("There are the relevant environment variables:")

	envVars := []string{"LAMBDA", "COMMAND"}

	for _, e := range envVars {
		fmt.Printf("=== %s\n", e)
		fmt.Println(os.Getenv(e))
	}

	m, err := url.ParseQuery(request.Body)

	if err != nil {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, err
	}

	jsonBody, err := json.Marshal(m)

	if err != nil {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, err
	}

	fmt.Printf("Request: %s", request.Body)
	fmt.Printf("json body: %s", string(jsonBody))

	postalCode := m.Get("text")
	channelId := m.Get("channel_id")

	if postalCode == "" {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, errors.New("No postal code supplied")
	}

	if channelId == "" {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, errors.New("Could not determine channel to post to")
	}

	fmt.Printf("Requested postal code `%s` in channel id `%s`", postalCode, channelId)
	coordinates, err := mapboxClient.GeocodePostalCode(postalCode)

	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, errors.New("Could not geocode the postal code")
	}

	req := vaccines.ApiRequest{
		Vaccine: vaccines.Moderna,
		Lat:     coordinates.Latitude,
		Long:    coordinates.Longitude,
	}

	providers, err := vaccinesClient.FindVaccines(req)

	if err != nil {
		fmt.Println("Could not load response")
		fmt.Println(err)
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, errors.New("Unable to retrieve providers")
	}

	blocks := BuildSlackBlocksForProviders(postalCode, providers)
	slackClient.PostMessage(channelId, slack.MsgOptionBlocks(blocks...))

	return events.APIGatewayProxyResponse{Body: string(jsonBody), StatusCode: 200}, nil
}
