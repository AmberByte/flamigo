package realtime

import (
	"errors"
)

var (
	ErrInvalidEvent = errors.New("invalid event")
)

// Event is the interface that all events must implement.
type Event interface {
	// Topic returns the topics under which the event is published.
	//
	// This is used to determine which events are sent to which subscribers.
	Topics() []Topic
}
