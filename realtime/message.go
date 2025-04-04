package realtime

import "encoding/json"

type ClientMessageMarshaller interface {
	MarshalClientPayload() ([]byte, error)
}

// Client event is a event that can be sent to a client
type ClientEvent interface {
	ClientMessage() ClientMessage
}

// IsClientEvent checks if the given Event implements ClientEvent
func IsClientEvent(i any) bool {
	_, ok := i.(ClientEvent)
	return ok
}

type ClientMessage interface {
	ClientMessageMarshaller
	Topic() string
	Payload() any
}

type message struct {
	topic   string
	payload interface{}
}

// MarshalClientPayload marshals the payload to a json encoded stirng.
// If the payload implements ClientMessageMarshaller interface, this will be used.
func (m *message) MarshalClientPayload() ([]byte, error) {
	if marshaller, ok := m.payload.(ClientMessageMarshaller); ok {
		return marshaller.MarshalClientPayload()
	}
	return json.Marshal(m.payload)
}

// Topic returns the topic of the message
//
// IMPORTANT: This is not the topic of the event, but the topic of the message.
// e.g. updated:foo
func (m *message) Topic() string {
	return m.topic
}

// Payload returns the payload of the message
func (m *message) Payload() any {
	return m.payload
}

// NewClientMessage creates a new ClientMessage with the given topic and payload
func NewClientMessage(topic string, payload any) *message {
	return &message{topic: topic, payload: payload}
}
