include .env

.PHONY: all
all: build run

.PHONY: build
build:
	go build -o file-sentinel ./cmd/file-sentinel

.PHONY: run
run:
	go run ./cmd/file-sentinel

.PHONY: test
test:
	go test ./...
