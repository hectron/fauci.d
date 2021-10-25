package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/hectron/fauci.d/mapbox"
	"github.com/hectron/fauci.d/slack"
	"github.com/hectron/fauci.d/vaccines"
	"github.com/pkg/errors"
	slackGo "github.com/slack-go/slack"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	mapboxClient   mapbox.Client
	vaccinesClient vaccines.Client
	slackClient    *slackGo.Client
	lambdaInvoked  bool
)

func init() {
	mapboxClient = mapbox.Client{
		ApiToken: os.Getenv("MAPBOX_API_TOKEN"),
		ApiUrl:   os.Getenv("MAPBOX_API_URL"),
	}
	slackClient = slackGo.New(os.Getenv("SLACK_API_TOKEN"))
	vaccinesClient = vaccines.Client{ApiUrl: os.Getenv("VACCINE_API_URL")}
	lambdaInvoked = os.Getenv("LAMBDA") == "true"
}

func main() {
	lambda.Start(SimpleHandler)
}

func SimpleHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	m, err := url.ParseQuery(request.Body)

	if err != nil {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, err
	}

	jsonBody, err := json.Marshal(m)

	if err != nil {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, err
	}

	fmt.Printf("=== Request: %s\n", request.Body)
	fmt.Printf("=== json body: %s\n", string(jsonBody))

	postalCode := m.Get("text")
	channelId := m.Get("channel_id")
	vaccineCommand := m.Get("command")

	if postalCode == "" {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, errors.New("No postal code supplied")
	}

	if channelId == "" {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, errors.New("Could not determine channel to post to")
	}

	fmt.Printf("=== Requested postal code `%s` in channel id `%s`", postalCode, channelId)
	coordinates, err := mapboxClient.GeocodePostalCode(postalCode)

	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, errors.New("Could not geocode the postal code")
	}

	var vaccine vaccines.Vaccine

	if vaccineCommand == "/pfizer" {
		vaccine = vaccines.Pfizer
	} else if vaccineCommand == "/moderna" {
		vaccine = vaccines.Moderna
	} else if vaccineCommand == "/jj" {
		vaccine = vaccines.JJ
	}

	req := vaccines.ApiRequest{
		Vaccine: vaccine,
		Lat:     coordinates.Latitude,
		Long:    coordinates.Longitude,
	}

	providers, err := vaccinesClient.FindVaccines(req)

	if err != nil {
		fmt.Println("Could not load response")
		fmt.Println(err)
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, errors.New("Unable to retrieve providers")
	}

	blocks := slack.BuildBlocksForProviders(postalCode, vaccine.String(), providers)
	slackClient.PostMessage(channelId, slackGo.MsgOptionBlocks(blocks...))

	return events.APIGatewayProxyResponse{Body: string(jsonBody), StatusCode: 200}, nil
}
