package realtime

import (
	"sync"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Subscription interface {
	Cancel()
	SubscribeTopic(topic string)
	UnsubscribeTopic(topic string)
	OnlyClientMessages()
	SubscribeAll()
}

var _ Subscription = (*subscription[Event])(nil)

type subscription[T Event] struct {
	id     string
	topics map[string]bool

	all                bool
	onlyClientMessages bool
	channel            chan published[T]
	ended              bool
	topicsLock         sync.Mutex
}

// Cancel ends the subscription. This will close the channel and prevent any further messages from being sent to it.
func (s *subscription[T]) Cancel() {
	if s.ended {
		return
	}
	close(s.channel)
	s.ended = true
}

// SubscribeTopic adds a topic to the subscription. If the subscription is already set to all topics, this will have no effect.
func (s *subscription[T]) SubscribeTopic(topic string) {
	if s.all {
		return
	}
	if s.topics == nil {
		s.topics = make(map[string]bool)
	}
	s.topicsLock.Lock()
	defer s.topicsLock.Unlock()
	s.topics[topic] = true
}

// UnsubscribeTopic removes a topic from the subscription. If the subscription is already set to all topics, this will have no effect.
func (s *subscription[T]) UnsubscribeTopic(topic string) {
	if s.topics == nil {
		return
	}
	s.topicsLock.Lock()
	defer s.topicsLock.Unlock()
	delete(s.topics, topic)
}

// OnlyClientMessages sets the subscription to only receive messages which implement ClientMessage. This is useful when implementing a client interface
func (s *subscription[T]) OnlyClientMessages() {
	s.onlyClientMessages = true
}

func (s *subscription[T]) matchesTopic(topic Topic) bool {
	if s.all {
		return true
	}
	for t := range s.topics {
		if topic.DoesMatch(t) {
			return true
		}
	}
	return false
}

// SubscribeAll sets the subscription to receive all messages. This will override any topics set before.
func (s *subscription[T]) SubscribeAll() {
	s.all = true
}

func newSubscription[T Event](channel chan published[T]) *subscription[T] {
	id, _ := gonanoid.New()
	return &subscription[T]{
		id:      id,
		channel: channel,
	}
}
