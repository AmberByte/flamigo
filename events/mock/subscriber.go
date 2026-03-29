package mock_events

import (
	"github.com/amberbyte/flamigo/events"
	"github.com/stretchr/testify/mock"
)

var _ events.Subscription = (*MockSubscriber[events.Event])(nil)

type MockSubscriber[T events.Event] struct {
	mock.Mock
}

func (m *MockSubscriber[T]) HandleEvents(event T) {
	m.MethodCalled("HandleEvents", event.Topics())
}

func (m *MockSubscriber[T]) EXPECT() *MockSubscriberExpected[T] {
	return &MockSubscriberExpected[T]{m: &m.Mock}
}

type MockSubscriberExpected[T events.Event] struct {
	m *mock.Mock
}

func (m *MockSubscriberExpected[T]) HandleEvents(event T) *mock.Call {
	return m.m.On("HandleEvents", event.Topics())
}

func (m *MockSubscriberExpected[T]) Cancel() *mock.Call {
	return m.m.On("Cancel")
}

func (m *MockSubscriberExpected[T]) AddTopic(topic events.Topic) *mock.Call {
	return m.m.On("AddTopic", topic)
}

func (m *MockSubscriberExpected[T]) RemoveTopic(topic events.Topic) *mock.Call {
	return m.m.On("RemoveTopic", topic)
}

func (m *MockSubscriberExpected[T]) SubscribeAll() *mock.Call {
	return m.m.On("SubscribeAll")
}

func (m *MockSubscriber[T]) Cancel() {
	m.MethodCalled("Cancel")
}

func (m *MockSubscriber[T]) AddTopic(topic events.Topic) {
	m.MethodCalled("AddTopic", topic)
}

func (m *MockSubscriber[T]) RemoveTopic(topic events.Topic) {
	m.MethodCalled("RemoveTopic", topic)
}

func (m *MockSubscriber[T]) SubscribeAll() {
	m.MethodCalled("SubscribeAll")
}

func NewMockSubscriber[T events.Event]() *MockSubscriber[T] {
	return &MockSubscriber[T]{}
}

func NewMockAppSubscriber() *MockSubscriber[events.Event] {
	return &MockSubscriber[events.Event]{}
}
