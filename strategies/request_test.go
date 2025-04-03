package strategies_test

import (
	"encoding/json"
	"testing"

	"github.com/amberbyte/flamigo/strategies"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalPayload(t *testing.T) {
	t.Run("it unmarshals a JSON []byte payload", func(t *testing.T) {
		data := []byte(`{"foo":"bar"}`)
		ctx := strategies.NewRequest("mock", data)
		var target map[string]string
		err := ctx.Bind(&target)
		assert.Nil(t, err)
		assert.Equal(t, "bar", target["foo"])
	})
	t.Run("it unmarshals a JSON string payload", func(t *testing.T) {
		data := `{"foo":"bar"}`
		ctx := strategies.NewRequest("mock", data)
		var target map[string]string
		err := ctx.Bind(&target)
		assert.Nil(t, err)
		assert.Equal(t, "bar", target["foo"])
	})
	t.Run("it unmarshals a JSON json.RawMessage payload", func(t *testing.T) {
		data := json.RawMessage(`{"foo":"bar"}`)
		ctx := strategies.NewRequest("mock", data)
		var target map[string]string
		err := ctx.Bind(&target)
		assert.Nil(t, err)
		assert.Equal(t, "bar", target["foo"])
	})
	t.Run("it returns an error if the payload is not a valid JSON", func(t *testing.T) {
		data := []byte(`{"foo":"bar"`)
		ctx := strategies.NewRequest("mock", data)
		var target map[string]string
		err := ctx.Bind(&target)
		assert.NotNil(t, err)
	})
	t.Run("it returns an error if the payload is not a pointer", func(t *testing.T) {
		data := []byte(`{"foo":"bar"}`)
		ctx := strategies.NewRequest("mock", data)
		var target map[string]string
		err := ctx.Bind(target)
		assert.NotNil(t, err)
	})
}
