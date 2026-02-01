.PHONY: build test lint clean install

BINARY_NAME ?= codeez
VERSION ?= dev
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/codeez

test:
	go test ./...

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY_NAME)
	go clean -cache

install: build
	go install $(LDFLAGS) ./cmd/codeez
