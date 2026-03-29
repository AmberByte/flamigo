package realtime

import (
	"context"
	"sync"

	flamigo "github.com/amberbyte/flamigo/core"
)

// Bus is a generic interface for event buses.
type Bus[T Event] interface {
	// Subscribe adds a listener to the bus.
	Subscribe(listener BusListener[T], opts ...SubscribeOpt) Subscription
	// Publish publishes the event to all subscribing subscribers without waiting for them to finish.
	Publish(message T, actor ...flamigo.Actor)
	// PublishSync publishes the event to all subscribing subscribers and waits for them to finish.
	PublishSync(message T, actor ...flamigo.Actor)
}

// Default buffer size if none is provided.
const defaultBufferSize = 200

// BusOptions allows optional configuration of the bus.
type BusOptions struct {
	BufferSize int
}

// AppBus is the default Bus using the base Event type.
type AppBus = Bus[Event]

var _ Bus[Event] = (*bus[Event])(nil)

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

type bus[T Event] struct {
	listeners     map[string]*subscription[T]
	listenersLock sync.RWMutex
	bufferSize    int
}

func (b *bus[T]) getAllSubscribers(topic Topic) []*subscription[T] {
	b.listenersLock.RLock()
	defer b.listenersLock.RUnlock()

	var subscribers []*subscription[T]
	for _, sub := range b.listeners {
		if sub.matchesTopic(topic) {
			subscribers = append(subscribers, sub)
		}
	}
	return subscribers
}

func (b *bus[T]) removeListener(id string) {
	b.listenersLock.Lock()
	defer b.listenersLock.Unlock()
	delete(b.listeners, id)
}

func (b *bus[T]) addListener(subscription *subscription[T]) {
	b.listenersLock.Lock()
	defer b.listenersLock.Unlock()
	b.listeners[subscription.id] = subscription
}

func (b *bus[T]) Subscribe(listener BusListener[T], opts ...SubscribeOpt) Subscription {
	channel := make(chan published[T], b.bufferSize)
	subscription := newSubscription(channel, nil, newSubscribeConfig(opts...))
	subscription.onCancel = func() {
		b.removeListener(subscription.id)
	}
	b.addListener(subscription)

	go func() {
		for {
			select {
			case <-subscription.done:
				return
			case msg := <-channel:
				appCtx := flamigo.NewContext(context.Background(), msg.actor)
				listener(NewContext(appCtx), msg.event)
				msg.Done()
			}
		}
	}()

	return subscription
}

func (b *bus[T]) Publish(message T, actor ...flamigo.Actor) {
	normActor := normalizeActor(actor)
	topics := message.Topics()

	if len(topics) == 1 {
		for _, sub := range b.getAllSubscribers(topics[0]) {
			sub.publish(published[T]{event: message, actor: normActor})
		}
		return
	}

	alreadyReceived := make(map[*subscription[T]]struct{})

	for _, topic := range topics {
		subscribers := b.getAllSubscribers(topic)

		for _, sub := range subscribers {
			if _, ok := alreadyReceived[sub]; ok {
				continue
			}
			alreadyReceived[sub] = struct{}{}

			sub.publish(published[T]{event: message, actor: normActor})
		}
	}
}

func (b *bus[T]) PublishSync(message T, actor ...flamigo.Actor) {
	normActor := normalizeActor(actor)
	topics := message.Topics()
	var wg sync.WaitGroup

	if len(topics) == 1 {
		for _, sub := range b.getAllSubscribers(topics[0]) {
			wg.Add(1)
			if !sub.publish(published[T]{event: message, actor: normActor, syncWg: &wg}) {
				wg.Done()
			}
		}
		wg.Wait()
		return
	}

	alreadyReceived := make(map[*subscription[T]]struct{})

	for _, topic := range topics {
		subscribers := b.getAllSubscribers(topic)

		for _, sub := range subscribers {
			if _, ok := alreadyReceived[sub]; ok {
				continue
			}
			alreadyReceived[sub] = struct{}{}

			wg.Add(1)
			if !sub.publish(published[T]{event: message, actor: normActor, syncWg: &wg}) {
				wg.Done()
			}
		}
	}

	wg.Wait()
}

func normalizeActor(actor []flamigo.Actor) flamigo.Actor {
	if len(actor) == 0 || actor[0] == nil {
		return flamigo.NewServerActor("unknown")
	}
	return actor[0]
}

// NewBus creates a new Bus with optional configuration.
func NewBus[T Event](opts ...BusOptions) Bus[T] {
	bufferSize := defaultBufferSize
	if len(opts) > 0 && opts[0].BufferSize > 0 {
		bufferSize = opts[0].BufferSize
	}

	return &bus[T]{
		listeners:  make(map[string]*subscription[T]),
		bufferSize: bufferSize,
	}
}
