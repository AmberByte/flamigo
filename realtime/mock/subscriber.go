package mock_realtime

import (
	"github.com/amberbyte/flamigo/realtime"
	"github.com/stretchr/testify/mock"
)

var _ realtime.Subscription = (*MockSubscriber[realtime.Event])(nil)

type MockSubscriber[T realtime.Event] struct {
	mock.Mock
}

func (m *MockSubscriber[T]) HandleEvents(event T) {
	m.MethodCalled("HandleEvents", event.Topics())
}

func (m *MockSubscriber[T]) EXPECT() *MockSubscriberExpected[T] {
	return &MockSubscriberExpected[T]{m: &m.Mock}
}

type MockSubscriberExpected[T realtime.Event] struct {
	m *mock.Mock
}

func (m *MockSubscriberExpected[T]) HandleEvents(event T) *mock.Call {
	return m.m.On("HandleEvents", event.Topics())
}

func (m *MockSubscriberExpected[T]) Cancel() *mock.Call {
	return m.m.On("Cancel")
}

func (m *MockSubscriberExpected[T]) AddTopic(topic realtime.Topic) *mock.Call {
	return m.m.On("AddTopic", topic)
}

func (m *MockSubscriberExpected[T]) RemoveTopic(topic realtime.Topic) *mock.Call {
	return m.m.On("RemoveTopic", topic)
}

func (m *MockSubscriberExpected[T]) SubscribeAll() *mock.Call {
	return m.m.On("SubscribeAll")
}

func (m *MockSubscriber[T]) Cancel() {
	m.MethodCalled("Cancel")
}

func (m *MockSubscriber[T]) AddTopic(topic realtime.Topic) {
	m.MethodCalled("AddTopic", topic)
}

func (m *MockSubscriber[T]) RemoveTopic(topic realtime.Topic) {
	m.MethodCalled("RemoveTopic", topic)
}

func (m *MockSubscriber[T]) SubscribeAll() {
	m.MethodCalled("SubscribeAll")
}

func NewMockSubscriber[T realtime.Event]() *MockSubscriber[T] {
	return &MockSubscriber[T]{}
}

func NewMockAppSubscriber() *MockSubscriber[realtime.Event] {
	return &MockSubscriber[realtime.Event]{}
}
