# Hynek POI Environment Variable Specification

This document defines all supported environment variables for Hynek POI.

All variables use the required prefix:

```
HYNEK_POI_
```

This prevents conflicts with host applications and ensures safe deployment in shared environments such as Docker, Kubernetes, Laravel, Spring Boot, and Node.js applications.

---

# Configuration Precedence

Configuration is loaded in the following order (highest priority first):

1. Environment variables
2. config.yaml
3. Built-in defaults

---

# Server Configuration

## HYNEK_POI_SERVER_PORT

Port used by the HTTP server.

Default:

```
8080
```

Example:

```
HYNEK_POI_SERVER_PORT=9090
```

---

## HYNEK_POI_SERVER_READ_TIMEOUT

HTTP read timeout.

Default:

```
5s
```

---

## HYNEK_POI_SERVER_WRITE_TIMEOUT

HTTP write timeout.

Default:

```
5s
```

---

# Redis Configuration

## HYNEK_POI_REDIS_ADDR

Redis server address.

Default:

```
localhost:6379
```

Docker example:

```
HYNEK_POI_REDIS_ADDR=redis:6379
```

---

## HYNEK_POI_REDIS_PASSWORD

Redis password.

Default:

```
(empty)
```

---

## HYNEK_POI_REDIS_DB

Redis database index.

Default:

```
0
```

---

# Cache Configuration

## HYNEK_POI_CACHE_TTL

Cache TTL duration.

Default:

```
5m
```

Examples:

```
HYNEK_POI_CACHE_TTL=1m
HYNEK_POI_CACHE_TTL=10m
```

---

## HYNEK_POI_CACHE_L1_SIZE

Maximum in-memory cache entries.

Default:

```
10000
```

---

# Provider Configuration

Format:

```
HYNEK_POI_PROVIDERS_<PROVIDER>_<SETTING>
```

Supported providers:

```
GOOGLE
OSM
HERE (future)
FOURSQUARE (future)
```

---

# Google Provider

## HYNEK_POI_PROVIDERS_GOOGLE_ENABLED

Enable Google Places provider.

Default:

```
false
```

---

## HYNEK_POI_PROVIDERS_GOOGLE_API_KEY

Google Places API key.

Required if provider enabled.

Example:

```
HYNEK_POI_PROVIDERS_GOOGLE_API_KEY=your_api_key
```

---

## HYNEK_POI_PROVIDERS_GOOGLE_PRIORITY

Provider priority.

Lower value = higher priority.

Default:

```
1
```

---

## HYNEK_POI_PROVIDERS_GOOGLE_TIMEOUT

Request timeout.

Default:

```
2s
```

---

## HYNEK_POI_PROVIDERS_GOOGLE_RETRIES

Retry attempts.

Default:

```
2
```

---

## HYNEK_POI_PROVIDERS_GOOGLE_CB_FAILURES

Circuit breaker failure threshold.

Default:

```
3
```

---

## HYNEK_POI_PROVIDERS_GOOGLE_CB_RESET_TIMEOUT

Circuit breaker reset timeout.

Default:

```
30s
```

---

# OpenStreetMap Provider

## HYNEK_POI_PROVIDERS_OSM_ENABLED

Enable OSM provider.

Default:

```
true
```

---

## HYNEK_POI_PROVIDERS_OSM_PRIORITY

Default:

```
10
```

---

## HYNEK_POI_PROVIDERS_OSM_TIMEOUT

Default:

```
5s
```

---

## HYNEK_POI_PROVIDERS_OSM_RETRIES

Default:

```
1
```

---

## HYNEK_POI_PROVIDERS_OSM_CB_FAILURES

Default:

```
3
```

---

## HYNEK_POI_PROVIDERS_OSM_CB_RESET_TIMEOUT

Default:

```
30s
```

---

# Router Configuration

## HYNEK_POI_ROUTER_TIMEOUT

Maximum total routing time.

Default:

```
3s
```

---

## HYNEK_POI_ROUTER_MAX_PROVIDERS

Maximum parallel providers.

Default:

```
5
```

---

# Metrics

## HYNEK_POI_METRICS_ENABLED

Enable Prometheus metrics.

Default:

```
true
```

---

## HYNEK_POI_METRICS_PATH

Metrics endpoint path.

Default:

```
/metrics
```

---

# Logging

## HYNEK_POI_LOG_LEVEL

Options:

```
debug
info
warn
error
```

Default:

```
info
```

---

## HYNEK_POI_LOG_FORMAT

Options:

```
json
text
```

Default:

```
json
```

---

# Config File Override

## HYNEK_POI_CONFIG_FILE

Override config file location.

Example:

```
HYNEK_POI_CONFIG_FILE=/etc/hynek-poi/config.yaml
```

---

# Full Example

```
HYNEK_POI_SERVER_PORT=8080

HYNEK_POI_REDIS_ADDR=redis:6379

HYNEK_POI_CACHE_TTL=5m

HYNEK_POI_PROVIDERS_GOOGLE_ENABLED=true
HYNEK_POI_PROVIDERS_GOOGLE_API_KEY=xxx
HYNEK_POI_PROVIDERS_GOOGLE_PRIORITY=1
HYNEK_POI_PROVIDERS_GOOGLE_TIMEOUT=2s
HYNEK_POI_PROVIDERS_GOOGLE_RETRIES=2

HYNEK_POI_PROVIDERS_OSM_ENABLED=true
HYNEK_POI_PROVIDERS_OSM_PRIORITY=10
HYNEK_POI_PROVIDERS_OSM_TIMEOUT=5s
HYNEK_POI_PROVIDERS_OSM_RETRIES=1
```

---

# Stability Guarantee

All environment variables in this document are considered part of the public API and follow semantic versioning.

---
