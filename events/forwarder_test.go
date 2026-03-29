package events_test

import (
	"testing"
	"time"

	"github.com/amberbyte/flamigo/events"
)

type testEvent struct {
	topics []events.Topic
}

func (e *testEvent) Topics() []events.Topic {
	return e.topics
}

func newTestEvent(topics ...events.Topic) events.Event {
	return &testEvent{topics: topics}
}

func TestForwarder(t *testing.T) {
	t.Run("It Forwards the events", func(t *testing.T) {
		f := events.NewForwarder[events.Event]()
		event := newTestEvent()
		done := make(chan bool)
		f.Subscribe(func(msg events.Event) {
			if msg != event {
				t.Errorf("expected event to be forwarded")
			}
			done <- true
		})
		f.Publish(event)
		select {
		case <-done:
		case <-time.After(1 * time.Second):
			t.Fatal("expected message to be received")
		}
	})
}
