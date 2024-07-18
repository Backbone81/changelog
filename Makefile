.PHONY: all
all: build

.PHONY: prepare
prepare:
	go mod tidy
	go fmt ./...
	go vet ./...

.PHONY: build
build: prepare
	go build -o changelog.bin ./cmd/changelog

.PHONY: test
test: prepare
	go test ./...
