package panobi

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	MaxChangeEvents    int           = 64
	bufferedSendPeriod time.Duration = 10 * time.Second
)

// Client for pushing feature flag events to your Panobi workspace.
type client struct {
	t *transport
}

// Creates a new client with the given key information.
func CreateClient(k KeyInfo) *client {
	c := &client{
		t: createTransport(k),
	}

	return c
}

func (client *client) Close() {
}

// Sends a single feature flag event to your Panobi workspace.
func (client *client) SendEvent(event Event) error {
	return client.SendEvents([]Event{event})
}

// Sends multiple feature flag events to your Panobi workspace.
func (client *client) SendEvents(events []Event) error {
	if len(events) > MaxChangeEvents {
		return fmt.Errorf(errMaxNumberSize, "batch", MaxChangeEvents, "events")
	}

	b, err := json.Marshal(&ChangeEvents{
		Events: events,
	})
	if err != nil {
		return err
	}

	_, err = client.t.post(b)
	return err
}
