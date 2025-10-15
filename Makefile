.PHONY: build test clean install benchmark

# Build variables
BINARY_NAME=compstat
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Build the binary
build:
	go build ${LDFLAGS} -o ${BINARY_NAME} ./cmd/compstat

# Build for all platforms
build-all:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-amd64 ./cmd/compstat
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-arm64 ./cmd/compstat
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-amd64 ./cmd/compstat
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-arm64 ./cmd/compstat
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-windows-amd64.exe ./cmd/compstat

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...

# Run benchmarks
benchmark:
	./scripts/run_benchmark.sh

# Install locally
install:
	go install ${LDFLAGS} ./cmd/compstat

# Clean build artifacts
clean:
	rm -f ${BINARY_NAME}
	rm -rf bin/
	rm -f coverage.out

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run