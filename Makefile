.PHONY: build lint

build:
	go build ./cmd/idl

lint:
	golangci-lint run
