package strategies

import (
	"context"

	flamigo "github.com/amberbyte/flamigo/core"
)

type StrategyResult interface {
	IsOk() bool
	IsError() bool
	Result() interface{}
	Err() error
}

type Context interface {
	context.Context
	Request() *Request
	Response() *Response
}

type strategyContext struct {
	context.Context
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
		request:  &Request{action: action, payload: payload},
		response: &Response{strategyTopic: action},
	}
}
