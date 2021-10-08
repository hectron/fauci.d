package main

import (
	"testing"

	"github.com/hectron/fauci.d/vaccines"
)

func TestProviderAsString(t *testing.T) {
	testCases := []struct {
		Description, Want string
		Provider          vaccines.VaccineProvider
	}{
		{
			"It formats a provider correctly",
			"<https://www.vaccines.gov/provider/?id=1|**Tester**> located at 123 Fake St, Milwaukee, WI 53201 (about 0.5 miles away). Phone Number: (312) 111-5555",
			vaccines.VaccineProvider{
				Guid:     "1",
				Name:     "Tester",
				Address1: "123 Fake St",
				City:     "Milwaukee",
				State:    "WI",
				Zipcode:  "53201",
				Distance: 0.5,
				Phone:    "(312) 111-5555",
			},
		},
		{
			"It formats a provider with a single-decimal distance correctly",
			"<https://www.vaccines.gov/provider/?id=1|**Tester**> located at 123 Fake St, Milwaukee, WI 53201 (about 0.5 miles away). Phone Number: ",
			vaccines.VaccineProvider{
				Guid:     "1",
				Name:     "Tester",
				Address1: "123 Fake St",
				City:     "Milwaukee",
				State:    "WI",
				Zipcode:  "53201",
				Distance: 0.5,
				Phone:    "",
			},
		},
		{
			"It formats a provider with a two-decimal distance correctly",
			"<https://www.vaccines.gov/provider/?id=1|**Tester**> located at 123 Fake St, Milwaukee, WI 53201 (about 0.5 miles away). Phone Number: ",
			vaccines.VaccineProvider{
				Guid:     "1",
				Name:     "Tester",
				Address1: "123 Fake St",
				City:     "Milwaukee",
				State:    "WI",
				Zipcode:  "53201",
				Distance: 0.5,
				Phone:    "",
			},
		},
		{
			"It limits the distance to two decimal places maximum",
			"<https://www.vaccines.gov/provider/?id=1|**Tester**> located at 123 Fake St, Milwaukee, WI 53201 (about 0.84 miles away). Phone Number: ",
			vaccines.VaccineProvider{
				Guid:     "1",
				Name:     "Tester",
				Address1: "123 Fake St",
				City:     "Milwaukee",
				State:    "WI",
				Zipcode:  "53201",
				Distance: 0.843595239,
				Phone:    "",
			},
		},
		{
			"It does not show the decimal distance if the distance is a whole number",
			"<https://www.vaccines.gov/provider/?id=1|**Tester**> located at 123 Fake St, Milwaukee, WI 53201 (about 3 miles away). Phone Number: ",
			vaccines.VaccineProvider{
				Guid:     "1",
				Name:     "Tester",
				Address1: "123 Fake St",
				City:     "Milwaukee",
				State:    "WI",
				Zipcode:  "53201",
				Distance: 3.0,
				Phone:    "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			got := ProviderAsString(tc.Provider)

			if got != tc.Want {
				t.Errorf("\ngot:  %s\nwant: %s", got, tc.Want)
			}
		})
	}
}
