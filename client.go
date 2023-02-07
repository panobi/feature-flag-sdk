package panobi

import (
	"encoding/json"
	"fmt"
)

// Client for pushing feature flag events to your Panobi workspace.
type Client struct {
	t *transport
}

// Creates a new client with the given key information.
func CreateClient(k *KeyInfo) *Client {
	return &Client{
		t: createTransport(k),
	}
}

// Sends a single feature flag event to your Panobi workspace.
func (client *Client) SendEvent(event Event) error {
	return client.SendEvents([]Event{event})
}

// Sends multiple feature flag events to your Panobi workspace.
func (client *Client) SendEvents(events []Event) error {
	if len(events) > maxChangeEvents {
		return fmt.Errorf(errMaxNumberSize, "batch", maxChangeEvents, "events")
	}

	b, err := json.Marshal(&ChangeEvents{
		Events: events,
	})
	if err != nil {
		return err
	}

	return client.t.post(b)
}
