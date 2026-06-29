package events

import (
	"fmt"
	"log/slog"

	flamigo "github.com/amberbyte/flamigo"
)

type ForwarderListener[T Event] func(message T)
type BusListener[T Event] func(context Context, message T)

type AppListener = BusListener[Event]

// ListenerForwarder is a type helper to be able to use custom event types in listeners.
//
// Example:
//
//	events.ListenerOnEvent(func(ctx events.Context, message MyEvent) {})
func ListenerOnEvent[T Event](listener BusListener[T]) AppListener {
	return func(context Context, message Event) {
		if typed, ok := message.(T); ok {
			listener(context, typed)
			return
		}
		flamigo.Logger().Debug(
			"unsupported event type",
			slog.String("component", "events"),
			slog.String("code", "unsupported_event_type"),
			slog.String("event_type", fmt.Sprintf("%T", message)),
		)
	}
}
