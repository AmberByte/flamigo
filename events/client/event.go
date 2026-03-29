package client

// Event is a event that can be sent to a client
type Event interface {
	ClientMessage() Message
}

// IsEvent checks if the given Event implements ClientEvent
func IsEvent(i any) bool {
	_, ok := i.(Event)
	return ok
}
