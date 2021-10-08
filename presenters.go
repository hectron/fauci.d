package main

import (
	"fmt"
	"strings"

	"github.com/hectron/fauci.d/vaccines"
	"github.com/slack-go/slack"
)

func FormatForSlack(providers []vaccines.VaccineProvider) string {
	// var s strings.Builder
	//
	// for _, p := range providers {
	//
	// }
	return ""
	var (
		headerText *slack.TextBlockObject
	)

	headerText = slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("We found *%d* appointments! We will only be displaying the 5 closest to you."))

	msg := slack.NewBlockMessage()

	return msg
}

func ProviderAsString(provider vaccines.VaccineProvider) string {
	return fmt.Sprintf(
		"<%s|**%s**> located at %s, %s, %s %s (about %s miles away). Phone Number: %s",
		provider.Website(),
		provider.Name,
		provider.Address1,
		provider.City,
		provider.State,
		provider.Zipcode,
		strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", provider.Distance), "0"), "."),
		provider.Phone,
	)
}
