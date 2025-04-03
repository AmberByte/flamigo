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

func (t Topic) String() string {
	return strings.Join(t, "/")
}

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

func TopicParseFromString(rawTopicString string) Topic {
	return Topic(strings.Split(rawTopicString, "/"))
}

func NewTopic(parts ...string) Topic {
	return Topic(parts)
}
