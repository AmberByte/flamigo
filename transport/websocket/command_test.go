package websocket_test

import (
	"testing"

	transportws "github.com/amberbyte/flamigo/transport/websocket"
	"github.com/stretchr/testify/assert"
)

func TestDecodeCommand(t *testing.T) {
	cmd, err := transportws.DecodeCommand([]byte(`{"topic":"app::users:get","payload":{"id":"42"},"ack":"abc"}`))
	assert.NoError(t, err)
	assert.Equal(t, "app::users:get", cmd.Action())
	assert.Equal(t, "abc", cmd.AckKey())
}

func TestDecodeCommandRejectsMissingTopic(t *testing.T) {
	_, err := transportws.DecodeCommand([]byte(`{"payload":{"id":"42"}}`))
	assert.ErrorIs(t, err, transportws.ErrInvalidCommand)
}
