package realtime

import (
	"context"
	"sync"
	"time"

	flamigo "github.com/amberbyte/flamigo/core"
	"github.com/amberbyte/flamigo/internal"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

var _ AppBus = (*Bus[Event])(nil)

type published[T Event] struct {
	event  T
	actor  flamigo.Actor
	syncWg *sync.WaitGroup
}

func (p published[T]) Done() {
	if p.syncWg != nil {
		p.syncWg.Done()
	}
}

type Subscription[T Event] struct {
	id     string
	topics map[string]bool

	all            bool
	onlyReceivales bool
	channel        chan published[T]
	ended          bool
	topicsLock     sync.Mutex
}

func (s *Subscription[T]) Cancel() {
	if s.ended {
		return
	}
	close(s.channel)
	s.ended = true
}

func (s *Subscription[T]) SubscribeTopic(topic string) {
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

func (s *Subscription[T]) UnsubscribeTopic(topic string) {
	if s.topics == nil {
		return
	}
	s.topicsLock.Lock()
	defer s.topicsLock.Unlock()
	delete(s.topics, topic)
}

func (s *Subscription[T]) OnlyReceivables() {
	s.onlyReceivales = true
}

// func (s *Subscription[T]) HasTopic(topic string) bool {
// 	if s.all {
// 		return true
// 	}
// 	if s.topics == nil {
// 		return false
// 	}
// 	return s.topics[topic]
// }

func (s *Subscription[T]) matchesTopic(topic Topic) bool {
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

func (s *Subscription[T]) AllTopics() {
	s.all = true
}

func newSubscription[T Event](channel chan published[T]) *Subscription[T] {
	id, _ := gonanoid.New()
	return &Subscription[T]{
		id:      id,
		channel: channel,
	}
}

type Bus[T Event] struct {
	listeners     map[string]*Subscription[T]
	listenersLock sync.RWMutex
}

func (bus *Bus[T]) getAllSubscribers(topic []string, receivable bool) []*Subscription[T] {
	bus.listenersLock.RLock()
	defer bus.listenersLock.RUnlock()
	// Find subscribers from topic
	subscribers := make([]*Subscription[T], 0)
	for _, sub := range bus.listeners {
		if sub.onlyReceivales && !receivable {
			continue
		}
		if sub.matchesTopic(topic) {
			subscribers = append(subscribers, sub)
		}
	}
	return subscribers
}

func (bus *Bus[T]) removeListener(id string) {
	// Request write lock for listeners
	bus.listenersLock.Lock()
	defer bus.listenersLock.Unlock()

	// Remove listener
	delete(bus.listeners, id)
}

func (bus *Bus[T]) addListener(subscription *Subscription[T]) {
	// Request write lock for listeners
	bus.listenersLock.Lock()
	defer bus.listenersLock.Unlock()

	// Add listener
	bus.listeners[subscription.id] = subscription
}

func (bus *Bus[T]) Subscribe(listener BusListener[T]) *Subscription[T] {
	// Create subscription
	channel := make(chan published[T], 200)
	subscription := newSubscription[T](channel)
	bus.addListener(subscription)

	// Start reading loop
	go func() {
		for {
			message, ok := <-channel
			if !ok {
				break
			}
			// We go concurrent here to avoid blocking further processing of events if one listener take some time to process
			go func(msg published[T]) {
				rctx, c := context.WithTimeout(context.Background(), 5*time.Second)
				defer c()
				appCtx := flamigo.NewContext(rctx, message.actor)
				ctx := NewContext(appCtx)

				listener(ctx, message.event)
				message.Done()
			}(message)
		}

		// Once the reading loop ends, delete the subscription
		bus.removeListener(subscription.id)
	}()
	return subscription
}

func (bus *Bus[T]) Publish(message T, actor ...flamigo.Actor) {
	normActor := internal.ParseOptionalParam[flamigo.Actor](actor, flamigo.NewServerActor("unknown"))
	alreadyReceived := make(map[string]bool)
	for _, topic := range message.Topics() {
		isReceivable := IsClientEvent(message)
		subscribers := bus.getAllSubscribers(topic, isReceivable)
		if len(subscribers) == 0 {
			continue
		}
		for _, subscription := range subscribers {
			if alreadyReceived[subscription.id] {
				continue
			}
			alreadyReceived[subscription.id] = true
			subscription.channel <- published[T]{event: message, actor: normActor}
		}
	}
}

// PublishSync publishes a message to all subscribers and blocks until all subscribers have processed the message
func (bus *Bus[T]) PublishSync(message T, actor ...flamigo.Actor) {
	normActor := internal.ParseOptionalParam[flamigo.Actor](actor, flamigo.NewServerActor("unknown"))
	alreadyReceived := make(map[string]bool)
	wg := sync.WaitGroup{}
	for _, topic := range message.Topics() {
		isReceivable := IsClientEvent(message)
		subscribers := bus.getAllSubscribers(topic, isReceivable)
		if len(subscribers) == 0 {
			continue
		}

		for _, subscription := range subscribers {
			if alreadyReceived[subscription.id] {
				continue
			}
			wg.Add(1)
			alreadyReceived[subscription.id] = true
			subscription.channel <- published[T]{event: message, actor: normActor, syncWg: &wg}
		}
	}
	wg.Wait()
}

func NewBus[T Event]() *Bus[T] {
	return &Bus[T]{listeners: make(map[string]*Subscription[T], 0)}
}
