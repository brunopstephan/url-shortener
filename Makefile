BINARY=build/shortener
SRC_DIR=./cmd/api

.PHONY: all build clean run test

all: clean test build

build:
		go build -o $(BINARY) $(SRC_DIR)

run:
		go run $(SRC_DIR)

test:
		go test ./...

clean:
		rm -f $(BINARY)
