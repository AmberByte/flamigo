package mock_strategies

import (
	"errors"

	"github.com/amberbyte/flamigo/strategies"
	"github.com/stretchr/testify/mock"
)

var _ strategies.AppRegistry = (*MockRegistry)(nil)

type MockRegistry struct {
	mock.Mock
	strategies map[string]strategies.AppStrategy
}

func (m *MockRegistry) Register(topic string, fn strategies.AppStrategy) error {
	args := m.MethodCalled("Register", topic, fn)
	m.strategies[topic] = fn
	return args.Error(0)
}

func (m *MockRegistry) Invoke(action string, ctx strategies.Context) strategies.Result {
	args := m.MethodCalled("Invoke", action, ctx)
	return args.Get(0).(strategies.Result)

}

// TestCallStrategy is a helper method to test the strategy. It directly returns the response and a possible error for easier evaluataion
func (m *MockRegistry) TestCallStrategy(topic string, ctx strategies.Context) (any, error) {
	fn, ok := m.strategies[topic]
	if !ok {
		return nil, errors.New("strategy not found")
	}
	fn(ctx)
	return ctx.Response().Payload(), ctx.Response().Err()
}

type MockRegistryExpecter struct {
	m *mock.Mock
}

func (m *MockRegistry) EXPECT() *MockRegistryExpecter {
	return &MockRegistryExpecter{m: &m.Mock}
}

func (m *MockRegistryExpecter) Register(topic any, fn any) *mock.Call {
	return m.m.On("Register", topic, fn)
}

func (m *MockRegistryExpecter) Invoke(action any, ctx any) *mock.Call {
	return m.m.On("Invoke", action, ctx)
}

func NewRegistry() *MockRegistry {
	return &MockRegistry{
		strategies: make(map[string]strategies.AppStrategy),
	}
}
