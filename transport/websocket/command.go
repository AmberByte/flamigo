package websocket

import (
	"encoding/json"
	"errors"
)

var ErrInvalidCommand = errors.New("invalid websocket command")

type Command struct {
	topic   string
	payload json.RawMessage
	ackKey  string
}

type commandData struct {
	Topic   string          `json:"topic"`
	Payload json.RawMessage `json:"payload,omitempty"`
	AckKey  string          `json:"ack,omitempty"`
}

func (c *Command) Action() string {
	return c.topic
}

func (c *Command) Payload() any {
	return c.payload
}

func (c *Command) AckKey() string {
	return c.ackKey
}

func (c *Command) UnmarshalPayload(target any) error {
	return json.Unmarshal(c.payload, target)
}

func DecodeCommand(raw []byte) (*Command, error) {
	var data commandData
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}
	if data.Topic == "" {
		return nil, ErrInvalidCommand
	}
	return &Command{
		topic:   data.Topic,
		payload: data.Payload,
		ackKey:  data.AckKey,
	}, nil
}
