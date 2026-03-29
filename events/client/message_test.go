package client

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type payloadMarshaller struct {
	raw []byte
}

func (p payloadMarshaller) MarshalClientPayload() ([]byte, error) {
	return p.raw, nil
}

func TestMessage(t *testing.T) {
	t.Run("marshals plain payloads as json", func(t *testing.T) {
		msg := NewMessage("topic", map[string]string{"foo": "bar"})
		raw, err := msg.MarshalClientPayload()
		assert.NoError(t, err)

		var payload map[string]string
		err = json.Unmarshal(raw, &payload)
		assert.NoError(t, err)
		assert.Equal(t, "topic", msg.Topic())
		assert.Equal(t, map[string]string{"foo": "bar"}, payload)
	})

	t.Run("uses custom payload marshaller when available", func(t *testing.T) {
		msg := NewMessage("topic", payloadMarshaller{raw: []byte(`{"ok":true}`)})
		raw, err := msg.MarshalClientPayload()
		assert.NoError(t, err)
		assert.JSONEq(t, `{"ok":true}`, string(raw))
	})
}
