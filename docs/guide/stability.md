# Stability

Flamigo is still pre-1.0.

That means the framework is already usable, but you should still expect some API and template changes while the architecture settles through real app usage.

## What feels stable already

These parts are now structurally solid enough that I would expect only incremental changes:

- the `events` package and event bus direction
- the `strategies` split between local registries and namespaced routers
- the `injection` API and startup wiring model
- the `transport/http` and `transport/websocket` package boundaries
- the generated project shape around `domains`, `api`, and `adapters`

## What may still change before 1.0

These parts are still more likely to move as apps put pressure on the framework:

- generated template details
- exact transport adapter scaffolding
- naming and ergonomics around helpers and optional packages
- documentation structure and examples

The core architectural direction should be much more stable than the exact convenience APIs around it.

## When Flamigo is a good fit today

Flamigo is already a good fit if:

- you are comfortable with pre-1.0 libraries
- you want a domain-first, hexagonal Go backend
- you are willing to update app code when templates or convenience APIs improve

You may want to wait a bit longer if:

- you want a large batteries-included web platform
- you do not want to absorb any migration cost before 1.0

## Versioning

Flamigo plans to follow semantic versioning.

That means the intent is:

- `1.x` and above should use major versions for breaking changes
- minor versions should add functionality in a backwards-compatible way
- patch versions should focus on fixes and smaller improvements

Until Flamigo reaches `1.0.0`, versions below `1.0.x` may still introduce breaking changes. That is intentional for the current stage of the project while the architecture, APIs, and project templates continue to settle through real usage.
