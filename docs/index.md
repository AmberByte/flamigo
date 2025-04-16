---
# https://vitepress.dev/reference/default-theme-home-page
layout: home

hero:
  name: "Flamigo (Preview)"
  text: "Domain Driven Go"
  tagline: Scalable Backend Framework for building Domain Driven Go Backends
  image:
    src: ./logo.png
    alt: Logo
  actions:
    - theme: brand
      text: Getting Started
      link: /guide/introduction.html
features:
  - title: Dependency Injection
    details: Flamigo uses a fully decoupled architecture, with all components provided via dependency injection.
  - title: Event-Driven
    details: Domains communicate seamlessly through a robust event-driven system.
  - title: Real-Time Support
    details: Built-in support for real-time communication with the frontend.

---

## Why Flamigo?

Flamigo is designed to simplify backend development by providing:

- **Domain-Driven Design**: Encourages clean architectural patterns for scalable and maintainable codebases.
- **Event-Driven Communication**: Built-in support for event-driven systems to decouple domains.
- **Real-Time Features**: Seamless integration with WebSockets for real-time communication.
- **Dependency Injection**: Fully decoupled architecture with dependency injection for better testability and modularity.
- **Go Ecosystem**: Built on proven Go libraries like `gorilla`, ensuring reliability and performance.

Whether you're building microservices or monoliths, Flamigo helps you focus on your business logic without getting bogged down by boilerplate.

