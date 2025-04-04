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
	Subscribe(listener BusListener[T]) Subscription
	Publish(message T, actor ...flamigo.Actor)
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

func (b *bus[T]) getAllSubscribers(topic []string, receivable bool) []*subscription[T] {
	b.listenersLock.RLock()
	defer b.listenersLock.RUnlock()

	var subscribers []*subscription[T]
	for _, sub := range b.listeners {
		if sub.onlyClientMessages && !receivable {
			continue
		}
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

func (b *bus[T]) Subscribe(listener BusListener[T]) Subscription {
	channel := make(chan published[T], b.bufferSize)
	subscription := newSubscription[T](channel)
	b.addListener(subscription)

	go func() {
		defer b.removeListener(subscription.id)
		for msg := range channel {
			go func(msg published[T]) {
				ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				appCtx := flamigo.NewContext(ctxWithTimeout, msg.actor)
				listener(NewContext(appCtx), msg.event)
				msg.Done()
			}(msg)
		}
	}()

	return subscription
}

func (b *bus[T]) Publish(message T, actor ...flamigo.Actor) {
	normActor := internal.ParseOptionalParam[flamigo.Actor](actor, flamigo.NewServerActor("unknown"))
	alreadyReceived := make(map[string]bool)

	for _, topic := range message.Topics() {
		isReceivable := IsClientEvent(message)
		subscribers := b.getAllSubscribers(topic, isReceivable)

		for _, sub := range subscribers {
			if alreadyReceived[sub.id] {
				continue
			}
			alreadyReceived[sub.id] = true

			select {
			case sub.channel <- published[T]{event: message, actor: normActor}:
				// Successfully published
			default:
				// Channel full; could log, drop, or handle differently
			}
		}
	}
}

func (b *bus[T]) PublishSync(message T, actor ...flamigo.Actor) {
	normActor := internal.ParseOptionalParam[flamigo.Actor](actor, flamigo.NewServerActor("unknown"))
	alreadyReceived := make(map[string]bool)
	var wg sync.WaitGroup

	for _, topic := range message.Topics() {
		isReceivable := IsClientEvent(message)
		subscribers := b.getAllSubscribers(topic, isReceivable)

		for _, sub := range subscribers {
			if alreadyReceived[sub.id] {
				continue
			}
			alreadyReceived[sub.id] = true

			wg.Add(1)
			select {
			case sub.channel <- published[T]{event: message, actor: normActor, syncWg: &wg}:
			default:
				wg.Done() // If full, drop and still mark as done
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
