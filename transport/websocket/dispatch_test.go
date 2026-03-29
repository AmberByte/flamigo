package websocket_test

import (
	"context"
	"errors"
	"testing"

	flamigo "github.com/amberbyte/flamigo/core"
	"github.com/amberbyte/flamigo/strategies"
	transportws "github.com/amberbyte/flamigo/transport/websocket"
	"github.com/stretchr/testify/assert"
)

type stubRouter struct {
	invoke func(ctx strategies.Context) strategies.Result
}

func (s *stubRouter) Register(namespace string, registry strategies.AppRegistry) error {
	return nil
}

func (s *stubRouter) Invoke(ctx strategies.Context) strategies.Result {
	return s.invoke(ctx)
}

func TestDispatcherHandleMessageSuccess(t *testing.T) {
	router := &stubRouter{
		invoke: func(ctx strategies.Context) strategies.Result {
			ack, ok := transportws.AckKey(ctx)
			assert.True(t, ok)
			assert.Equal(t, "ack-1", ack)
			ctx.Response().SetPayload(map[string]string{"ok": "true"})
			return ctx.Response()
		},
	}

	dispatcher := transportws.NewDispatcher(router)
	var wrote []byte
	err := dispatcher.HandleMessage(context.Background(), []byte(`{"topic":"app::users:get","payload":{"id":"42"},"ack":"ack-1"}`), func(raw []byte) error {
		wrote = raw
		return nil
	})
	assert.NoError(t, err)
	assert.JSONEq(t, `{"topic":"app::users:get","ack":"ack-1","payload":{"ok":"true"}}`, string(wrote))
}

func TestDispatcherHandleMessageError(t *testing.T) {
	router := &stubRouter{
		invoke: func(ctx strategies.Context) strategies.Result {
			ctx.Response().SetError(flamigo.NewError("forbidden", flamigo.Public("Forbidden"), flamigo.Kind(flamigo.ErrorTypeForbidden)))
			return ctx.Response()
		},
	}

	dispatcher := transportws.NewDispatcher(router)
	var wrote []byte
	err := dispatcher.HandleMessage(context.Background(), []byte(`{"topic":"app::users:get","ack":"ack-1"}`), func(raw []byte) error {
		wrote = raw
		return nil
	})
	assert.NoError(t, err)
	assert.Contains(t, string(wrote), `"topic":"error"`)
	assert.Contains(t, string(wrote), `"ack":"ack-1"`)
}

func TestDispatcherHandleMessageActorFactoryError(t *testing.T) {
	router := &stubRouter{
		invoke: func(ctx strategies.Context) strategies.Result {
			t.Fatal("router should not be called")
			return nil
		},
	}

	dispatcher := transportws.NewDispatcher(router, transportws.WithActorFactory(func(ctx context.Context, cmd *transportws.Command) (flamigo.Actor, error) {
		return nil, errors.New("boom")
	}))

	var wrote []byte
	err := dispatcher.HandleMessage(context.Background(), []byte(`{"topic":"app::users:get","ack":"ack-1"}`), func(raw []byte) error {
		wrote = raw
		return nil
	})
	assert.NoError(t, err)
	assert.Contains(t, string(wrote), `"topic":"error"`)
	assert.Contains(t, string(wrote), `"ack":"ack-1"`)
}
