# Stage 1 — build
FROM golang:1.23-bullseye AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=ssh go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# Stage 2 — build healthcheck binary
FROM golang:1.23-bullseye AS builder-hc
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=ssh  go mod download
COPY . .

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
    | sh -s -- -b /usr/local/bin v1.60.1
RUN golangci-lint run ./...


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /healthcheck ./cmd/healthcheck

# Stage 3 — runtime
FROM gcr.io/distroless/static-debian11

COPY --from=builder /app/server /server
COPY --from=builder-hc /healthcheck /healthcheck

EXPOSE 8080
ENTRYPOINT ["/server"]

