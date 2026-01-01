# syntax=docker/dockerfile:1.5

# Stage 1 — build
FROM golang:1.23-bullseye AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# Stage 2 — build healthcheck binary
FROM golang:1.23-bullseye AS builder-hc
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download
COPY . .

ARG SKIP_LINT=false

RUN if [ "${SKIP_LINT}" != "true" ]; then \
      curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
        | sh -s -- -b /usr/local/bin v2.6.2; \
    else \
      echo "SKIP_LINT=true — skipping golangci-lint install"; \
    fi

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    if [ "${SKIP_LINT}" != "true" ]; then \
      golangci-lint run --config .golangci.yml ./...; \
    else \
      echo "SKIP_LINT=true — skipping golangci-lint run"; \
    fi

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /healthcheck ./cmd/healthcheck

# Stage 3 — runtime
FROM gcr.io/distroless/static-debian11

COPY --from=builder /app/server /server
COPY --from=builder-hc /healthcheck /healthcheck

EXPOSE 8080
ENTRYPOINT ["/server"]

