# Contributing to Flamigo

Thanks for contributing to Flamigo.

Flamigo is an open source Go backend framework focused on hexagonal architecture, domain-first design, and reusable transport adapters. Contributions are welcome across code, docs, templates, bug reports, and examples.

## Before You Start

Flamigo is still pre-`1.0`.

That means contributions should favor:

- clear architecture over convenience-driven complexity
- small, composable framework primitives
- app-owned policy where possible
- consistency between framework code, docs, and generated templates

## Good Ways to Contribute

Useful contributions include:

- bug reports with clear reproduction steps
- documentation improvements
- template fixes and generated project cleanup
- tests for missing edge cases
- API ergonomics improvements
- transport adapter improvements
- examples from real app usage

## Discussing Changes

For small fixes, typo fixes, missing tests, and focused bug fixes, a pull request is usually fine directly.

For larger changes, architecture changes, new public APIs, or new optional features, please open an issue first so the direction can be discussed before implementation starts.

Examples of changes that should be discussed first:

- new top-level packages
- changes to generated project structure
- new framework-level abstractions
- changes that introduce or remove dependencies
- changes that affect public APIs or expected architecture

## Development Setup

Requirements:

- Go installed locally

Run the test suite:

```bash
go test ./...
```

If you change generated templates, documentation, or public APIs, make sure the related parts stay aligned.

## Project Expectations

When contributing, please keep these expectations in mind:

- keep the framework small and focused
- avoid moving app-specific policy into framework core
- prefer explicit boundaries between domains, strategies, events, and adapters
- keep transport concerns inside transport adapters
- keep generated templates consistent with the framework and the docs

In practice, that often means:

- `events` should stay transport-agnostic
- `strategies` should express application actions, not transport details
- adapters should connect the outside world to the application core
- convenience helpers should not force unnecessary dependencies into core packages

## Pull Requests

A good pull request usually:

- has a clear and narrow scope
- includes tests when behavior changes
- updates docs when public behavior changes
- updates templates when generated app behavior changes
- explains any architectural tradeoffs

If your change introduces breaking behavior, call that out clearly in the pull request description.

## Documentation and Templates

For Flamigo, docs and templates are part of the product.

If you change:

- public package APIs
- generated project structure
- transport scaffolding
- framework terminology

then please update the relevant docs and templates in the same change when possible.

## Release Expectations

Flamigo intends to follow semantic versioning.

Before `1.0.0`, some breaking changes may still happen between releases. Contributions should therefore aim to improve clarity and long-term direction, not just add more surface area.

## Questions

If you are unsure whether an idea fits the project, open an issue and describe:

- the problem you are trying to solve
- why the change belongs in the framework instead of app code
- what tradeoffs the proposal introduces

That makes it much easier to review the idea well.
