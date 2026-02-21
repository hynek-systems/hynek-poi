# ---- Build stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# install git (needed for go modules sometimes)
RUN apk add --no-cache git

# copy go mod files first (better caching)
COPY go.mod go.sum ./

RUN go mod download

# copy source
COPY . .

# build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o hynek-poi cmd/api/main.go


# ---- Runtime stage ----
FROM alpine:3.19

WORKDIR /app

# copy binary
COPY --from=builder /app/hynek-poi .

# copy config
COPY config.yaml .

# expose port
EXPOSE 8080

# run binary
CMD ["./hynek-poi"]