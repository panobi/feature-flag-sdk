package main

import (
	"log"
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
		log.Fatal("Error parsing key:", err)
	}

	//
	// Create a client with the signing key information.
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
		Key:          "beta-feature-xyz",
		DateModified: time.Now(),
	}

	// `Name` is a human-readable name for the feature flag that will
	// show up in the Panobi UI.
	event.SetName("Beta Feature XYZ")
	event.SetEnabled(true)

	if err := client.SendEvent(event); err != nil {
		log.Fatal("Error sending event:", err)
	}
}
