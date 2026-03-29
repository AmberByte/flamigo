package strategies

import (
	flamigo "github.com/amberbyte/flamigo/core"
)

type Result interface {
	IsOk() bool
	IsError() bool
	Payload() interface{}
	Err() error
}

type Context interface {
	flamigo.Context
	Request() *Request
	Response() *Response
}

var _ Context = (*strategyContext)(nil)
var _ flamigo.Context = (*strategyContext)(nil)

type strategyContext struct {
	flamigo.Context
	request  *Request
	response *Response
}

func (c *strategyContext) Request() *Request {
	return c.request
}

func (c *strategyContext) Response() *Response {
	return c.response
}

func NewContext(ctx flamigo.Context, action string, payload interface{}) Context {
	return &strategyContext{
		Context:  ctx,
		request:  NewRequest(action, payload),
		response: &Response{strategyTopic: action},
	}
}
