package panobi

import "time"

type Event struct {
	Project      string    `json:"project"`
	Key          string    `json:"key"`
	DateModified time.Time `json:"dateModified"`
	IsEnabled    *bool     `json:"isEnabled,omitempty"`
	Name         *string   `json:"name,omitempty"`
}

func (event *Event) SetEnabled(isEnabled bool) {
	event.IsEnabled = &isEnabled
}

func (event *Event) SetName(name string) {
	event.Name = &name
}

type ChangeEvents struct {
	Events []Event `json:"events"`
}
