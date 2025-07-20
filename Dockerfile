# Multi-stage Dockerfile for local TelDrive builds
# Replaces goreleaser.dockerfile for local development and deployment

# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build UI assets from source (simplified approach)
ARG UI_SOURCE_TYPE="build"  # Primary option: "build" from source
ARG UI_REPO_URL="https://github.com/cdgmx/teldrive-ui.git"
ARG UI_REPO_BRANCH="main"
ARG GITHUB_TOKEN=""
ENV GITHUB_TOKEN=$GITHUB_TOKEN

# Install Node.js and pnpm for UI building
RUN echo "Installing Node.js and pnpm for UI building..."; \
    apk add --no-cache nodejs npm; \
    npm install -g pnpm

# Build UI from source repository
RUN echo "Building UI from source repository: $UI_REPO_URL"; \
    cd /tmp && \
    git clone --depth 1 --branch "$UI_REPO_BRANCH" "$UI_REPO_URL" teldrive-ui-src && \
    cd teldrive-ui-src && \
    echo "Installing UI dependencies..." && \
    pnpm install && \
    echo "Building UI..." && \
    pnpm run build && \
    echo "Copying built UI assets..." && \
    mkdir -p /app/ui/dist && \
    cp -r dist/* /app/ui/dist/ && \
    cd /app && \
    rm -rf /tmp/teldrive-ui-src && \
    echo "✅ UI built successfully from source"

# Generate API code (matching taskfile.yml) - Skip for faster builds
# RUN go generate ./...

# Build the application
ARG BUILD_MODE="production"
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags "-s -w -extldflags=-static" \
    -o teldrive \
    main.go

# Runtime stage
FROM scratch

# Copy CA certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /app/teldrive /teldrive

# Expose port
EXPOSE 8080

# Set entrypoint matching current goreleaser.dockerfile
ENTRYPOINT ["/teldrive","run","--tg-storage-file","/storage.db"]