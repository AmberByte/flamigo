package realtime

import "github.com/sirupsen/logrus"

type ForwarderListener[T Event] func(message T)
type BusListener[T Event] func(context Context, message T)

type AppListener = BusListener[Event]

func ListenerOnEvent[T Event](listener BusListener[T]) AppListener {
	return func(context Context, message Event) {
		if typed, ok := message.(T); ok {
			listener(context, typed)
			return
		}
		logrus.Debugf("unsupported event type: %T", message)
	}
}
