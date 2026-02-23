# Makefile for Hynek POI

APP_NAME=hynek-poi
DOCKER_IMAGE=hyneksystems/hynek-poi

GO=go

# =========================

# Build

# =========================

.PHONY: build
build:
$(GO) build -o $(APP_NAME) ./cmd/api

.PHONY: run
run:
$(GO) run cmd/api/main.go

.PHONY: clean
clean:
rm -f $(APP_NAME)

# =========================

# Test

# =========================

.PHONY: test
test:
$(GO) test ./...

.PHONY: test-race
test-race:
$(GO) test -race ./...

.PHONY: test-cover
test-cover:
$(GO) test -coverprofile=coverage.out ./...
$(GO) tool cover -html=coverage.out

# =========================

# Lint

# =========================

.PHONY: lint
lint:
golangci-lint run

.PHONY: fmt
fmt:
$(GO) fmt ./...

.PHONY: vet
vet:
$(GO) vet ./...

# =========================

# Dependencies

# =========================

.PHONY: deps
deps:
$(GO) mod download

.PHONY: tidy
tidy:
$(GO) mod tidy

# =========================

# Docker

# =========================

.PHONY: docker-build
docker-build:
docker build -t $(DOCKER_IMAGE):local .

.PHONY: docker-run
docker-run:
docker run -p 8080:8080 $(DOCKER_IMAGE):local

.PHONY: docker-compose-up
docker-compose-up:
docker compose up -d

.PHONY: docker-compose-down
docker-compose-down:
docker compose down

.PHONY: docker-compose-logs
docker-compose-logs:
docker compose logs -f

# =========================

# Redis

# =========================

.PHONY: redis-flush
redis-flush:
docker compose exec redis redis-cli flushall

# =========================

# Dev

# =========================

.PHONY: dev
dev: fmt vet lint test

# =========================

# CI parity

# =========================

.PHONY: ci
ci: deps fmt vet lint test-race build

# =========================

# Release

# =========================

.PHONY: release
release:
ifndef VERSION
$(error VERSION is required, usage: make release VERSION=v0.1.0)
endif
git tag $(VERSION)
git push origin $(VERSION)

# =========================

# Help

# =========================

.PHONY: help
help:
@echo ""
@echo "Hynek POI Makefile"
@echo ""
@echo "Build:"
@echo "  make build"
@echo "  make run"
@echo ""
@echo "Testing:"
@echo "  make test"
@echo "  make test-race"
@echo "  make test-cover"
@echo ""
@echo "Lint:"
@echo "  make lint"
@echo "  make fmt"
@echo ""
@echo "Docker:"
@echo "  make docker-build"
@echo "  make docker-compose-up"
@echo ""
@echo "CI:"
@echo "  make ci"
@echo ""
