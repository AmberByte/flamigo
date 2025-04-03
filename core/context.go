package flamigo

import (
	"context"
)

type Context interface {
	context.Context
	Actor() Actor
}

type flamigoContext struct {
	context.Context
	actor Actor
}

func (c *flamigoContext) Actor() Actor {
	return c.actor
}

func NewContext(ctx context.Context, actor Actor) Context {
	sCtx := &flamigoContext{
		Context: ctx,
		actor:   actor,
	}

	return sCtx
}
