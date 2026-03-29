# What is Flamigo?

::: warning
Flamigo is still pre-1.0. The architecture is already usable, but APIs, templates, and tooling may still change.
:::

Flamigo is a backend framework for Go applications that want a domain-first, hexagonal architecture without turning the whole codebase into framework glue.

It gives you a small set of architectural building blocks:

- `domains` for business logic and contracts
- `events` for domain event flow
- `strategies` for application actions
- `adapters` for HTTP, WebSocket, and external integrations
- `injection` for startup wiring

## What Flamigo is for

Flamigo is a good fit when you want:

- clear domain boundaries
- transport-independent application logic
- event-driven coordination between modules
- reusable HTTP and WebSocket transport primitives
- a generated project structure that encourages consistency

Typical examples:

- game backends
- modular monoliths
- service backends with multiple integration points
- systems where the same domain events drive internal reactions and client updates

## What Flamigo is not for

Flamigo is not trying to be:

- a full web platform
- an ORM
- a database abstraction layer
- a framework that decides your auth, persistence, or validation policy

Those decisions are expected to stay in the app unless Flamigo provides an optional helper for them.

## Architectural direction

Flamigo follows a hexagonal direction:

- domains expose contracts and business concepts
- adapters connect the app to the outside world
- strategies represent application actions
- events connect internal workflows

The framework tries to keep the core small and push policy into the app when that policy is not broadly reusable.
