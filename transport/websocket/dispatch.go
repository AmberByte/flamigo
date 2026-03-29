package websocket

import (
	"context"

	flamigo "github.com/amberbyte/flamigo/core"
	"github.com/amberbyte/flamigo/strategies"
)

type Writer func(raw []byte) error
type CommandDecoder func(raw []byte) (*Command, error)
type ActorFactory func(ctx context.Context, cmd *Command) (flamigo.Actor, error)
type ErrorEncoder func(err error, cmd *Command) ([]byte, error)
type SuccessEncoder func(cmd *Command, payload any) ([]byte, error)

type DispatchOption func(*Dispatcher)

type Dispatcher struct {
	router         strategies.AppRouter
	commandDecoder CommandDecoder
	actorFactory   ActorFactory
	errorEncoder   ErrorEncoder
	successEncoder SuccessEncoder
}

func WithCommandDecoder(decoder CommandDecoder) DispatchOption {
	return func(d *Dispatcher) {
		d.commandDecoder = decoder
	}
}

func WithActorFactory(factory ActorFactory) DispatchOption {
	return func(d *Dispatcher) {
		d.actorFactory = factory
	}
}

func WithErrorEncoder(encoder ErrorEncoder) DispatchOption {
	return func(d *Dispatcher) {
		d.errorEncoder = encoder
	}
}

func WithSuccessEncoder(encoder SuccessEncoder) DispatchOption {
	return func(d *Dispatcher) {
		d.successEncoder = encoder
	}
}

func defaultActorFactory(ctx context.Context, cmd *Command) (flamigo.Actor, error) {
	return flamigo.NewServerActor("websocket"), nil
}

func defaultErrorEncoder(err error, cmd *Command) ([]byte, error) {
	if cmd == nil {
		return EncodeError(err)
	}
	return EncodeError(err, FromCommand(cmd))
}

func defaultSuccessEncoder(cmd *Command, payload any) ([]byte, error) {
	return EncodeSuccess(cmd.Action(), payload, FromCommand(cmd))
}

func NewDispatcher(router strategies.AppRouter, opts ...DispatchOption) *Dispatcher {
	dispatcher := &Dispatcher{
		router:         router,
		commandDecoder: DecodeCommand,
		actorFactory:   defaultActorFactory,
		errorEncoder:   defaultErrorEncoder,
		successEncoder: defaultSuccessEncoder,
	}
	for _, opt := range opts {
		opt(dispatcher)
	}
	return dispatcher
}

func (d *Dispatcher) HandleMessage(ctx context.Context, raw []byte, write Writer) error {
	cmd, err := d.commandDecoder(raw)
	if err != nil {
		encoded, encodeErr := d.errorEncoder(err, nil)
		return d.writeEncoded(write, encoded, encodeErr)
	}

	actor, err := d.actorFactory(ctx, cmd)
	if err != nil {
		encoded, encodeErr := d.errorEncoder(err, cmd)
		return d.writeEncoded(write, encoded, encodeErr)
	}

	appCtx := flamigo.NewContext(ctx, actor)
	strategyCtx := strategies.NewContext(appCtx, cmd.Action(), cmd.Payload())
	AttachMetadata(strategyCtx.Request(), Metadata{
		AckKey: cmd.AckKey(),
	})

	result := d.router.Invoke(strategyCtx)
	if err := result.Err(); err != nil {
		encoded, encodeErr := d.errorEncoder(err, cmd)
		return d.writeEncoded(write, encoded, encodeErr)
	}
	if !result.IsOk() {
		return nil
	}

	encoded, encodeErr := d.successEncoder(cmd, result.Payload())
	return d.writeEncoded(write, encoded, encodeErr)
}

func (d *Dispatcher) writeEncoded(write Writer, raw []byte, err error) error {
	if err != nil {
		return err
	}
	return write(raw)
}
