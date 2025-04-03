package mock_realtime

import (
	"errors"

	flamigo "github.com/amberbyte/flamigo/core"
	"github.com/amberbyte/flamigo/internal"
	"github.com/amberbyte/flamigo/realtime"
	"github.com/stretchr/testify/mock"
)

type MockPublisher[Evt realtime.Event] struct {
	mock.Mock
}

func (m *MockPublisher[Evt]) Publish(event realtime.Event) {
	for _, topic := range event.Topics() {
		m.MethodCalled("Publish", topic, event)
	}
}

func (m *MockPublisher[Evt]) ExpectPublish(topic string, payload ...any) *mock.Call {
	payloadD := internal.ParseOptionalParam(payload, mock.Anything)
	return m.On("Publish", topic, payloadD)
}

func NewMockPublisher() *MockPublisher[realtime.Event] {
	return &MockPublisher[realtime.Event]{}
}

func NewCustomMockPublisher[Evt realtime.Event]() *MockPublisher[Evt] {
	return &MockPublisher[Evt]{}
}

var _ realtime.AppBus = (*MockBus[realtime.Event])(nil)

type MockBus_Expecter struct {
	mock *mock.Mock
}

func (m *MockBus_Expecter) Subscribe(subscription any) *mock.Call {
	return m.mock.On("Subscribe", subscription)
}

func (m *MockBus_Expecter) Publish(event any, actor any) *mock.Call {
	return m.mock.On("Publish", event, actor)
}

func (m *MockBus_Expecter) PublishSync(subscription any, actor any) *mock.Call {
	// In mocking we do not distinguish between sync and async
	return m.Publish(subscription, actor)
}

type MockAppBus = MockBus[realtime.Event]

type MockBus[Evt realtime.Event] struct {
	mock.Mock
	listeners map[realtime.Subscription[Evt]]realtime.BusListener[realtime.Event]
}

func (m *MockBus[Evt]) Subscribe(subscription realtime.BusListener[realtime.Event]) realtime.Subscription[Evt] {
	args := m.Called(subscription)
	subscriber := args.Get(0).(realtime.Subscription[Evt])
	m.listeners[subscriber] = subscription
	return subscriber
}

func (m *MockBus[Evt]) Publish(event realtime.Event, actor ...flamigo.Actor) {
	actorD := internal.ParseOptionalParam(actor, nil)
	m.Called(event, actorD)
}

func (m *MockBus[Evt]) PublishSync(event realtime.Event, actor ...flamigo.Actor) {
	// In mocking we do not distinguish between sync and async
	m.Publish(event, actor...)
}

func (m *MockBus[Evt]) TRIGGER(subscription realtime.Subscription[Evt], ctx realtime.Context, event realtime.Event) error {
	if listener, ok := m.listeners[subscription]; ok {
		listener(ctx, event)
		return nil
	}
	return errors.New("suscription not yet registered an listener")
}

func (m *MockBus[Evt]) EXPECT() *MockBus_Expecter {
	return &MockBus_Expecter{&m.Mock}
}

func NewBus() *MockBus[realtime.Event] {
	return &MockBus[realtime.Event]{
		listeners: make(map[realtime.Subscription[realtime.Event]]realtime.BusListener[realtime.Event]),
	}
}
