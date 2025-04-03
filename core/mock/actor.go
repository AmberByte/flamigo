package mock_flamigo

import (
	flamigo "github.com/amberbyte/flamigo/core"
	"github.com/stretchr/testify/mock"
)

type InternalActorMock struct {
	mock.Mock
}

func (a *InternalActorMock) Type() string {
	return flamigo.TypeActorServer
}

func (a *InternalActorMock) Logger(map[string]any) {
	// Since this is a mock it is not added to the logger
}

func NewMockServerActor() *InternalActorMock {
	return &InternalActorMock{}
}
