package realtime

import "encoding/json"

type ClientMessageMarshaller interface {
	MarshalClientPayload() ([]byte, error)
}

// Client event is a event that can be sent to a client
type ClientEvent interface {
	ClientMessage() ClientMessage
}

func IsClientEvent(i interface{}) bool {
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

func (m *message) MarshalClientPayload() ([]byte, error) {
	if marshaller, ok := m.payload.(ClientMessageMarshaller); ok {
		return marshaller.MarshalClientPayload()
	}
	return json.Marshal(m.payload)
}

func (m *message) Topic() string {
	return m.topic
}

func (m *message) Payload() any {
	return m.payload
}

func NewClientMessage(topic string, payload any) *message {
	return &message{topic: topic, payload: payload}
}
