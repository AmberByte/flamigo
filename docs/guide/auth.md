# Auth

Authentication in Flamigo is built as a dedicated domain, responsible for handling user identities and login state.  

::: warning
The config module can only be enabled when initializing a new project.
:::

The auth domain provides:

- A `User` model
- A specialized `UserActor`
- Built-in actor claim validation
- Hooks to customize authentication logic for your app

You are expected to **customize the authentication logic** to fit your application's specific needs — whether that's integrating OAuth, JWTs, session tokens, or any other identity mechanism.

---

## The `UserActor`

The `UserActor` is a special kind of actor provided by the auth domain. It extends the core `flamigo.Actor` interface and includes user identity information:

```go
type UserActor interface {
  flamigo.Actor
  User() *User
}
```

This allows any service, strategy, or listener to easily access the current user and perform user-specific logic or authorization checks.

---

### Validating `UserActor` with Claims

To simplify the use of `UserActor` in strategies or other parts of the application, the auth domain provides a helper function:

```go
func RequireUserActorWithClaims(ctx flamigo.Context, opts ...flamigo.ActorClaimValidator) (UserActor, error)
```

This function behaves similarly to `flamigo.RequireActorWithClaims`, but will **only succeed if the actor is a `UserActor`**. It's perfect for enforcing user-based permissions or ensuring a valid session before executing logic.

Example:

```go
userActor, err := auth.RequireUserActorWithClaims(ctx, auth.IsAuthenticated())
if err != nil {
  // handle unauthenticated access
}
```

---

## Actor Claim Validators

The auth domain also introduces **claim validators**, which are pluggable checks that can be used to verify actor properties before running your logic.

### `IsAuthenticated`

Checks whether the actor is a logged-in user.

```go
actor, err := auth.RequireUserActorWithClaims(ctx, auth.IsAuthenticated())
```

---

### `IsUnauthenticated`

Checks whether the actor is not logged in — useful for routes like login, registration, or public actions.

```go
actor, err := auth.RequireUserActorWithClaims(ctx, auth.IsUnauthenticated())
```

---

## Integrating with OAuth2

To integrate OAuth2, a login service could look like the following

```go
func OAuth2Login(ctx flamigo.Context, token string) (*auth.User, error) {
    userInfo, err := oauth2.ValidateToken(token)
    if err != nil {
        return nil, flamigo.NewError("invalid token", flamigo.WithPublicResponse("Authentication failed."))
    }

    user := &auth.User{
        ID:    userInfo.ID,
        Email: userInfo.Email,
    }
    return user, nil
}
```

---

## Integrating

To check if a user is validated, you can create ActorClaimValidators

```go
func IsAuthenticated(ctx flamigo.Context, token string) flamigo.ActorClaimValidator {
   return func(ctx context.Context, actor flamigo.Actor) error {
		if uA, ok := actor.(UserActor); ok {
			if uA.User() != nil {
				return auth.ErrAuthenticated
			}
			return nil
		}
		return auth.ErrNoUserActor
	}
}
```

This can then be used together with `RequireUserActorWithClaims`
```go
func MyStrategy(ctx strategies.Context) error {
  actor, err := auth.RequireUserActorWithClaims(ctx, auth.IsAuthenticated())
}
```