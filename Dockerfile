# Stage 1 — build
FROM golang:1.23-bullseye AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=ssh go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# Stage 2 — minimal runtime
FROM gcr.io/distroless/static-debian11

COPY --from=builder /app/server /server

EXPOSE 8080

ENTRYPOINT ["/server"]