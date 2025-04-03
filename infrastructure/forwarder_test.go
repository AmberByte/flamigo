package flamigo_infra_test

import (
	"testing"
	"time"

	flamigo_infra "github.com/amberbyte/flamigo/infrastructure"
	"github.com/amberbyte/flamigo/realtime"
)

type testEvent struct {
	topics []realtime.Topic
}

func (e *testEvent) Topics() []realtime.Topic {
	return e.topics
}

func newTestEvent(topics ...realtime.Topic) realtime.Event {
	return &testEvent{topics: topics}
}

func TestForwarder(t *testing.T) {
	t.Run("It Forwards the events", func(t *testing.T) {
		f := flamigo_infra.NewForwarder[realtime.Event]()
		event := newTestEvent()
		done := make(chan bool)
		f.Subscribe(func(msg realtime.Event) {
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
