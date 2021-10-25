package slack

import (
	"fmt"
	"strings"

	"github.com/hectron/fauci.d/vaccines"
	slackGo "github.com/slack-go/slack"
)

const (
	maxNumberOfProviders = 5
	markdown             = "mrkdwn"
)

func BuildBlocksForProviders(postalCode string, vaccineName string, providers []vaccines.VaccineProvider) []slackGo.Block {
	blocks := []slackGo.Block{}
	divSection := slackGo.NewDividerBlock()

	text := fmt.Sprintf("Found %d providers near %s offering %s appointments.", len(providers), postalCode, vaccineName)

	if len(providers) > maxNumberOfProviders {
		text = text + fmt.Sprintf(" Only displaying the closest %d.", maxNumberOfProviders)
	}
	// header section
	headerText := slackGo.NewTextBlockObject(markdown, text, false, false)
	blocks = append(blocks, slackGo.NewSectionBlock(headerText, nil, nil))

	for idx, provider := range providers {
		if idx > maxNumberOfProviders {
			break
		}

		text := slackGo.NewTextBlockObject(markdown, ProviderAsString(provider), false, false)
		section := slackGo.NewSectionBlock(text, nil, nil)
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
