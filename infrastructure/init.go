package flamigo_infra

import (
	"github.com/amberbyte/flamigo/injection"

	"github.com/amberbyte/flamigo/realtime"
	"github.com/amberbyte/flamigo/strategies"
)

func Init(injector injection.DependencyManager) error {
	// Initialize the strategie registry
	registry := NewRegistry[strategies.Context]("app")

	// Initialize the event bus
	events := NewBus[realtime.Event]()
	if err := injector.AddInjectable(registry); err != nil {
		return err
	}

	if err := injector.AddInjectable(events); err != nil {
		return err
	}
	return nil
}
