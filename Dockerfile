# Dashboard build stage
FROM oven/bun:latest AS dashboard
WORKDIR /build
COPY dashboard/package.json dashboard/bun.lock ./
RUN bun install --frozen-lockfile
COPY dashboard/ .
RUN bun run build

# Go build stage
FROM golang:1.26-alpine AS builder
RUN apk add --no-cache git
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=dashboard /build/dist/ internal/dashboard/dashboard/
RUN mv internal/dashboard/dashboard/index.html internal/dashboard/dashboard/dashboard.html
RUN go build -ldflags="-s -w" -o pinchtab ./cmd/pinchtab

# Runtime stage
FROM alpine:latest

LABEL org.opencontainers.image.source="https://github.com/pinchtab/pinchtab"
LABEL org.opencontainers.image.description="High-performance browser automation bridge"

# Install Chromium and dependencies
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    dumb-init

# Create non-root user and persistent config/state directory
RUN adduser -D -h /data -g '' pinchtab && \
    mkdir -p /data && \
    chown pinchtab:pinchtab /data

# Copy binary from builder
COPY --from=builder /build/pinchtab /usr/local/bin/pinchtab
COPY --chmod=0755 docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

# Switch to non-root user
USER pinchtab
WORKDIR /data

# Environment variables
ENV HOME=/data \
    XDG_CONFIG_HOME=/data/.config

# Expose port
EXPOSE 9867

# Use dumb-init to properly handle signals
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

# Run pinchtab
CMD ["/usr/local/bin/docker-entrypoint.sh", "pinchtab"]
