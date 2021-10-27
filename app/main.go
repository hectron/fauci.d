package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/hectron/fauci.d/mapbox"
	"github.com/hectron/fauci.d/vaccines"
	"github.com/pkg/errors"
	slackGo "github.com/slack-go/slack"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
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

	if postalCode = m.Get("text"); postalCode == "" {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, errors.New("No postal code supplied")
	}

	if channelId = m.Get("channel_id"); channelId == "" {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, errors.New("Could not determine channel to post to")
	}

	fmt.Printf("=== Requested postal code `%s` in channel id `%s`", postalCode, channelId)

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

	return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, errors.New("Invalid request")

	// coordinates, err := mapboxClient.GeocodePostalCode(postalCode)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, errors.New("Could not geocode the postal code")
	// }

	// req := vaccines.ApiRequest{
	// 	Vaccine: vaccine,
	// 	Lat:     coordinates.Latitude,
	// 	Long:    coordinates.Longitude,
	// }

	// providers, err := vaccinesClient.FindVaccines(req)

	// if err != nil {
	// 	fmt.Println("Could not load response")
	// 	fmt.Println(err)
	// 	return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, errors.New("Unable to retrieve providers")
	// }

	// blocks := slack.BuildBlocksForProviders(postalCode, vaccine.String(), providers)
	// slackClient.PostMessage(channelId, slackGo.MsgOptionBlocks(blocks...))
}

func invokeVaccineFinderLambda(channelId string, postalCode string, vaccine vaccines.Vaccine) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambda.New(sess, &aws.Config{Region: aws.String("us-east-2")})
}
