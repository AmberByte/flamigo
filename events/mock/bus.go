package mock_events

import (
	"errors"

	flamigo "github.com/amberbyte/flamigo/core"
	"github.com/amberbyte/flamigo/internal"
	"github.com/amberbyte/flamigo/events"
	"github.com/stretchr/testify/mock"
)

type MockPublisher[Evt events.Event] struct {
	mock.Mock
}

func (m *MockPublisher[Evt]) Publish(event events.Event) {
	for _, topic := range event.Topics() {
		m.MethodCalled("Publish", topic, event)
	}
}

func (m *MockPublisher[Evt]) ExpectPublish(topic string, payload ...any) *mock.Call {
	payloadD := internal.ParseOptionalParam(payload, mock.Anything)
	return m.On("Publish", topic, payloadD)
}

func NewMockPublisher() *MockPublisher[events.Event] {
	return &MockPublisher[events.Event]{}
}

func NewCustomMockPublisher[Evt events.Event]() *MockPublisher[Evt] {
	return &MockPublisher[Evt]{}
}

var _ events.AppBus = (*MockBus[events.Event])(nil)

type MockBus_Expecter struct {
	mock *mock.Mock
}

func (m *MockBus_Expecter) Subscribe(subscription any, opts ...any) *mock.Call {
	args := append([]any{subscription}, opts...)
	return m.mock.On("Subscribe", args...)
}

func (m *MockBus_Expecter) Publish(event any, actor any) *mock.Call {
	return m.mock.On("Publish", event, actor)
}

func (m *MockBus_Expecter) PublishSync(subscription any, actor any) *mock.Call {
	// In mocking we do not distinguish between sync and async
	return m.Publish(subscription, actor)
}

type MockAppBus = MockBus[events.Event]

type MockBus[Evt events.Event] struct {
	mock.Mock
	listeners map[events.Subscription]events.BusListener[events.Event]
}

func (m *MockBus[Evt]) Subscribe(subscription events.BusListener[events.Event], opts ...events.SubscribeOpt) events.Subscription {
	callArgs := []any{subscription}
	for range opts {
		callArgs = append(callArgs, mock.Anything)
	}
	args := m.Called(callArgs...)
	subscriber := args.Get(0).(events.Subscription)
	m.listeners[subscriber] = subscription
	return subscriber
}

func (m *MockBus[Evt]) Publish(event events.Event, actor ...flamigo.Actor) {
	actorD := internal.ParseOptionalParam(actor, nil)
	m.Called(event, actorD)
}

func (m *MockBus[Evt]) PublishSync(event events.Event, actor ...flamigo.Actor) {
	// In mocking we do not distinguish between sync and async
	m.Publish(event, actor...)
}

func (m *MockBus[Evt]) TRIGGER(subscription events.Subscription, ctx events.Context, event events.Event) error {
	if listener, ok := m.listeners[subscription]; ok {
		listener(ctx, event)
		return nil
	}
	return errors.New("suscription not yet registered an listener")
}

func (m *MockBus[Evt]) EXPECT() *MockBus_Expecter {
	return &MockBus_Expecter{&m.Mock}
}

func NewBus() *MockBus[events.Event] {
	return &MockBus[events.Event]{
		listeners: make(map[events.Subscription]events.BusListener[events.Event]),
	}
}
