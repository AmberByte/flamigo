package strategies

import (
	"encoding/json"
	"errors"

	"github.com/amberbyte/flamigo/internal"
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

// BindAndValidate binds the payload to the target and validates it using validator/v10.
func (c *Request) BindAndValidate(target any) error {
	if err := c.Bind(target); err != nil {
		return err
	}

	if err := internal.Validate(target); err != nil {
		return err
	}
	return nil
}

func NewRequest(action string, payload any) *Request {
	return &Request{
		action:  action,
		payload: payload,
	}
}
