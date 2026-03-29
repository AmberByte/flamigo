package websocket_test

import (
	"testing"

	flamigo "github.com/amberbyte/flamigo/core"
	transportws "github.com/amberbyte/flamigo/transport/websocket"
	"github.com/stretchr/testify/assert"
)

func TestEncodeSuccess(t *testing.T) {
	raw, err := transportws.EncodeSuccess("app::users:get", map[string]string{"id": "42"}, transportws.WithAckKey("abc"))
	assert.NoError(t, err)
	assert.JSONEq(t, `{"topic":"app::users:get","ack":"abc","payload":{"id":"42"}}`, string(raw))
}

func TestEncodeError(t *testing.T) {
	raw, err := transportws.EncodeError(flamigo.NewError("nope", flamigo.Public("Forbidden"), flamigo.StatusCode(403)), transportws.WithAckKey("abc"))
	assert.NoError(t, err)
	assert.Contains(t, string(raw), `"topic":"error"`)
	assert.Contains(t, string(raw), `"ack":"abc"`)
	assert.Contains(t, string(raw), `"message":"Forbidden"`)
}
