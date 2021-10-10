package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hectron/fauci.d/vaccines"
	"github.com/slack-go/slack"
)

func FormatForSlack(providers []vaccines.VaccineProvider) string {
	blocks := []slack.Block{}
	divSection := slack.NewDividerBlock()

	// header section
	headerText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Found %d providers near you!", len(providers)), false, false)
	blocks = append(blocks, slack.NewSectionBlock(headerText, nil, nil))

	for idx, provider := range providers {
		if idx > 10 {
			break
		}

		text := slack.NewTextBlockObject("mrkdwn", ProviderAsString(provider), false, false)
		section := slack.NewSectionBlock(text, nil, nil)
		blocks = append(blocks, section, divSection)
	}

	msg := slack.NewBlockMessage(blocks...)
	b, err := json.MarshalIndent(msg, "", "    ")

	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(b)
}

func FormatForSlackUsingBlocks(providers []vaccines.VaccineProvider) []slack.Block {
	blocks := []slack.Block{}
	divSection := slack.NewDividerBlock()

	// header section
	headerText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Found %d providers near you!", len(providers)), false, false)
	blocks = append(blocks, slack.NewSectionBlock(headerText, nil, nil))

	for idx, provider := range providers {
		if idx > 10 {
			break
		}

		text := slack.NewTextBlockObject("mrkdwn", ProviderAsString(provider), false, false)
		section := slack.NewSectionBlock(text, nil, nil)
		blocks = append(blocks, section, divSection)
	}

	return blocks
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
