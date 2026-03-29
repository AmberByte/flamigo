package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildClientTopic(t *testing.T) {
	topic := BuildClientTopic("topic", "a", "b", "c")
	assert.Equal(t, "topic:a:b:c", topic)
}
