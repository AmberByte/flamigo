# Auth

Authentication in Flamigo is built as a dedicated domain, responsible for handling user identities and login state.  
When creating a new project, **auth must be explicitly enabled** to activate this functionality.

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