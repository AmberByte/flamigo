---
# https://vitepress.dev/reference/default-theme-home-page
layout: home

hero:
  name: "Flamigo"
  text: "Hexagonal, Event-Driven Go"
  tagline: A backend framework for domain-first Go applications
  image:
    src: ./logo.png
    alt: Logo
  actions:
    - theme: brand
      text: Getting Started
      link: /guide/introduction.html
    - theme: alt
      text: GitHub
      link: https://github.com/amberbyte/flamigo
features:
  - title: Domain-First Structure
    details: Organize applications around domains, strategies, events, and adapters instead of transport-first folders.
  - title: Event Bus
    details: Coordinate domains through explicit domain events without coupling everything directly together.
  - title: Transport Adapters
    details: Reuse framework transport primitives for HTTP and WebSocket adapters while keeping app behavior local.

---

## Why Flamigo?

Flamigo is designed for Go backends that want architectural structure without turning into a giant framework.

It works especially well for:

- game backends
- modular monoliths
- internal backends with multiple adapters
- services where domain boundaries and event flow matter

It is not trying to be an ORM, a frontend stack, or a highly opinionated full web platform.
