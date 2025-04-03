package strategies

import "fmt"

type Response struct {
	strategyTopic string
	result        any
	rejected      error
}

func (s *Response) SetResult(payload any) {
	s.result = payload
}
func (s *Response) SetError(err error) {
	s.rejected = fmt.Errorf("call strategy(%s): %w", s.strategyTopic, err)
}

func (s *Response) IsOk() bool {
	return s.result != nil || s.result != nil && s.rejected == nil
}
func (s *Response) IsError() bool {
	return s.rejected != nil
}

func (s *Response) Result() any {
	return s.result
}
func (s *Response) Err() error {
	return s.rejected
}
