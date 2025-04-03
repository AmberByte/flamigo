package strategies_test

import (
	"testing"

	"github.com/amberbyte/flamigo/strategies"
	"github.com/stretchr/testify/assert"
)

func TestResultResolve(t *testing.T) {
	req := &strategies.Response{}
	req.SetResult("test")
	assert.True(t, req.IsOk(), "should be resolved")
	assert.False(t, req.IsError(), "should not be rejected")

	assert.Equal(t, "test", req.Result(), "should have resolved with 'test'")
}

func TestResultReject(t *testing.T) {
	req := &strategies.Response{}
	req.SetError(assert.AnError)
	assert.False(t, req.IsOk(), "should not be resolved")
	assert.True(t, req.IsError(), "should be rejected")

	assert.Equal(t, "call strategy(): "+assert.AnError.Error(), req.Err().Error(), "should have rejected with ErrStrategyNotFound")
}
