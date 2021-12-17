package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/hectron/fauci.d/vaccines"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaSdk "github.com/aws/aws-sdk-go/service/lambda"
)

var (
	slackClient                                              *slack.Client
	successfulAsyncStatusCode                                int64
	functionName, backendFunctionName, somethingWentWrongMsg string
	somethingWentWrongSlackMsg                               slack.MsgOption
	sentryEnabled                                            bool
)

type BackendRequest struct {
	ChannelId  string           `json:"channelId"`
	PostalCode string           `json:"postalCode"`
	Vaccine    vaccines.Vaccine `json:"vaccine"`
}

func init() {
	successfulAsyncStatusCode = 202
	somethingWentWrongMsg = "I'm sorry :( \nSomething went wrong and I'm unable to process request."
	somethingWentWrongSlackMsg = slack.MsgOptionText(somethingWentWrongMsg, false)

	functionName = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	backendFunctionName = fmt.Sprintf("%s_backend", functionName)
	slackClient = slack.New(os.Getenv("SLACK_API_TOKEN"))

	sentryEnabled = os.Getenv("SENTRY_DSN") != ""

	if sentryEnabled {
		sentry.Init(sentry.ClientOptions{
			Dsn:         os.Getenv("SENTRY_DSN"),
			Debug:       true,
			DebugWriter: os.Stderr,
			Environment: os.Getenv("SENTRY_ENVIRONMENT"),
			Release:     os.Getenv("SENTRY_RELEASE"),
		})

		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("function", functionName)
		})
	}
}

func main() {
	if sentryEnabled {
		log.Println("Sentry enabled and bootstrapped!")
		lambda.Start(withSentry(IncomingMessageHandler))
	} else {
		lambda.Start(IncomingMessageHandler)
	}
}

func IncomingMessageHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		m                     url.Values
		err                   error
		channelId, postalCode string
		vaccine               vaccines.Vaccine
	)

	m, err = url.ParseQuery(request.Body)

	if err != nil {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, err
	}

	log.Printf("=== Request: %s\n", request.Body)

	if channelId = m.Get("channel_id"); channelId == "" {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, errors.New("Could not determine channel to post to")
	}

	if postalCode = m.Get("text"); postalCode == "" {
		return failAndNotifyInSlack("No postal code supplied.", channelId)
	}

	log.Printf("=== Requested postal code `%s` in channel id `%s`\n", postalCode, channelId)

	switch vaccineCommand := m.Get("command"); vaccineCommand {
	case "/pfizer":
		vaccine = vaccines.Pfizer
	case "/moderna":
		vaccine = vaccines.Moderna
	case "/jj":
		vaccine = vaccines.JJ
	default:
		vaccine = vaccines.Null
	}

	if vaccine != vaccines.Null {
		invokeVaccineFinderLambda(channelId, postalCode, vaccine)
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 200}, nil
	}

	return failAndNotifyInSlack(somethingWentWrongMsg, channelId)
}

func invokeVaccineFinderLambda(channelId string, postalCode string, vaccine vaccines.Vaccine) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigStateFromEnv,
	}))

	lambdaClient := lambdaSdk.New(sess)

	payload, err := json.Marshal(BackendRequest{
		ChannelId:  channelId,
		PostalCode: postalCode,
		Vaccine:    vaccine,
	})

	if err != nil {
		msg := fmt.Sprintf("Could not generate request for backend lambda: %s", err)
		log.Print(msg)
		slackClient.PostMessage(channelId, slack.MsgOptionText(msg, false))
		return
	}

	result, err := lambdaClient.Invoke(&lambdaSdk.InvokeInput{
		FunctionName:   aws.String(backendFunctionName),
		Payload:        payload,
		InvocationType: aws.String("Event"),
	})

	if err != nil {
		log.Printf("Error invoking backend lambda: %s", err)
		slackClient.PostMessage(channelId, somethingWentWrongSlackMsg)
		return
	}

	if *result.StatusCode != successfulAsyncStatusCode {
		log.Printf("Expected a status code of 202, got %d", result.StatusCode)
		slackClient.PostMessage(channelId, somethingWentWrongSlackMsg)
		return
	}
}

func failAndNotifyInSlack(message string, channelId string) (events.APIGatewayProxyResponse, error) {
	log.Print(message)
	slackClient.PostMessage(channelId, slack.MsgOptionText(message, false))
	return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, errors.New(message)
}

func withSentry(f func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	function := f

	return func(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		log.Println("Starting invocation")

		defer sentry.Recover()

		resp, err := function(ctx, e)

		if err != nil {
			log.Printf("Something went wrong! %s\n", err.Error())
			sentry.CaptureException(err)
		}

		log.Println("Finished invocation")

		return resp, err
	}
}
