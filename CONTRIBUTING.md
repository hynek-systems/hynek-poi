# Contributing to Hynek POI

Thank you for your interest in contributing to Hynek POI.

Hynek POI aims to be a high-performance, open, provider-agnostic POI routing engine. Contributions of all kinds are welcome.

---

# Ways to Contribute

You can contribute by:

* Reporting bugs
* Suggesting features
* Improving documentation
* Adding new providers (HERE, Foursquare, etc.)
* Improving performance
* Writing tests
* Improving observability
* Fixing bugs

---

# Development Setup

## Requirements

* Go 1.25+
* Docker (recommended)
* Redis (optional if using Docker)

---

## Clone the repository

```id="z8qk4m"
git clone https://github.com/hynek-systems/hynek-poi.git
cd hynek-poi
```

---

## Run locally

```id="dhtz3x"
go run cmd/api/main.go
```

---

## Run with Docker

```id="8y9ywx"
docker compose up -d
```

---

## Run tests

```id="5g7t0l"
go test ./...
```

---

## Run with race detection

```id="cz3s37"
go test -race ./...
```

---

# Project Structure

```id="g3swnq"
cmd/api/                 HTTP entrypoint
internal/cache/          Cache layer
internal/config/         Config system
internal/provider/       Provider implementations
internal/orchestrator/   Routing engine
internal/dedupe/         Deduplication engine
internal/ranking/        Ranking engine
internal/metrics/        Prometheus metrics
```

---

# Coding Guidelines

## Go Standards

Follow standard Go conventions:

* gofmt
* go vet
* idiomatic Go

Format code before committing:

```id="v0h3yr"
go fmt ./...
```

---

## Naming

Use clear, descriptive names.

Avoid abbreviations unless standard.

---

## Error Handling

Always return errors.

Do not ignore errors silently.

Example:

```id="4k14x2"
if err != nil {
    return nil, err
}
```

---

## Concurrency

Use:

* context for cancellation
* channels carefully
* avoid race conditions

---

# Adding a New Provider

To add a new provider:

1. Create new provider file:

```id="l7ejb7"
internal/provider/<provider>_provider.go
```

2. Implement Provider interface:

```id="z7fndh"
type Provider interface {
    Name() string
    Search(SearchQuery) ([]POI, error)
}
```

3. Register provider in:

```id="p5g8rt"
internal/provider/registry.go
```

4. Add config support

5. Add tests

---

# Pull Request Process

1. Fork the repository
2. Create feature branch:

```id="x7ff4e"
git checkout -b feature/my-feature
```

3. Commit changes:

```id="7fjys0"
git commit -m "Add HERE provider"
```

4. Push branch:

```id="33b0ym"
git push origin feature/my-feature
```

5. Open Pull Request

---

# Pull Request Requirements

PRs must:

* Compile successfully
* Pass tests
* Follow coding standards
* Include tests when applicable
* Not break existing functionality

---

# Reporting Bugs

Open an issue with:

* Description
* Steps to reproduce
* Expected behavior
* Actual behavior
* Logs if available

---

# Feature Requests

Open an issue describing:

* Problem
* Proposed solution
* Use case

---

# Code of Conduct

Be respectful.

No harassment or abusive behavior.

---

# License

By contributing, you agree your contributions will be licensed under the MIT License.

---

Thank you for helping make Hynek POI better.
