package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hectron/fauci.d/vaccines"
)

type SlackRequest struct {
	ChannelId  string           `json:"channelId"`
	PostalCode string           `json:"postalCode"`
	Vaccine    vaccines.Vaccine `json:"vaccine"`
}

func main() {
	lambda.Start(MessageHandler)
}

func MessageHandler(ctx context.Context, request SlackRequest) {
	fmt.Print("Received a lambda event!")
	fmt.Print(request.ChannelId)
	fmt.Print(request.PostalCode)
	fmt.Print(request.Vaccine.Guid())
}
