package main

import (
	"fmt"
	"os"

	panobi "github.com/panobi/feature-flag-sdk"
)

func main() {
	events := []panobi.Event{
		{
			ExternalID: "slackbot-greeting",
		},
	}

	k, err := panobi.ParseKey(os.Getenv("FEATURE_FLAG_SDK_SECRET_KEY"))
	if err != nil {
		fmt.Printf("Error parsing key: %v\n", err)
		os.Exit(1)
	}

	client := panobi.CreateClient(k)

	if err := client.SendEvents(events); err != nil {
		fmt.Printf("Error sending events: %v\n", err)
		os.Exit(1)
	}
}
