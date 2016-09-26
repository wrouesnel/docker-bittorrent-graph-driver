
GO_SRC := $(shell find -type f -name '*.go' ! -path '*/vendor/*')

SRC_ROOT = github.com/wrouesnel/docker-bittorrent-graph-driver
PROGNAME := docker-bittorrent-graph-driver
VERSION ?= git:$(shell git describe --long --dirty)
TAG ?= latest
CONTAINER_NAME ?= wrouesnel/$(PROGNAME):$(TAG)
BUILD_CONTAINER ?= $(PROGNAME)_build

all: vet test test-style $(PROGNAME)

vet:
	go vet .

# Check code conforms to go fmt
test-style:
	! gofmt -l $(GO_SRC) 2>&1 | read

# Test everything
test:
	go test -v ./...
	
# Format the code
fmt:
	go fmt ./...

# Simple go build
$(PROGNAME): $(GO_SRC)
	GOOS=linux go build -a \
	-ldflags "-extldflags '-static' -X main.Version=$(VERSION)" \
	-o $(PROGNAME) .
	
.PHONY: vet test test-style
