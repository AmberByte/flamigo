package realtime

import (
	"errors"
)

var (
	ErrAlreadySubscribed = errors.New("already subscribed")
)

type Forwarder[T Event] struct {
	channel           chan T
	alreadySubscribed bool
}

func (bus *Forwarder[T]) Subscribe(listener ForwarderListener[T]) (func(), error) {
	if bus.alreadySubscribed {
		return nil, ErrAlreadySubscribed
	}

	go func() {
		for {
			message, ok := <-bus.channel
			if !ok {
				return
			}
			listener(message)
		}
	}()
	bus.alreadySubscribed = true
	return func() {
		close(bus.channel)
	}, nil
}

func (forwarder *Forwarder[T]) Publish(message T) {
	forwarder.channel <- message
}

func NewForwarder[T Event]() *Forwarder[T] {
	return &Forwarder[T]{
		channel: make(chan T, 5),
	}
}
