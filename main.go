package main

import (
	"fmt"

	"github.com/hectron/fauci.d/vaccines"
)

func main() {
	vaccineApi := vaccines.Api{Vaccine: vaccines.Moderna}

	providers, err := vaccineApi.Request()

	if err == nil {
		fmt.Println(providers)
	} else {
		fmt.Println("Could not load response")
		fmt.Println(err)
	}
}
