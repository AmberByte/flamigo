package events

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestTopicMatch(t *testing.T) {
	cases := []struct {
		topic        string
		compareTopic string
		expect       bool
	}{
		{topic: "topic", compareTopic: "topic", expect: true},
		{topic: "topic/a/b/c", compareTopic: "*", expect: false},
		{topic: "topic", compareTopic: "topic1", expect: false},
		{topic: "topic/a", compareTopic: "topic/a", expect: true},
		{topic: "topic/a", compareTopic: "topic/b", expect: false},
		{topic: "topic/a", compareTopic: "topic/*", expect: true},
		{topic: "topic/a", compareTopic: "topic/*", expect: true},
		{topic: "topic/a", compareTopic: "topic/*/*", expect: false},
		{topic: "topic/a/b", compareTopic: "topic/a/b/f", expect: false},
		{topic: "topic/a/b", compareTopic: "topic/*/*", expect: true},
		{topic: "topic/a/b", compareTopic: "topic/a", expect: false},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d Match Test", i+1), func(t *testing.T) {
			topic := ParseTopic(c.topic)
			result := topic.DoesMatch(c.compareTopic)
			assert.Equal(t, c.expect, result)
		})
	}
}

func TestTopic_String(t *testing.T) {
	topic := Topic{"topic", "a", "b", "c"}
	assert.Equal(t, "topic/a/b/c", topic.String())
}

func TestParseTopic(t *testing.T) {
	topic := ParseTopic("topic/a/b/c")
	assert.Equal(t, Topic{"topic", "a", "b", "c"}, topic)
}

func TestNewTopic(t *testing.T) {
	topic := NewTopic("topic", "a", "b", "c")
	assert.Equal(t, Topic{"topic", "a", "b", "c"}, topic)
}
