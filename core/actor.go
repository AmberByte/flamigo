package flamigo

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrNoUserActor      = errors.New("actor is not a user actor")
	ErrNotAuthenticated = errors.New("actor is not authenticated")
	ErrNoPermission     = errors.New("actor does not have permission")
	ErrAuthenticated    = errors.New("actor is authenticated")
	ErrInvalidActorType = errors.New("actor does not fullfill the required type")
)

const (
	TypeActorServer = "backend"
)

type Actor interface {
	Type() string
}

type ActorClaimValidator = func(ctx context.Context, actor Actor) error

func OfType(requiredType string) ActorClaimValidator {
	return func(ctx context.Context, actor Actor) error {
		if actor.Type() != requiredType {
			return ErrInvalidActorType
		}
		return nil
	}
}

func IsServer() ActorClaimValidator {
	return func(ctx context.Context, actor Actor) error {
		if sA, ok := actor.(Actor); ok {
			if sA.Type() != TypeActorServer {
				return ErrInvalidActorType
			}
			return nil
		}
		return ErrInvalidActorType
	}
}

func IsActorOfType(interfaceName string) ActorClaimValidator {
	return func(ctx context.Context, actor Actor) error {
		if sA, ok := actor.(Actor); ok {
			if sA.Type() != interfaceName {
				return ErrInvalidActorType
			}
			return nil
		}
		return ErrInvalidActorType
	}
}

func wrapRequireActorErr(err error) error {
	return fmt.Errorf("validating actor: %w", err)
}

func RequireActorWithClaims[T Actor](ctx Context, opts ...ActorClaimValidator) (parsedActor T, err error) {
	actor := ctx.Actor()
	for _, opt := range opts {
		if err = opt(ctx, actor); err != nil {
			err = wrapRequireActorErr(err)
			return
		}
	}

	parsed, ok := actor.(T)
	if !ok {
		err = wrapRequireActorErr(ErrInvalidActorType)
		return
	}
	parsedActor = parsed
	return
}

func VerifyActorClaims(ctx Context, opts ...ActorClaimValidator) error {
	_, err := RequireActorWithClaims[Actor](ctx, opts...)
	return err
}
