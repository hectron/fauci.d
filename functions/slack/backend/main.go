package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hectron/fauci.d/mapbox"
	"github.com/hectron/fauci.d/slack"
	"github.com/hectron/fauci.d/vaccines"

	slackGo "github.com/slack-go/slack"
)

type SlackRequest struct {
	ChannelId  string           `json:"channelId"`
	PostalCode string           `json:"postalCode"`
	Vaccine    vaccines.Vaccine `json:"vaccine"`
}

var (
	mapboxClient   mapbox.Client
	vaccinesClient vaccines.Client
	slackClient    *slackGo.Client
)

func init() {
	vaccinesClient = vaccines.Client{ApiUrl: os.Getenv("VACCINE_API_URL")}
	mapboxClient = mapbox.Client{
		ApiToken: os.Getenv("MAPBOX_API_TOKEN"),
		ApiUrl:   os.Getenv("MAPBOX_API_URL"),
	}
	slackClient = slackGo.New(os.Getenv("SLACK_API_TOKEN"))
}

func main() {
	fmt.Println("Starting the lambda")
	lambda.Start(MessageHandler)
	fmt.Println("Done with the lambda")
}

func MessageHandler(ctx context.Context, request SlackRequest) (events.APIGatewayProxyResponse, error) {
	notifyStartToUser(request)

	fmt.Printf("Looking up coordinates for %s", request.PostalCode)

	coordinates, err := mapboxClient.GeocodePostalCode(request.PostalCode)

	if err != nil {
		return failAndNotifyInSlack(
			err,
			"Something went wrong with the request. Please try again later.",
			request.ChannelId,
		)
	}

	req := vaccines.ApiRequest{Vaccine: request.Vaccine, Lat: coordinates.Latitude, Long: coordinates.Longitude}

	fmt.Printf("Looking up %s vaccine appointments", request.Vaccine.String())
	providers, err := vaccinesClient.FindVaccines(req)

	if err != nil {
		return failAndNotifyInSlack(
			err,
			"Unable to find vaccination appointments at the moment. Please try again later.",
			request.ChannelId,
		)
	}

	fmt.Println("Posting to Slack")
	blocks := slack.BuildBlocksForProviders(request.PostalCode, request.Vaccine.String(), providers)
	slackClient.PostMessage(request.ChannelId, slackGo.MsgOptionBlocks(blocks...))

	fmt.Printf("Successfully sent a message to %s", request.ChannelId)

	return events.APIGatewayProxyResponse{Body: "", StatusCode: 200}, nil
}

func failAndNotifyInSlack(err error, message string, channelId string) (events.APIGatewayProxyResponse, error) {
	fmt.Println(err)
	slackClient.PostMessage(channelId, slackGo.MsgOptionText(message, false))
	return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, errors.New(message)
}

func notifyStartToUser(request SlackRequest) {
	msg := fmt.Sprintf("Looking up appointments for the %s vaccine in %s... :eyes:", request.Vaccine.String(), request.PostalCode)
	fmt.Println(msg)

	slackClient.PostMessage(request.ChannelId, slackGo.MsgOptionText(msg, true))
}
