package mock_strategies

import (
	"github.com/amberbyte/flamigo/internal"
	"github.com/stretchr/testify/mock"
)

type MockRequest struct {
	mock.Mock
}

func (s *MockRequest) Resolve(payload interface{}) {
	s.MethodCalled("Resolve", payload)

}
func (s *MockRequest) Reject(err error) {
	s.MethodCalled("Reject", err)
}

func (s *MockRequest) ExpectResolved(resolvePayload ...interface{}) *mock.Call {
	resolvePayloadDefaulted := internal.ParseOptionalParam[interface{}](resolvePayload, mock.Anything)
	return s.On("Resolve", resolvePayloadDefaulted)
}

func (s *MockRequest) ExpectRejected(err ...interface{}) *mock.Call {
	errDefaulted := internal.ParseOptionalParam[interface{}](err, mock.Anything)
	return s.On("Reject", errDefaulted)
}

func (s *MockRequest) IsResolved() bool {
	result := s.MethodCalled("IsResolved")
	return result.Get(0).(bool)
}
func (s *MockRequest) IsRejected() bool {
	result := s.MethodCalled("IsRejected")
	return result.Get(0).(bool)
}

func (s *MockRequest) Resolved() interface{} {
	result := s.MethodCalled("Resolved")
	return result.Get(0)
}
func (s *MockRequest) Rejected() error {
	result := s.MethodCalled("Rejected")
	return result.Get(0).(error)
}

func NewMockRequest() *MockRequest {
	return &MockRequest{}
}
