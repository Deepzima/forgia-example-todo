# Multi-stage Dockerfile for todo-operator
# Stage 1: Build the Go binary
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src

# Copy dependency files first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY kinds/ kinds/

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /out/operator ./cmd/operator/

# Stage 2: Minimal runtime image
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /out/operator /operator

# Run as non-root user (distroless nonroot uid=65532)
USER 65532:65532

EXPOSE 9090

ENTRYPOINT ["/operator"]
