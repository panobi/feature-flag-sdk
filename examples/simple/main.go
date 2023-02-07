package main

import (
	"fmt"
	"os"

	panobi "github.com/panobi/feature-flag-sdk"
)

func main() {
	//
	// You can find your key in your Panobi workspace's integration settings.
	// It's safer to load it from an environment variable than to paste it
	// directly into this code; don't put secret keys in GitHub!
	//

	k, err := panobi.ParseKey(os.Getenv("FEATURE_FLAG_SDK_SECRET_KEY"))
	if err != nil {
		fmt.Printf("Error parsing key: %v\n", err)
		os.Exit(1)
	}

	//
	// Create a client with the secret key information.
	//

	client := panobi.CreateClient(k)

	//
	// Push a single event. The ExternalID identifies the feature flag. It can
	// be an existing flag, or a new flag. It must be unique. Here we're
	// the status of the flag to enabled, that is, the flag (or experiment) is
	// on in your code, and users are being bucketed.
	//

	event := panobi.Event{
		ExternalID: "slackbot-greeting",
	}
	event.SetEnabled(true)

	if err := client.SendEvent(event); err != nil {
		fmt.Printf("Error sending events: %v\n", err)
		os.Exit(1)
	}
}
