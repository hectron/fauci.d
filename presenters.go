package main

import (
	"fmt"
	"strings"

	"github.com/hectron/fauci.d/vaccines"
	"github.com/slack-go/slack"
)

const (
	maxNumberOfProviders = 10
	markdown             = "mrkdwn"
)

func BuildSlackBlocksForProviders(postalCode string, vaccineName string, providers []vaccines.VaccineProvider) []slack.Block {
	blocks := []slack.Block{}
	divSection := slack.NewDividerBlock()

	text := fmt.Sprintf("Found %d providers near %s offering appointments for %s!", len(providers), postalCode, vaccineName)

	if len(providers) > maxNumberOfProviders {
		text = text + fmt.Sprintf(" Only displaying the closest %d.", maxNumberOfProviders)
	}
	// header section
	headerText := slack.NewTextBlockObject(markdown, text, false, false)
	blocks = append(blocks, slack.NewSectionBlock(headerText, nil, nil))

	for idx, provider := range providers {
		if idx > maxNumberOfProviders {
			break
		}

		text := slack.NewTextBlockObject(markdown, ProviderAsString(provider), false, false)
		section := slack.NewSectionBlock(text, nil, nil)
		blocks = append(blocks, section, divSection)
	}

	return blocks
}

func ProviderAsString(provider vaccines.VaccineProvider) string {
	return fmt.Sprintf(
		"<%s|%s> - %s, %s, %s %s (about %s miles away). Phone Number: %s",
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
