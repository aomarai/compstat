FROM golang:1.21-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o compstat ./cmd/compstat

# Final stage
FROM alpine:latest
LABEL author="Ash Omaraie"

RUN apk add --no-cache \
    zstd \
    xz \
    pigz \
    lz4 \
    bzip2 \
    pbzip2 \
    brotli \
    bash \
    ca-certificates

# Create non-root user
RUN addgroup -g 1000 compstat && \
    adduser -D -u 1000 -G compstat compstat

# Copy binary from builder
COPY --from=builder /build/compstat /usr/local/bin/compstat

WORKDIR /data
RUN chown -R compstat:compstat /data

USER compstat

ENTRYPOINT ["compstat"]
CMD ["--help"]