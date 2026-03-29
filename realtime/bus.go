package realtime

import (
	"context"
	"sync"
	"time"

	flamigo "github.com/amberbyte/flamigo/core"
	"github.com/amberbyte/flamigo/internal"
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
				ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				appCtx := flamigo.NewContext(ctxWithTimeout, msg.actor)
				listener(NewContext(appCtx), msg.event)
				cancel()
				msg.Done()
			}
		}
	}()

	return subscription
}

func (b *bus[T]) Publish(message T, actor ...flamigo.Actor) {
	normActor := internal.ParseOptionalParam[flamigo.Actor](actor, flamigo.NewServerActor("unknown"))
	alreadyReceived := make(map[string]bool)

	for _, topic := range message.Topics() {
		subscribers := b.getAllSubscribers(topic)

		for _, sub := range subscribers {
			if alreadyReceived[sub.id] {
				continue
			}
			alreadyReceived[sub.id] = true

			sub.publish(published[T]{event: message, actor: normActor})
		}
	}
}

func (b *bus[T]) PublishSync(message T, actor ...flamigo.Actor) {
	normActor := internal.ParseOptionalParam(actor, flamigo.NewServerActor("unknown"))
	alreadyReceived := make(map[string]bool)
	var wg sync.WaitGroup

	for _, topic := range message.Topics() {
		subscribers := b.getAllSubscribers(topic)

		for _, sub := range subscribers {
			if alreadyReceived[sub.id] {
				continue
			}
			alreadyReceived[sub.id] = true

			wg.Add(1)
			if !sub.publish(published[T]{event: message, actor: normActor, syncWg: &wg}) {
				wg.Done()
			}
		}
	}

	wg.Wait()
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
