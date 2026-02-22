# Hynek POI Architecture

Hynek POI is designed as a high-performance, provider-agnostic POI routing engine.

This document describes its internal architecture and execution flow.

---

# Overview

Hynek POI acts as a routing layer between clients and POI providers.

Responsibilities:

* Query multiple providers
* Execute providers in parallel
* Apply retry and timeout policies
* Apply circuit breakers
* Deduplicate results
* Rank results
* Cache results
* Return optimized response

---

# High-Level Architecture

```
Client
  ↓
HTTP API Layer
  ↓
Cached Orchestrator
  ↓
Parallel Orchestrator
  ↓
Provider Execution Layer
  ↓
Deduplication Engine
  ↓
Ranking Engine
  ↓
Response
```

---

# Component Breakdown

## HTTP Layer

Location:

```
cmd/api/
```

Responsibilities:

* Parse HTTP requests
* Validate input
* Convert to SearchQuery
* Return JSON response

Endpoints:

```
/v1/search
/health
/ready
/metrics
```

---

## Orchestrator Layer

Location:

```
internal/orchestrator/
```

Core component responsible for coordinating providers.

Execution flow:

```
Parallel execution
Timeout enforcement
Aggregation
Deduplication
Ranking
```

Key implementations:

```
ParallelOrchestrator
CachedOrchestrator
```

---

## Cache Layer

Location:

```
internal/cache/
```

Two-level cache:

L1 Cache:

```
In-memory cache
Fast access (<1ms)
```

L2 Cache:

```
Redis cache
Shared across instances
```

Cache key includes:

```
GeoHash
Categories
Radius or BBox
```

---

## Provider Layer

Location:

```
internal/provider/
```

Provider interface:

```
type Provider interface {
    Name() string
    Search(query SearchQuery) ([]POI, error)
}
```

Providers:

```
Google Places
OpenStreetMap
HERE (planned)
```

---

## Provider Execution Stack

Each provider is wrapped in multiple resilience layers:

```
RetryProvider
  ↓
TimeoutProvider
  ↓
CircuitBreakerProvider
  ↓
BaseProvider
```

Responsibilities:

RetryProvider:

```
Retry failed requests
```

TimeoutProvider:

```
Abort slow providers
```

CircuitBreakerProvider:

```
Prevent cascading failures
```

---

## Deduplication Engine

Location:

```
internal/dedupe/
```

Removes duplicate POIs using:

```
Name normalization
Distance threshold
Provider merging
```

---

## Ranking Engine

Location:

```
internal/ranking/
```

Sorts POIs using:

```
Provider priority
Distance from query
Future: ratings, popularity
```

---

## Config System

Location:

```
internal/config/
```

Supports:

```
config.yaml
Environment variables
Docker/Kubernetes config
```

Uses Viper for configuration loading.

---

## Metrics System

Location:

```
internal/metrics/
```

Exposes Prometheus metrics:

```
hynek_poi_requests_total
hynek_poi_cache_hits_total
hynek_poi_cache_misses_total
hynek_poi_request_duration_seconds
```

---

# Execution Flow Example

```
Client sends request
  ↓
HTTP handler creates SearchQuery
  ↓
Cache lookup (L1 → L2)
  ↓
If miss:
  ↓
ParallelOrchestrator executes providers
  ↓
RetryProvider retries failures
  ↓
TimeoutProvider enforces timeout
  ↓
CircuitBreaker prevents failing providers
  ↓
Results collected
  ↓
Deduplication removes duplicates
  ↓
Ranking sorts results
  ↓
Results cached
  ↓
Response returned
```

---

# Concurrency Model

Hynek POI uses Go concurrency primitives:

```
goroutines
channels
context cancellation
```

Providers execute fully in parallel.

---

# Scaling Model

Hynek POI is stateless.

Horizontal scaling supported:

```
Multiple instances
Shared Redis cache
Load balancer
Kubernetes deployment
```

---

# Failure Handling

Failure isolation per provider:

```
Timeout isolation
Retry policies
Circuit breaker protection
```

System continues operating even if providers fail.

---

# Performance Characteristics

Typical latency:

```
Cache hit: < 5ms
Google provider: 50–150ms
OSM provider: 200–800ms
Parallel multi-provider: 50–200ms
```

Throughput:

```
100k+ requests/min per instance
```

---

# Deployment Architecture

Typical production deployment:

```
Load Balancer
    ↓
Hynek POI instances
    ↓
Redis cluster
```

---

# Future Architecture Extensions

Planned features:

```
Additional providers
Adaptive provider scoring
Distributed cache
GraphQL endpoint
SDK integrations
```

---

# Design Goals

Hynek POI is designed for:

```
Performance
Reliability
Provider independence
Scalability
Observability
Extensibility
```

---

# Summary

Hynek POI is a resilient, scalable routing layer for POI aggregation, designed to operate as critical infrastructure in modern applications.

---
