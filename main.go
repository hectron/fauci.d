package main

import (
	"fmt"
	"net/http"

	"github.com/hectron/fauci.d/vaccines"
)

func main() {
	vaccineApi := vaccines.Api{Vaccine: vaccines.Moderna}
	client := &http.Client{}

	providers, err := vaccineApi.Request(client)

	if err == nil {
		fmt.Println(providers)
	} else {
		fmt.Println("Could not load response")
		fmt.Println(err)
	}
}
