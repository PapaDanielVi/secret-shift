.PHONY: build test lint cover clean

build:
	go build -ldflags="-s -w -X github.com/PapaDanielVi/secret-shift/cmd.Version=$(shell git describe --tags --always --dirty)" -o bin/secret-shift .

test:
	go test -v -race -count=1 ./...

lint:
	golangci-lint run

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf bin/
