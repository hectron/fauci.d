package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

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
	slackClient                                *slack.Client
	successfulAsyncStatusCode                  int64
	backendFunctionName, somethingWentWrongMsg string
	somethingWentWrongSlackMsg                 slack.MsgOption
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

	backendFunctionName = os.Getenv("AWS_LAMBDA_FUNCTION_NAME") + "_backend"
	slackClient = slack.New(os.Getenv("SLACK_API_TOKEN"))
}

func main() {
	lambda.Start(IncomingMessageHandler)
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

	fmt.Printf("=== Request: %s\n", request.Body)

	if channelId = m.Get("channel_id"); channelId == "" {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, errors.New("Could not determine channel to post to")
	}

	if postalCode = m.Get("text"); postalCode == "" {
		return failAndNotifyInSlack("No postal code supplied.", channelId)
	}

	fmt.Printf("=== Requested postal code `%s` in channel id `%s`\n", postalCode, channelId)

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
		msg := fmt.Sprintf("Could not generate requst for backend lambda: %s", err)
		fmt.Print(msg)
		slackClient.PostMessage(channelId, slack.MsgOptionText(msg, false))
		return
	}

	result, err := lambdaClient.Invoke(&lambdaSdk.InvokeInput{
		FunctionName:   aws.String(backendFunctionName),
		Payload:        payload,
		InvocationType: aws.String("Event"),
	})

	if err != nil {
		fmt.Printf("Error invoking backend lambda: %s", err)
		slackClient.PostMessage(channelId, somethingWentWrongSlackMsg)
		return
	}

	if *result.StatusCode != successfulAsyncStatusCode {
		fmt.Printf("Expected a status code of 202, got %d", result.StatusCode)
		slackClient.PostMessage(channelId, somethingWentWrongSlackMsg)
		return
	}

	fmt.Printf("Successfully invoked %s", backendFunctionName)
}

func failAndNotifyInSlack(message string, channelId string) (events.APIGatewayProxyResponse, error) {
	fmt.Print(message)
	slackClient.PostMessage(channelId, slack.MsgOptionText(message, false))
	return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, errors.New(message)
}
