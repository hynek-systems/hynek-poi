# Hynek POI

**Hynek POI** is a high-performance, provider-agnostic Points-of-Interest routing engine.
It aggregates POIs from multiple providers (OpenStreetMap, Google Places, HERE, and more), deduplicates, ranks, and returns the best results via a fast, cache-optimized API.

Hynek POI is designed to be used as:

* A standalone microservice
* A backend for map-based applications
* A routing layer between applications and POI providers
* A self-hosted alternative to proprietary POI APIs

---

# Features

## Core Engine

* Multi-provider aggregation (OSM, Google Places, Foursquare, more coming)
* Parallel provider execution
* Deduplication engine (distance-based)
* Ranking engine (configurable provider priority)
* Category filtering
* Radius search
* Bounding box search
* Paginated results
* Rich POI metadata (ratings, hours, contact info, accessibility, and more)

## Performance

* Redis L2 cache
* In-memory L1 cache
* GeoHash-based cache keys
* 100k+ requests/min capability
* Timeout and retry policies per provider

## Reliability

* Circuit breakers per provider
* Provider-specific timeout configuration
* Automatic retry policies
* Graceful degradation

## Observability

* Prometheus metrics
* Grafana dashboards
* Health and readiness endpoints

## Production-Ready

* Docker support
* Kubernetes-ready
* Config-driven architecture
* Environment variable configuration
* SDK-friendly API design

---

# Architecture

```
Client
  ↓
Hynek POI API
  ↓
Cache Layer (L1 Memory + L2 Redis)
  ↓
Parallel Orchestrator
  ↓
Retry Provider
  ↓
Timeout Provider
  ↓
Circuit Breaker
  ↓
Providers (Google, OSM, Foursquare)
  ↓
Deduplication Engine
  ↓
Ranking Engine
  ↓
Response
```

---

# Quick Start

## Run with Docker Compose

```
docker compose up -d
```

Service will start at:

```
http://localhost:8080
```

---

# API

## Search by Radius

```
GET /v1/search?lat=59.3293&lng=18.0686&radius=1000&categories=restaurant,cafe
```

Example:

```
curl "http://localhost:8080/v1/search?lat=59.3293&lng=18.0686&categories=restaurant"
```

---

## Search by Bounding Box

```
GET /v1/search?bbox=59.32,18.05,59.35,18.10&categories=restaurant
```

Format:

```
bbox=minLat,minLng,maxLat,maxLng
```

---

## Pagination

Results are paginated by default.

Parameters:

* `page` — Page number (default: 1)
* `page_size` — Results per page (default: 20, max: 100)

Example:

```
GET /v1/search?lat=59.3293&lng=18.0686&categories=restaurant&page=2&page_size=10
```

---

## Health Check

```
GET /health
```

---

## Readiness Check

```
GET /ready
```

---

## Metrics

```
GET /metrics
```

Prometheus-compatible.

---

# Example Response

```json
{
  "data": [
    {
      "id": "12345",
      "name": "McDonald's",
      "latitude": 59.3293,
      "longitude": 18.0686,
      "category": "restaurant",
      "source": "google",
      "rating": 4.2,
      "rating_count": 312,
      "address": "Storgatan 1, Stockholm",
      "phone": "+46812345678",
      "website": "https://mcdonalds.se",
      "opening_hours": ["Mon-Sun 06:00-23:00"],
      "open_now": true,
      "price_level": 1,
      "cuisine": "Fast Food",
      "wheelchair_accessible": true,
      "delivery": true,
      "takeaway": true
    }
  ],
  "total": 42,
  "page": 1,
  "page_size": 20,
  "total_pages": 3
}
```

---

# Response Fields

Core fields (always present):

* `id` — Provider-specific place ID
* `name` — Place name
* `latitude` / `longitude` — Coordinates
* `category` — Place category
* `source` — Provider name (google, osm, foursquare)

Enriched fields (included when available):

* `rating` — Rating score
* `rating_count` — Number of ratings
* `address` — Formatted address
* `phone` — Phone number
* `email` — Email address
* `website` — Website URL
* `menu_url` — Menu URL
* `opening_hours` — Opening hours
* `open_now` — Whether the place is currently open
* `cuisine` — Cuisine type or taste tags
* `price_level` — Price level (1-4)
* `description` — Place description
* `wheelchair_accessible` — Wheelchair accessibility
* `outdoor_seating` — Outdoor seating available
* `takeaway` — Takeaway available
* `delivery` — Delivery available
* `verified` — Whether the place is verified
* `popularity` — Popularity score (0-1)

Fields are omitted from the response when not available from the provider.

---

# Configuration

Configuration can be provided via:

* config.yaml
* Environment variables
* Docker / Kubernetes env

---

# Environment Variables

All variables use prefix:

```
HYNEK_POI_
```

Example:

```
HYNEK_POI_SERVER_PORT=8080

HYNEK_POI_REDIS_ADDR=redis:6379

HYNEK_POI_PROVIDERS_GOOGLE_ENABLED=true
HYNEK_POI_PROVIDERS_GOOGLE_API_KEY=xxx

HYNEK_POI_PROVIDERS_OSM_ENABLED=true
```

Full specification: see ENVIRONMENT.md

---

# Example config.yaml

```
server:
  port: 8080

redis:
  addr: redis:6379

cache:
  ttl: 5m

providers:

  google:
    enabled: true
    api_key: xxx
    priority: 1
    timeout: 2s
    retries: 2

  osm:
    enabled: true
    priority: 10
    timeout: 5s
    retries: 1

  foursquare:
    enabled: true
    api_key: xxx
    priority: 5
    timeout: 3s
    retries: 2
```

---

# Supported Providers

| Provider      | Status    |
| ------------- | --------- |
| OpenStreetMap | Supported |
| Google Places | Supported |
| Foursquare    | Supported |
| HERE Maps     | Planned   |

---

# Performance

Typical latency:

| Provider                  | Latency   |
| ------------------------- | --------- |
| Cache hit                 | < 5ms     |
| Google Places             | 50–150ms  |
| OSM                       | 200–800ms |
| Multi-provider aggregated | 50–200ms  |

---

# Metrics

Prometheus metrics include:

```
hynek_poi_requests_total
hynek_poi_cache_hits_total
hynek_poi_cache_misses_total
hynek_poi_request_duration_seconds
```

---

# Docker Deployment

```
docker build -t hynek-poi .
docker run -p 8080:8080 hynek-poi
```

---

# Kubernetes Deployment

Example readiness probe:

```
readinessProbe:
  httpGet:
    path: /ready
    port: 8080
```

---

# SDK Integration

Hynek POI works with any backend:

Laravel:

```
HYNEK_POI_URL=http://localhost:8080
```

Spring Boot:

```
hynek.poi.url=http://localhost:8080
```

Node:

```
process.env.HYNEK_POI_URL
```

---

# Roadmap

Upcoming features:

* HERE Maps provider
* Adaptive provider scoring
* Distributed cache support
* GraphQL endpoint
* Official SDKs

---

# Development

Run locally:

```
go run cmd/api/main.go
```

Run tests:

```
go test ./...
```

---

# Contributing

Contributions are welcome.

Please open an issue or submit a pull request.

---

# License

MIT License

---

# Project Status

Hynek POI is production-ready and actively developed.

---

# Vision

Hynek POI aims to become the open standard routing layer for Points-of-Interest data.

---
