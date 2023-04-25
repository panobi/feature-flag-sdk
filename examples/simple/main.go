package main

import (
	"fmt"
	"os"
	"time"

	panobi "github.com/panobi/feature-flag-sdk"
)

func main() {
	//
	// You can find your key in your Panobi workspace's integration settings.
	// It is safer to load it from an environment variable than to paste it
	// directly into this code; do not put secrets in GitHub.
	//

	k, err := panobi.ParseKey(os.Getenv("FEATURE_FLAG_SDK_SIGNING_KEY"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing key:", err)
		os.Exit(1)
	}

	//
	// Create a client with the secret key information.
	//

	client := panobi.CreateClient(k)
	defer client.Close()

	//
	// Push a single event.
	//
	// The `Project` is a way to namespace flags. It can be the project
	// you're working on, or the team you're on, or it can be empty.
	//
	// The `Key` identifies the feature flag. It can be a new flag, or an
	// existing flag. If it is new, then it will be inserted; if it is an
	// existing flag, then it will be updated. The key must be unique within
	// the `Project`.
	//
	// Here we're updating the status of the flag to enabled, meaning it is
	// live.
	//

	event := panobi.Event{
		Project:      "growth-team",
		Key:          "slackbot-greeting",
		DateModified: time.Now(),
	}
	event.SetEnabled(true)

	if err := client.SendEvent(event); err != nil {
		fmt.Fprintln(os.Stderr, "Error sending event:", err)
		os.Exit(1)
	}
}
