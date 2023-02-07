package panobi

type Event struct {
	Description *string `json:"description,omitempty"`
	ExternalID  string  `json:"externalID"`
	IsEnabled   *bool   `json:"isEnabled,omitempty"`
	Name        *string `json:"name,omitempty"`
}

type ChangeEvents struct {
	Events []Event `json:"events"`
}
