# Hynek POI Release Process

This document defines the release process for Hynek POI.

Hynek POI follows Semantic Versioning:

```id="5fw9z5"
MAJOR.MINOR.PATCH
```

Example:

```id="jzgjlp"
v0.1.0
v1.0.0
v1.1.0
v1.1.1
```

---

# Version Definitions

## MAJOR

Breaking changes.

Example:

* API changes
* Config changes
* Architecture changes

---

## MINOR

New features without breaking changes.

Example:

* New provider
* New endpoint
* New config options

---

## PATCH

Bug fixes only.

Example:

* Bug fixes
* Performance improvements
* Documentation fixes

---

# Release Checklist

Before releasing:

## Code

* Tests pass

```id="nnt6s6"
go test ./...
```

* No race conditions

```id="hm7p8f"
go test -race ./...
```

* Code formatted

```id="hyb7v4"
go fmt ./...
```

---

## Build

Verify build works:

```id="i3dp7g"
go build ./cmd/api
```

---

## Docker

Verify Docker build:

```id="l7qypz"
docker build -t hynek-poi .
```

---

## Documentation

Ensure updated:

* README.md
* ENVIRONMENT.md
* ARCHITECTURE.md

---

# Release Steps

## 1. Update version

Create tag:

```id="k1y5r9"
git tag v0.1.0
```

---

## 2. Push tag

```id="f9f12j"
git push origin v0.1.0
```

---

## 3. Build Docker image

```id="g4e5c2"
docker build -t hynek-poi:v0.1.0 .
```

---

## 4. Publish Docker image

Example:

```id="okd2w4"
docker tag hynek-poi:v0.1.0 hyneksystems/hynek-poi:v0.1.0
docker push hyneksystems/hynek-poi:v0.1.0
```

---

## 5. Create GitHub Release

Include:

* Release notes
* Changes
* Breaking changes if any

---

# Release Notes Template

```id="3y1w0j"
# Hynek POI v0.1.0

Initial release.

Features:

- Multi-provider routing
- Google Places provider
- OpenStreetMap provider
- Parallel orchestrator
- Deduplication engine
- Ranking engine
- Redis cache
- Prometheus metrics
```

---

# Versioning Policy

Hynek POI guarantees:

* No breaking changes in PATCH releases
* No breaking changes in MINOR releases
* Breaking changes only in MAJOR releases

---

# Production Deployment Recommendation

Use tagged versions:

```id="scdqlg"
hyneksystems/hynek-poi:v0.1.0
```

Avoid:

```id="n3r4gh"
latest
```

---

# Long-Term Support

Future LTS versions will be designated.

---

# Summary

Release process ensures:

* Stability
* Compatibility
* Reliability
* Reproducibility

---
