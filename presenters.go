package main

import (
	"fmt"
	"strings"

	"github.com/hectron/fauci.d/vaccines"
)

func FormatForSlack(providers []vaccines.VaccineProvider) string {
	// var s strings.Builder
	//
	// for _, p := range providers {
	//
	// }
	return ""
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
