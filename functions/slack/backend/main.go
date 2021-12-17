package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/getsentry/sentry-go"
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
	sentryEnabled  bool
	functionName   string
)

func init() {
	vaccinesClient = vaccines.Client{ApiUrl: os.Getenv("VACCINE_API_URL")}
	mapboxClient = mapbox.Client{
		ApiToken: os.Getenv("MAPBOX_API_TOKEN"),
		ApiUrl:   os.Getenv("MAPBOX_API_URL"),
	}
	slackClient = slackGo.New(os.Getenv("SLACK_API_TOKEN"))
	sentryEnabled = os.Getenv("SENTRY_DSN") != ""
	functionName = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")

	if sentryEnabled {
		sentry.Init(sentry.ClientOptions{
			Dsn:         os.Getenv("SENTRY_DSN"),
			Debug:       true,
			DebugWriter: os.Stderr,
			Environment: os.Getenv("SENTRY_ENVIRONMENT"),
			Release:     os.Getenv("SENTRY_RELEASE"),
		})

	}
}

func main() {
	if sentryEnabled {
		lambda.Start(withSentry(MessageHandler))
	} else {
		lambda.Start(MessageHandler)
	}
}

func MessageHandler(ctx context.Context, request SlackRequest) (events.APIGatewayProxyResponse, error) {
	notifyStartToUser(request)

	log.Printf("Looking up coordinates for %s\n", request.PostalCode)

	coordinates, err := mapboxClient.GeocodePostalCode(request.PostalCode)

	if err != nil {
		return failAndNotifyInSlack(
			err,
			"Something went wrong with the request. Please try again later.",
			request.ChannelId,
		)
	}

	req := vaccines.ApiRequest{Vaccine: request.Vaccine, Lat: coordinates.Latitude, Long: coordinates.Longitude}

	log.Printf("Looking up %s vaccine appointments\n", request.Vaccine.String())
	providers, err := vaccinesClient.FindVaccines(req)

	if err != nil {
		return failAndNotifyInSlack(
			err,
			"Unable to find vaccination appointments at the moment. Please try again later.",
			request.ChannelId,
		)
	}

	blocks := slack.BuildBlocksForProviders(request.PostalCode, request.Vaccine.String(), providers)
	slackClient.PostMessage(request.ChannelId, slackGo.MsgOptionBlocks(blocks...))

	log.Printf("Successfully sent a message to %s\n", request.ChannelId)

	return events.APIGatewayProxyResponse{Body: "", StatusCode: 200}, nil
}

func failAndNotifyInSlack(err error, message string, channelId string) (events.APIGatewayProxyResponse, error) {
	log.Println(err)
	slackClient.PostMessage(channelId, slackGo.MsgOptionText(message, false))
	return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, errors.New(message)
}

func notifyStartToUser(request SlackRequest) {
	msg := fmt.Sprintf("Looking up appointments for the %s vaccine in %s... :eyes:", request.Vaccine.String(), request.PostalCode)
	log.Println(msg)

	slackClient.PostMessage(request.ChannelId, slackGo.MsgOptionText(msg, true))
}

func withSentry(f func(context.Context, SlackRequest) (events.APIGatewayProxyResponse, error)) func(context.Context, SlackRequest) (events.APIGatewayProxyResponse, error) {
	function := f

	return func(ctx context.Context, s SlackRequest) (events.APIGatewayProxyResponse, error) {
		log.Println("Starting invocation")

		sentryHub := sentry.CurrentHub().Clone()

		sentryHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("function", functionName)
		})

		defer sentryHub.Flush(time.Second * 2)

		resp, err := function(ctx, s)

		if err != nil {
			log.Printf("Something went wrong! %s\n", err.Error())
			sentryHub.CaptureException(err)
		}

		log.Println("Finished invocation")

		return resp, err
	}
}
