package panobi

import "time"

const (
	idLen int = 22

	maxInputBytes   int = 1_048_576
	maxChangeEvents int = 64

	bufferedSendPeriod time.Duration = 10 * time.Second

	defaultRetryAfter time.Duration = 1 * time.Second
	maxRetries        int           = 2

	eventsURI string = "http://localhost:8080/integrations/flags/events"

	errInvalidKey    string = "invalid key"
	errMaxNumberSize string = "%s cannot be larger than %d %s"
)
