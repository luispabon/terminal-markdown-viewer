.PHONY: build test check

build:
	mkdir -p bin
	go build -o bin/markdown-viewer ./cmd/markdown-viewer

test:
	go test ./...

check:
	test -z "$$(gofmt -l cmd/markdown-viewer internal/viewer)"
	go vet ./...
	go build ./...
	go test ./...
