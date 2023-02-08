package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

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
		fmt.Fprintln(os.Stderr, "Error parsing key:", err)
		os.Exit(1)
	}

	//
	// Create a client with the secret key information.
	//

	client := panobi.CreateClient(k)
	defer client.Close()

	//
	// Read lines from standard input. Each line is assumed to be an event
	// in JSON format. Events will be buffered and sent in batches.
	//

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			var event panobi.Event
			if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
				fmt.Fprintln(os.Stderr, "Error parsing JSON:", err)
				continue
			}

			if strings.TrimSpace(event.ExternalID) != "" {
				client.SendEventBuffered(event)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
		}
	}()

	wg.Wait()
}
