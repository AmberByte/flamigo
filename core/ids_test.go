package flamigo_test

import (
	"testing"

	flamigo "github.com/amberbyte/flamigo/core"
	"github.com/go-playground/assert/v2"
)

func TestNewRandomID(t *testing.T) {
	id1 := flamigo.NewRandomID()
	id2 := flamigo.NewRandomID()
	id3 := flamigo.NewRandomID()
	// Make sure all ids are unique
	assert.NotEqual(t, id1, id2)
	assert.NotEqual(t, id1, id3)
	assert.NotEqual(t, id2, id3)
}

func TestNewHashId(t *testing.T) {
	id1 := flamigo.NewHashId("test")
	id2 := flamigo.NewHashId("test")
	id3 := flamigo.NewHashId("test_other")
	// Make sure all ids are unique
	assert.Equal(t, id1, id2)
	assert.NotEqual(t, id1, id3)
	assert.NotEqual(t, id2, id3)
}
