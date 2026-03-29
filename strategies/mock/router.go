package mock_strategies

import (
	"github.com/amberbyte/flamigo/strategies"
	"github.com/stretchr/testify/mock"
)

type MockRouter struct {
	MockRegistry
}

var _ strategies.AppRouter = (*MockRouter)(nil)

func (m *MockRouter) Register(namespace string, registry strategies.AppRegistry) error {
	args := m.MethodCalled("Register", namespace, registry)
	return args.Error(0)
}

func (m *MockRouter) Invoke(ctx strategies.Context) strategies.Result {
	args := m.MethodCalled("Invoke", ctx)
	return args.Get(0).(strategies.Result)
}

type MockRouterExpecter struct {
	m *MockRouter
}

func (m *MockRouter) EXPECT_ROUTER() *MockRouterExpecter {
	return &MockRouterExpecter{m: m}
}

func (m *MockRouterExpecter) Register(namespace any, registry any) *mock.Call {
	return m.m.On("Register", namespace, registry)
}

func (m *MockRouterExpecter) Invoke(ctx any) *mock.Call {
	return m.m.On("Invoke", ctx)
}

func NewRouter() *MockRouter {
	return &MockRouter{
		MockRegistry: MockRegistry{
			strategies: make(map[string]strategies.AppStrategy),
		},
	}
}
