BINARY_NAME=stgui

.PHONY: all build run test clean fmt vet

all: build

build:
	mkdir -p bin
	go build -o bin/$(BINARY_NAME) cmd/main.go

run:
	go run cmd/main.go

test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

clean:
	rm -rf bin

fmt:
	go fmt ./...

vet:
	go vet ./...
