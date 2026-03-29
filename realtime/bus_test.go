package realtime

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type testEvent struct {
	topics []Topic
}

func (e *testEvent) Topics() []Topic {
	return e.topics
}

func newTestEvent(topics ...Topic) Event {
	return &testEvent{topics: topics}
}

func TestBus(t *testing.T) {
	t.Run("should subscribe and publish messages", func(t *testing.T) {
		bus := NewBus[Event]()
		done := make(chan bool)
		bus.Subscribe(func(ctx Context, msg Event) {
			for _, topic := range msg.Topics() {
				if topic.String() != "test" {
					t.Errorf("expected 'test' got %s", topic)
				}
			}
			done <- true
		}, WithTopic("test"))

		evt := newTestEvent(NewTopic("test"))

		bus.Publish(evt)
		select {
		case <-done:
		case <-time.After(1 * time.Second):
			t.Fatal("expected message to be received")
		}

	})
	t.Run("should receive callback only once", func(t *testing.T) {
		bus := NewBus[Event]()
		fnMock := mock.Mock{}
		bus.Subscribe(func(ctx Context, msg Event) {
			fnMock.MethodCalled("Called", ctx, msg)
		}, WithTopics("test", "test2"))
		// Expect the function to only be called once
		fnMock.On("Called", mock.Anything, mock.Anything).Once().Return(nil)

		evt := newTestEvent(NewTopic("test"), NewTopic("test2"))

		bus.Publish(evt)
		<-time.After(1 * time.Second)
		fnMock.AssertExpectations(t)
	})

	t.Run("should be able to send events synchronously", func(t *testing.T) {
		bus := NewBus[Event]()
		fnMock := mock.Mock{}
		bus.Subscribe(func(ctx Context, msg Event) {
			// this makes the test work. when the PublishSync would return early it would fail
			<-time.After(1 * time.Second)
			fnMock.MethodCalled("Called", ctx, msg)
		}, WithTopics("test", "test2"))
		// Expect the function to only be called once
		fnMock.On("Called", mock.Anything, mock.Anything).Once().Return(nil)

		evt := newTestEvent(NewTopic("test"), NewTopic("test2"))

		bus.PublishSync(evt)
		fnMock.AssertExpectations(t)
	})
	t.Run("should not deliver to canceled subscriptions", func(t *testing.T) {
		bus := NewBus[Event]()
		done := make(chan bool, 1)
		s := bus.Subscribe(func(ctx Context, msg Event) {
			done <- true
		}, WithTopic("test"))
		s.Cancel()

		bus.Publish(newTestEvent(NewTopic("test")))

		select {
		case <-done:
			t.Fatal("expected canceled subscription to stop receiving messages")
		case <-time.After(100 * time.Millisecond):
		}
	})
	t.Run("should allow adding topics after subscribing", func(t *testing.T) {
		bus := NewBus[Event]()
		done := make(chan bool, 1)
		s := bus.Subscribe(func(ctx Context, msg Event) {
			done <- true
		})

		s.AddTopic("dynamic")
		bus.Publish(newTestEvent(NewTopic("dynamic")))

		select {
		case <-done:
		case <-time.After(1 * time.Second):
			t.Fatal("expected dynamically added topic to receive messages")
		}
	})
}

type benchmarkEvent struct {
	topic    string
	payload  string
	receiver bool
}

func (e benchmarkEvent) Topics() []Topic {
	return []Topic{NewTopic(e.topic)}
}

func BenchmarkBusPublish(b *testing.B) {
	tests := []struct {
		subscribers int
		topics      int
	}{
		{subscribers: 1, topics: 1},
		{subscribers: 10, topics: 1},
		{subscribers: 100, topics: 1},
		{subscribers: 10, topics: 10},
		{subscribers: 100, topics: 10},
	}

	for _, tt := range tests {
		b.Run(fmt.Sprintf("subs=%d/topics=%d", tt.subscribers, tt.topics), func(b *testing.B) {
			bus := NewBus[benchmarkEvent]()

			// Setup subscribers
			for i := 0; i < tt.subscribers; i++ {
				bus.Subscribe(func(ctx Context, event benchmarkEvent) {
					// Simulate some work
					_ = event.payload
				}, benchmarkTopicOpts(tt.topics)...)
			}

			event := benchmarkEvent{
				topic:    "topic-0",
				payload:  "test payload",
				receiver: true,
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bus.Publish(event)
			}
		})
	}
}

func BenchmarkBusPublishSync(b *testing.B) {
	tests := []struct {
		subscribers int
		topics      int
	}{
		{subscribers: 1, topics: 1},
		{subscribers: 10, topics: 1},
		{subscribers: 100, topics: 1},
		{subscribers: 10, topics: 10},
		{subscribers: 100, topics: 10},
	}

	for _, tt := range tests {
		b.Run(fmt.Sprintf("subs=%d/topics=%d", tt.subscribers, tt.topics), func(b *testing.B) {
			bus := NewBus[benchmarkEvent]()

			// Setup subscribers
			for i := 0; i < tt.subscribers; i++ {
				bus.Subscribe(func(ctx Context, event benchmarkEvent) {
					// Simulate some work
					_ = event.payload
				}, benchmarkTopicOpts(tt.topics)...)
			}

			event := benchmarkEvent{
				topic:    "topic-0",
				payload:  "test payload",
				receiver: true,
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bus.PublishSync(event)
			}
		})
	}
}

func BenchmarkBusSubscribe(b *testing.B) {
	bus := NewBus[benchmarkEvent]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sub := bus.Subscribe(func(ctx Context, event benchmarkEvent) {
			// Simulate some work
			_ = event.payload
		}, WithAllTopics())
		sub.Cancel()
	}
}

func benchmarkTopicOpts(topics int) []SubscribeOpt {
	opts := make([]SubscribeOpt, 0, topics)
	for j := 0; j < topics; j++ {
		opts = append(opts, WithTopic(fmt.Sprintf("topic-%d", j)))
	}
	return opts
}
