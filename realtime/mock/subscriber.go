package mock_realtime

import (
	"github.com/amberbyte/flamigo/realtime"
	"github.com/stretchr/testify/mock"
)

var _ realtime.Subscription[realtime.Event] = (*MockSubscriber[realtime.Event])(nil)

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

func (m *MockSubscriberExpected[T]) SubscribeTopic(topic string) *mock.Call {
	return m.m.On("SubscribeTopic", topic)
}

func (m *MockSubscriberExpected[T]) UnsubscribeTopic(topic string) *mock.Call {
	return m.m.On("UnsubscribeTopic", topic)
}

func (m *MockSubscriberExpected[T]) OnlyReceivables() *mock.Call {
	return m.m.On("OnlyReceivables")
}

func (m *MockSubscriber[T]) Cancel() {
	m.MethodCalled("Cancel")
}

func (m *MockSubscriber[T]) SubscribeTopic(topic string) {
	m.MethodCalled("SubscribeTopic", topic)
}

func (m *MockSubscriber[T]) UnsubscribeTopic(topic string) {
	m.MethodCalled("UnsubscribeTopic", topic)
}

func (m *MockSubscriber[T]) OnlyReceivables() {
	m.MethodCalled("OnlyReceivables")
}

func NewMockSubscriber[T realtime.Event]() *MockSubscriber[T] {
	return &MockSubscriber[T]{}
}

func NewMockAppSubscriber() *MockSubscriber[realtime.Event] {
	return &MockSubscriber[realtime.Event]{}
}
