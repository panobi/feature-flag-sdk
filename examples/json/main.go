package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	panobi "github.com/panobi/feature-flag-sdk"
)

func main() {
	//
	// We need the name of the JSON file.
	//

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <filename>\n", os.Args[0])
	}

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
	// Open the file and read it. Each line is assumed to be an event
	// in JSON format.
	//

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	events := make([]panobi.Event, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event panobi.Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			log.Print("Error parsing JSON:", err)
			continue
		}

		//
		// Batch the event. When we reach the max batch size, fire the
		// entire batch to Panobi.
		//

		events = append(events, event)

		if len(events) == panobi.MaxChangeEvents {
			err := client.SendEvents(events)
			if err != nil {
				log.Fatal("Error sending events:", err)
			}

			log.Printf("Successfully sent %d event(s)", len(events))

			events = make([]panobi.Event, 0)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Print("Error reading input:", err)
	}

	//
	// Fire any remaining events to Panobi.
	//

	if len(events) > 0 {
		err := client.SendEvents(events)
		if err != nil {
			log.Fatal("Error sending events:", err)
		}

		log.Printf("Successfully sent %d event(s)", len(events))
	}
}
