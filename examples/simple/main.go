package main

import (
	"fmt"
	"os"

	panobi "github.com/panobi/feature-flag-sdk"
)

func main() {
	//
	// You can find your key in your Panobi workspace's integration settings.
	// It is safer to load it from an environment variable than to paste it
	// directly into this code; do not put secret keys in GitHub.
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
	// Push a single event.
	//
	// The `ExternalID` identifies the feature flag. It can be an existing
	// flag, or a new flag. It must be unique.
	//
	// Here we're updating the status of the flag to enabled, meaning it is
	// live, and users are presumably being bucketed.
	//
	// You can also update the name of the flag and the flag's description.
	// But the `ExternalID` is fixed.
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
