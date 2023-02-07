package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"

	panobi "github.com/panobi/feature-flag-sdk"
)

func main() {
	//
	// We need the name of the CSV file.
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
	// Open the file and read it. Each line is assumed to be a list of
	// comma-separated values, with the columns in the following order. The
	// last two columns are optional, and can be omitted.
	//
	// Project, Key, DateModified, IsEnabled, Name
	//

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	events := make([]panobi.Event, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		cols := strings.Split(line, ",")

		if len(cols) < 3 {
			log.Fatalf(`Expected at least three columns in "%s"`, line)
		}

		dateModified := coerceStringToDate(cols[2])
		if dateModified == nil {
			log.Fatalf(`Invalid timestamp "%s"`, cols[2])
		}

		event := panobi.Event{
			Project:      strings.TrimSpace(cols[0]),
			Key:          strings.TrimSpace(cols[1]),
			DateModified: *dateModified,
		}

		//
		// The last two columns are optional.
		//

		if len(cols) > 3 {
			event.SetEnabled(isTrue(cols[3]))
		}

		if len(cols) > 4 {
			event.SetName(strings.TrimSpace(cols[4]))
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

func coerceStringToDate(s string) *time.Time {
	formats := []string{
		time.DateOnly,
		time.DateTime,
		time.RFC3339,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC850,
		time.RFC822,
		time.UnixDate,
		time.RubyDate,
	}
	for _, format := range formats {
		t, err := time.Parse(format, s)
		if err == nil && t.Year() > 1 {
			return &t
		}
	}
	return nil
}

func isTrue(val string) bool {
	return val == "1" || strings.ToLower(val) == "true"
}
