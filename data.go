package panobi

type Event struct {
	Description *string `json:"description,omitempty"`
	ExternalID  string  `json:"externalID"`
	IsEnabled   *bool   `json:"isEnabled,omitempty"`
	Name        *string `json:"name,omitempty"`
}

func (event *Event) SetDescription(description string) {
	event.Description = &description
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
