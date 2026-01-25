.PHONY: build run test clean

BINARY_NAME=envm

build:
	go build -o bin/$(BINARY_NAME) ./cmd/envm

run: build
	./bin/$(BINARY_NAME)

test:
	go test -v ./...

clean:
	go clean
	rm -f bin/$(BINARY_NAME)
