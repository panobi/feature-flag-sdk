package panobi

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	maxChangeEvents    int           = 64
	bufferedSendPeriod time.Duration = 10 * time.Second
)

// Client for pushing feature flag events to your Panobi workspace.
type client struct {
	t    *transport
	q    chan Event
	done chan bool
	wg   sync.WaitGroup
}

// Creates a new client with the given key information.
func CreateClient(k *KeyInfo) *client {
	c := &client{
		t:    createTransport(k),
		q:    make(chan Event, maxChangeEvents),
		done: make(chan bool),
	}

	c.wg.Add(1)
	c.startBufferedSender()

	return c
}

func (client *client) Close() {
	client.done <- true
	client.wg.Wait()
}

// Sends a single feature flag event to your Panobi workspace.
func (client *client) SendEvent(event Event) error {
	return client.SendEvents([]Event{event})
}

// Sends multiple feature flag events to your Panobi workspace.
func (client *client) SendEvents(events []Event) error {
	if len(events) > maxChangeEvents {
		return fmt.Errorf(errMaxNumberSize, "batch", maxChangeEvents, "events")
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

// Buffers an event so that it can be sent in a batch.
func (client *client) SendEventBuffered(event Event) {
	client.q <- event
}

func (client *client) startBufferedSender() {
	ticker := time.NewTicker(bufferedSendPeriod)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				client.drain()
			case <-client.done:
				client.drain()
				return
			}
		}
	}()
}

func (client *client) drain() {
	events := drain(client.q)

	if len(events) > 0 {
		b, err := json.Marshal(&ChangeEvents{
			Events: events,
		})
		switch err {
		case nil:
			_, postErr := client.t.post(b)
			if postErr != nil {
				fmt.Fprintln(os.Stderr, "Error sending event:", postErr)
			}
		default:
			fmt.Fprintln(os.Stderr, "Error marshalling event:", err)
		}
	}
}

func drain(q chan Event) []Event {
	events := make([]Event, 0)

	for {
		select {
		case e := <-q:
			events = append(events, e)
		default:
			return events
		}
	}
}
