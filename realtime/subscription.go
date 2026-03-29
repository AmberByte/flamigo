package realtime

import (
	"sync"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Subscription interface {
	Cancel()
	SubscribeTopic(topic string)
	UnsubscribeTopic(topic string)
	SubscribeAll()
}

var _ Subscription = (*subscription[Event])(nil)

type subscription[T Event] struct {
	id     string
	topics map[string]bool

	all        bool
	channel    chan published[T]
	done       chan struct{}
	cancelOnce sync.Once
	onCancel   func()
	topicsLock sync.RWMutex
}

// Cancel ends the subscription and unregisters it from the bus.
func (s *subscription[T]) Cancel() {
	s.cancelOnce.Do(func() {
		close(s.done)
		if s.onCancel != nil {
			s.onCancel()
		}
	})
}

// SubscribeTopic adds a topic to the subscription. If the subscription is already set to all topics, this will have no effect.
func (s *subscription[T]) SubscribeTopic(topic string) {
	s.topicsLock.Lock()
	defer s.topicsLock.Unlock()

	if s.all {
		return
	}
	if s.topics == nil {
		s.topics = make(map[string]bool)
	}
	s.topics[topic] = true
}

// UnsubscribeTopic removes a topic from the subscription. If the subscription is already set to all topics, this will have no effect.
func (s *subscription[T]) UnsubscribeTopic(topic string) {
	s.topicsLock.Lock()
	defer s.topicsLock.Unlock()

	if s.all {
		return
	}
	if s.topics == nil {
		return
	}
	delete(s.topics, topic)
}

func (s *subscription[T]) matchesTopic(topic Topic) bool {
	s.topicsLock.RLock()
	defer s.topicsLock.RUnlock()

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
	s.topicsLock.Lock()
	defer s.topicsLock.Unlock()
	s.all = true
	s.topics = nil
}

func (s *subscription[T]) publish(msg published[T]) bool {
	select {
	case <-s.done:
		return false
	case s.channel <- msg:
		return true
	}
}

func newSubscription[T Event](channel chan published[T], onCancel func()) *subscription[T] {
	id, _ := gonanoid.New()
	return &subscription[T]{
		id:       id,
		channel:  channel,
		done:     make(chan struct{}),
		onCancel: onCancel,
	}
}
