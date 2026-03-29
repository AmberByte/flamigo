package client

import "encoding/json"

type MessageMarshaller interface {
	MarshalClientPayload() ([]byte, error)
}

type Message interface {
	MessageMarshaller
	Topic() string
	Payload() any
}

type message struct {
	topic   string
	payload any
}

// MarshalClientPayload marshals the payload as JSON unless it provides a
// custom client payload marshaller.
func (m *message) MarshalClientPayload() ([]byte, error) {
	if marshaller, ok := m.payload.(MessageMarshaller); ok {
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

func NewMessage(topic string, payload any) Message {
	return &message{topic: topic, payload: payload}
}
