package realtime

import flamigo "github.com/amberbyte/flamigo/core"

type Context interface {
	flamigo.Context
}

type eventContext struct {
	flamigo.Context
}

func NewContext(c flamigo.Context) Context {
	return &eventContext{
		Context: c,
	}
}
