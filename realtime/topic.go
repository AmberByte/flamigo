package realtime

import (
	"strings"

	flamigo "github.com/amberbyte/flamigo/core"
)

func BuildClientTopic(parts ...string) string {
	return strings.Join(parts, ":")
}

type Publisher interface {
	Publish(event Event, actor ...flamigo.Actor)
}

type Topic []string

// String returns the string representation of the topic
//
// e.g. "foo/bar/baz"
func (t Topic) String() string {
	return strings.Join(t, "/")
}

// DoesMatch checks if the topic matches the given pattern
//
// The pattern can contain wildcards (*) to match any part of the topic.
func (t Topic) DoesMatch(topicPattern string) bool {
	splitedTopicPattern := strings.Split(topicPattern, "/")
	if len(t) < len(splitedTopicPattern) {
		return false
	}
	for i, part := range splitedTopicPattern {
		if part != t[i] && part != "*" {
			return false
		}
	}
	return true
}

// TopicParseFromString parses a raw topic string into a Topic
//
// e.g. "foo/bar/baz" -> Topic{"foo", "bar", "baz"}.
func TopicParseFromString(rawTopicString string) Topic {
	return Topic(strings.Split(rawTopicString, "/"))
}

// NewTopic creates a new Topic from the given parts.
func NewTopic(parts ...string) Topic {
	return Topic(parts)
}
