EXAMPLE?=counter

.PHONY: all build examples run test coverage clean fmt vet check

all: check

build:
	go build ./...

examples: build
	mkdir -p bin
	go build -o bin/ ./examples/...

run:
	go run ./examples/$(EXAMPLE)

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

clean:
	rm -rf bin coverage.out

fmt:
	go fmt ./...

vet:
	go vet ./...

check: fmt vet test
