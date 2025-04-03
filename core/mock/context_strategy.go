package mock_flamigo

import (
	"context"

	flamigo "github.com/amberbyte/flamigo/core"
	"github.com/stretchr/testify/mock"
)

type MockContext struct {
	context.Context
	mock.Mock
	actor flamigo.Actor
}

var _ flamigo.Context = &MockContext{}

func (s *MockContext) Action() string {
	result := s.MethodCalled("Action")
	return result.Get(0).(string)
}

func (s *MockContext) Actor() flamigo.Actor {
	return s.actor
}

func (s *MockContext) AssertExpectations(t mock.TestingT) bool {
	return s.Mock.AssertExpectations(t)
}

// func (s *strategyMockContext) UseStrategy(key string, payload interface{}) *flamigo.Result {
// 	result := s.MethodCalled("Use", key, payload)
// 	return result.Get(0).(*flamigo.Result)
// }

func NewMockContext(actor flamigo.Actor) *MockContext {
	c := &MockContext{
		Context: context.Background(),
		actor:   actor,
	}
	return c
}

// func NewMockInternalContext() *MockContext {
// 	c := &MockContext{
// 		Context: context.Background(),
// 		actor:   &InternalActorMock{},
// 		request: &MockRequest{},
// 	}
// 	return c
// }
