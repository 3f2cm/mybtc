MAKEFILE_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

.PHONY: all
all: build test

.PHONY: build
build: mybtc uxto-summary

mybtc:
	CGO_ENABLED=0 go build -ldflags '-s -w' -o mybtc $(MAKEFILE_DIR)main.go

uxto-summary:
	CGO_ENABLED=0 go build -ldflags '-s -w' -o uxto-summary $(MAKEFILE_DIR)blockstream/cmd/uxto-summary/main.go

.PHONY: test
test:
	go test --count=1 -v ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: clean
clean:
	rm -f $(MAKEFILE_DIR){mybtc,uxto-summary}
