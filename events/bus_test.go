package events

import (
	"context"
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
		}, WithTopic(ParseTopic("test")))

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
		}, WithTopics(ParseTopic("test"), ParseTopic("test2")))
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
		}, WithTopics(ParseTopic("test"), ParseTopic("test2")))
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
		}, WithTopic(ParseTopic("test")))
		s.Cancel()

		bus.Publish(newTestEvent(NewTopic("test")))

		select {
		case <-done:
			t.Fatal("expected canceled subscription to stop receiving messages")
		case <-time.After(100 * time.Millisecond):
		}
	})
	t.Run("should auto cancel when lifecycle context ends", func(t *testing.T) {
		bus := NewBus[Event]()
		done := make(chan bool, 1)
		lifecycleCtx, cancel := context.WithCancel(context.Background())
		s := bus.Subscribe(func(ctx Context, msg Event) {
			done <- true
		}, WithTopic(ParseTopic("test")), WithLifecycleContext(lifecycleCtx))

		cancel()
		time.Sleep(10 * time.Millisecond)

		bus.Publish(newTestEvent(NewTopic("test")))

		select {
		case <-done:
			t.Fatal("expected lifecycle context cancelation to stop receiving messages")
		case <-time.After(100 * time.Millisecond):
		}

		// Manual cancellation remains safe after lifecycle-driven shutdown.
		s.Cancel()
	})
	t.Run("should allow adding topics after subscribing", func(t *testing.T) {
		bus := NewBus[Event]()
		done := make(chan bool, 1)
		s := bus.Subscribe(func(ctx Context, msg Event) {
			done <- true
		})

		s.AddTopic(ParseTopic("dynamic"))
		bus.Publish(newTestEvent(NewTopic("dynamic")))

		select {
		case <-done:
		case <-time.After(1 * time.Second):
			t.Fatal("expected dynamically added topic to receive messages")
		}
	})
	t.Run("should stop delivering after removing a topic", func(t *testing.T) {
		bus := NewBus[Event]()
		done := make(chan bool, 1)
		s := bus.Subscribe(func(ctx Context, msg Event) {
			done <- true
		}, WithTopic(ParseTopic("dynamic")))

		s.RemoveTopic(ParseTopic("dynamic"))
		bus.Publish(newTestEvent(NewTopic("dynamic")))

		select {
		case <-done:
			t.Fatal("expected removed topic to stop receiving messages")
		case <-time.After(100 * time.Millisecond):
		}
	})
	t.Run("should receive all topics after SubscribeAll", func(t *testing.T) {
		bus := NewBus[Event]()
		done := make(chan string, 2)
		s := bus.Subscribe(func(ctx Context, msg Event) {
			done <- msg.Topics()[0].String()
		}, WithTopic(ParseTopic("only/initial")))

		s.SubscribeAll()
		bus.Publish(newTestEvent(NewTopic("foo")))
		bus.Publish(newTestEvent(NewTopic("bar")))

		received := map[string]bool{}
		for i := 0; i < 2; i++ {
			select {
			case topic := <-done:
				received[topic] = true
			case <-time.After(1 * time.Second):
				t.Fatal("expected SubscribeAll to receive any topic")
			}
		}
		if !received["foo"] || !received["bar"] {
			t.Fatalf("expected to receive both topics, got %#v", received)
		}
	})
	t.Run("should block instead of dropping when subscriber buffer is full", func(t *testing.T) {
		bus := NewBus[Event](BusOptions{BufferSize: 1})
		started := make(chan struct{}, 1)
		release := make(chan struct{})
		publishDone := make(chan struct{})

		bus.Subscribe(func(ctx Context, msg Event) {
			select {
			case started <- struct{}{}:
			default:
			}
			<-release
		}, WithTopic(ParseTopic("backpressure")))

		go func() {
			bus.Publish(newTestEvent(NewTopic("backpressure")))
			bus.Publish(newTestEvent(NewTopic("backpressure")))
			bus.Publish(newTestEvent(NewTopic("backpressure")))
			close(publishDone)
		}()

		select {
		case <-started:
		case <-time.After(1 * time.Second):
			t.Fatal("expected first event to start processing")
		}

		select {
		case <-publishDone:
			t.Fatal("expected Publish to block while the subscriber queue is full")
		case <-time.After(100 * time.Millisecond):
		}

		close(release)

		select {
		case <-publishDone:
		case <-time.After(1 * time.Second):
			t.Fatal("expected Publish to continue once queue space is available")
		}
	})
}
