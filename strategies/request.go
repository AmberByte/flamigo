package strategies

import (
	"encoding/json"
	"errors"
)

type Request struct {
	action  string
	payload any
}

func (c *Request) Action() string {
	return c.action
}

func (c *Request) Payload() any {
	return c.payload
}

func (c *Request) Bind(target any) error {
	switch v := c.payload.(type) {
	case []byte:
		return json.Unmarshal(v, target)
	case string:
		return json.Unmarshal([]byte(v), target)
	case json.RawMessage:
		return json.Unmarshal(v, target)
	default:
		return errors.New("payload cannot be parsed as json")
	}
}

func NewRequest(action string, payload any) *Request {
	return &Request{
		action:  action,
		payload: payload,
	}
}
