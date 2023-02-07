package panobi

import (
	"encoding/json"
	"fmt"
)

type Client struct {
	t *transport
}

func CreateClient(k *KeyInfo) *Client {
	return &Client{
		t: createTransport(k),
	}
}

func (client *Client) SendEvent(event Event) error {
	return client.SendEvents([]Event{event})
}

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
