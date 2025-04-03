package realtime

import (
	"errors"
)

var (
	ErrInvalidEvent = errors.New("invalid event")
)

type Event interface {
	Topics() []Topic
}
