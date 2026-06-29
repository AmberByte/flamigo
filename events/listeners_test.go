package events

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	flamigo "github.com/amberbyte/flamigo"
)

type listenerTestEvent struct{}

func (listenerTestEvent) Topics() []Topic {
	return []Topic{NewTopic("listener")}
}

type listenerOtherEvent struct{}

func (listenerOtherEvent) Topics() []Topic {
	return []Topic{NewTopic("listener")}
}

func TestListenerOnEvent_DebugLogging(t *testing.T) {
	var output bytes.Buffer
	previousLogger := flamigo.Logger()

	logger := slog.New(slog.NewTextHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug}))
	flamigo.SetLogger(logger)
	defer flamigo.SetLogger(previousLogger)

	listener := ListenerOnEvent(func(ctx Context, message listenerTestEvent) {})
	listener(nil, listenerOtherEvent{})

	if output.Len() == 0 {
		t.Fatal("expected unsupported event type to be logged")
	}
}

func TestListenerOnEvent_FilteredDebugLogging(t *testing.T) {
	var output bytes.Buffer
	previousLogger := flamigo.Logger()

	logger := slog.New(flamigo.NewFilterHandler(
		slog.NewTextHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug}),
		func(_ context.Context, record slog.Record) bool {
			if record.Message != "unsupported event type" {
				return true
			}

			allowed := true
			record.Attrs(func(attr slog.Attr) bool {
				if attr.Key == "code" && attr.Value.String() == "unsupported_event_type" {
					allowed = false
					return false
				}
				return true
			})
			return allowed
		},
	))
	flamigo.SetLogger(logger)
	defer flamigo.SetLogger(previousLogger)

	listener := ListenerOnEvent(func(ctx Context, message listenerTestEvent) {})
	listener(nil, listenerOtherEvent{})

	if output.Len() != 0 {
		t.Fatal("expected unsupported event type log to be filtered")
	}
}
