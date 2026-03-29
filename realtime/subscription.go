package realtime

import (
	"sync"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Subscription interface {
	Cancel()
	AddTopic(topic string)
	RemoveTopic(topic string)
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

func (s *subscription[T]) AddTopic(topic string) {
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

func (s *subscription[T]) RemoveTopic(topic string) {
	s.topicsLock.Lock()
	defer s.topicsLock.Unlock()

	if s.all || s.topics == nil {
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

func newSubscription[T Event](channel chan published[T], onCancel func(), cfg subscribeConfig) *subscription[T] {
	id, _ := gonanoid.New()
	topics := make(map[string]bool, len(cfg.topics))
	for _, topic := range cfg.topics {
		topics[topic] = true
	}
	return &subscription[T]{
		id:       id,
		topics:   topics,
		all:      cfg.all,
		channel:  channel,
		done:     make(chan struct{}),
		onCancel: onCancel,
	}
}
