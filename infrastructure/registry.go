package flamigo_infra

import (
	"errors"
	"fmt"
	"strings"

	"github.com/amberbyte/flamigo/strategies"
)

type strategyRegistry[CTX strategies.Context] struct {
	prefix     string
	strategies map[string]strategies.Strategy[CTX]
}

func (r *strategyRegistry[CTX]) Register(topic string, fn strategies.Strategy[CTX]) error {
	if !strings.HasPrefix(topic, r.prefix+"::") {
		return fmt.Errorf("adding strategy: %w", fmt.Errorf("strategy %s should be %s::%s", topic, r.prefix, topic))
	}
	r.strategies[topic] = fn
	return nil
}

func (r *strategyRegistry[CTX]) Use(ctx CTX) strategies.StrategyResult {
	topic := ctx.Request().Action()
	if !strings.HasPrefix(topic, r.prefix+"::") {
		ctx.Response().SetError(fmt.Errorf("registry (%s): %w", r.prefix, fmt.Errorf("strategy %s should be %s::%s", topic, r.prefix, topic)))
		return ctx.Response()
	}
	fn, ok := r.strategies[topic]
	if !ok {
		ctx.Response().SetError(fmt.Errorf("registry (%s): %w", r.prefix, errors.New("no handler found for "+topic)))
		return ctx.Response()
	}
	fn(ctx)
	return ctx.Response()
}

func NewRegistry[CTX strategies.Context](namespace string) strategies.Registry[CTX] {
	return &strategyRegistry[CTX]{
		prefix:     namespace,
		strategies: map[string]strategies.Strategy[CTX]{},
	}
}
