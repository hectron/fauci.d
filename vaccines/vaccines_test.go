package vaccines

import "testing"

func TestVaccineGuid(t *testing.T) {
	testCases := []struct {
		Description, Want string
		Vaccine           Vaccine
	}{
		{
			Description: "Pfizer returns the correct GUID",
			Want:        "a84fb9ed-deb4-461c-b785-e17c782ef88b",
			Vaccine:     Pfizer,
		},
		{
			Description: "Moderna returns the correct GUID",
			Want:        "779bfe52-0dd8-4023-a183-457eb100fccc",
			Vaccine:     Moderna,
		},
		{
			Description: "Johnson & Johnson returns the correct GUID",
			Want:        "784db609-dc1f-45a5-bad6-8db02e79d44f",
			Vaccine:     JJ,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			got := testCase.Vaccine.Guid()

			if got != testCase.Want {
				t.Errorf("got %v, want %v", got, testCase.Want)
			}
		})
	}
}

func TestVaccineString(t *testing.T) {
	testCases := []struct {
		Description, Want string
		Vaccine           Vaccine
	}{
		{
			Description: "Pfizer returns \"pfizer\"",
			Want:        "pfizer",
			Vaccine:     Pfizer,
		},
		{
			Description: "Moderna returns \"moderna\"",
			Want:        "moderna",
			Vaccine:     Moderna,
		},
		{
			Description: "Johnson & Johnson returns \"jj\"",
			Want:        "jj",
			Vaccine:     JJ,
		},
		{
			Description: "A vaccine name that is not recognized returns an empty string",
			Want:        "",
			Vaccine:     4,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			got := testCase.Vaccine.String()

			if got != testCase.Want {
				t.Errorf("got %v, want %v", got, testCase.Want)
			}
		})
	}
}
