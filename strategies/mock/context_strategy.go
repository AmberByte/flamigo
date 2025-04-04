package mock_strategies

import (
	"context"
	"encoding/json"

	flamigo "github.com/amberbyte/flamigo/core"
	"github.com/amberbyte/flamigo/strategies"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

var _ strategies.Context = &MockContext{}

type MockContext struct {
	flamigo.Context
	mock.Mock
	request  *strategies.Request
	response *strategies.Response
}

var _ strategies.Context = &MockContext{}

func (s *MockContext) Request() *strategies.Request {
	return s.request
}

func (s *MockContext) Response() *strategies.Response {
	return s.response
}

func (s *MockContext) Logger() *logrus.Entry {
	return logrus.WithField("component", "strategy")
}

func (s *MockContext) AssertExpectations(t mock.TestingT) bool {
	return s.Mock.AssertExpectations(t)
}

// SetRequestPayload mocks the raw request payload
func (s *MockContext) SetRequestPayloadRaw(payload interface{}) error {
	s.request = strategies.NewRequest("", payload)
	return nil
}

// SetRequestPayload mocks the payload. it parses the given data as json first
func (s *MockContext) SetRequestPayload(payload interface{}) error {
	dt, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	s.request = strategies.NewRequest("", dt)
	return nil
}

// func (s *strategyMockContext) UseStrategy(key string, payload interface{}) *flamigo.Result {
// 	result := s.MethodCalled("Use", key, payload)
// 	return result.Get(0).(*flamigo.Result)
// }

func NewMockContext(actor flamigo.Actor) *MockContext {
	c := &MockContext{
		Context:  flamigo.NewContext(context.Background(), actor),
		request:  &strategies.Request{},
		response: &strategies.Response{},
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
