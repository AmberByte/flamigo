package flamigo_test

import (
	"testing"
)

func TestIsAuthenticated(t *testing.T) {
	// t.Run("Returns error when actor is not a user actor", func(t *testing.T) {
	// 	actor := mock_flamigo.NewMockServerActor()
	// 	err := flamigo.IsAuthenticated()(nil, actor)
	// 	assert.Equal(t, flamigo.ErrNoUserActor, err)
	// })

	// t.Run("Returns error when actor is not authenticated", func(t *testing.T) {
	// 	actor := mock_flamigo.NewMockWebsocketActor(false)
	// 	err := flamigo.IsAuthenticated()(nil, actor)
	// 	assert.Equal(t, flamigo.ErrNotAuthenticated, err)
	// })

	// t.Run("Returns nil when actor is authenticated", func(t *testing.T) {
	// 	actor := mock_flamigo.NewMockWebsocketActor(true)
	// 	err := flamigo.IsAuthenticated()(nil, actor)
	// 	assert.Nil(t, err)
	// })
}

func TestIsUnauthenticated(t *testing.T) {
	// t.Run("Returns error when actor is not a user actor", func(t *testing.T) {
	// 	actor := mock_flamigo.NewMockServerActor()
	// 	err := flamigo.IsUnauthenticated()(nil, actor)
	// 	assert.Equal(t, flamigo.ErrNoUserActor, err)
	// })

	// t.Run("Returns error when actor is authenticated", func(t *testing.T) {
	// 	actor := mock_flamigo.NewMockWebsocketActor(true)
	// 	err := flamigo.IsUnauthenticated()(nil, actor)
	// 	assert.Equal(t, flamigo.ErrAuthenticated, err)
	// })

	// t.Run("Returns nil when actor is not authenticated", func(t *testing.T) {
	// 	actor := mock_flamigo.NewMockWebsocketActor(false)
	// 	err := flamigo.IsUnauthenticated()(nil, actor)
	// 	assert.Nil(t, err)
	// })

}

func TestOfType(t *testing.T) {
	// t.Run("Returns error when actor type does not match", func(t *testing.T) {
	// 	actor := mock_flamigo.NewMockServerActor()
	// 	err := flamigo.OfType("mock")(nil, actor)
	// 	assert.Equal(t, flamigo.ErrInvalidActorType, err)
	// })

	// t.Run("Returns nil when actor type matches", func(t *testing.T) {
	// 	actor := mock_flamigo.NewMockWebsocketActor()
	// 	err := flamigo.OfType("websocket")(nil, actor)
	// 	assert.Nil(t, err)
	// })
}

func TestRequireActor(t *testing.T) {
	// t.Run("Returns error when actor canot be parsed", func(t *testing.T) {

	// 	actor := mock_flamigo.NewMockServerActor()
	// 	_, err := flamigo.RequireActorWithClaims[flamigo.UserActor](mock_flamigo.NewMockContext(actor))
	// 	assert.Equal(t, "validating actor: "+flamigo.ErrInvalidActorType.Error(), err.Error())
	// })

	// t.Run("Returns error of modifier", func(t *testing.T) {
	// 	actor := mock_flamigo.NewMockWebsocketActor(false)
	// 	_, err := flamigo.RequireActorWithClaims[flamigo.UserActor](mock_flamigo.NewMockContext(actor), flamigo.IsAuthenticated())
	// 	assert.Equal(t, "validating actor: "+flamigo.ErrNotAuthenticated.Error(), err.Error())
	// })

	// t.Run("Returns error of first modifier", func(t *testing.T) {
	// 	actor := mock_flamigo.NewMockWebsocketActor(false)
	// 	_, err := flamigo.RequireActorWithClaims[flamigo.UserActor](mock_flamigo.NewMockContext(actor), flamigo.OfType("mock"), flamigo.IsAuthenticated())
	// 	assert.Equal(t, "validating actor: "+flamigo.ErrInvalidActorType.Error(), err.Error())
	// })
}

func Test_IsAdmin(t *testing.T) {
	// t.Run("Returns error when actor is not a user actor", func(t *testing.T) {
	// 	actor := mock_flamigo.NewMockServerActor()
	// 	err := flamigo.IsAdmin()(nil, actor)
	// 	assert.Equal(t, flamigo.ErrNoUserActor, err)
	// })
	// t.Run("Returns error when actor is not authenticated", func(t *testing.T) {
	// 	actor := mock_flamigo.NewMockWebsocketActor(false)
	// 	err := flamigo.IsAdmin()(nil, actor)
	// 	assert.Equal(t, flamigo.ErrNotAuthenticated, err)
	// })
	// t.Run("Returns error when actor is not an admin", func(t *testing.T) {
	// 	actor := mock_flamigo.NewMockWebsocketActor(true)
	// 	err := flamigo.IsAdmin()(nil, actor)
	// 	assert.Equal(t, flamigo.ErrNoPermission, err)
	// })
	// t.Run("Returns nil when actor is an admin", func(t *testing.T) {
	// 	actor := mock_flamigo.NewMockWebsocketActor(true)
	// 	actor.User().Permissions = append(actor.User().Permissions, "game:admin")
	// 	err := flamigo.IsAdmin()(nil, actor)
	// 	assert.Nil(t, err)
	// })
}
