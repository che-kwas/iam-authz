.DEFAULT_GOAL := help

# ==============================================================================
# Build options

VERSION      := $(shell git describe --tags --always)
OUTPUT_DIR   := ./_output
GO_LDFLAGS   += -X main.Version=$(VERSION)
MAKEFLAGS    += --no-print-directory

# ==============================================================================
# Includes

include make-rules/tools.mk

# ==============================================================================
# Targets

## all: Build all.
.PHONY: all
all: lint test build

## lint: Check syntax and styling of go sources.
.PHONY: lint
lint: tools.verify.golangci-lint
	go mod tidy -compat=1.17
	golangci-lint run ./...

## test: Run unit test.
.PHONY: test
test:
	@-mkdir -p $(OUTPUT_DIR)
	go test -race -cover -coverprofile=$(OUTPUT_DIR)/coverage.out ./...

## cover: Run unit test and get test coverage.
.PHONY: cover
cover: test
	sed -i '/mock_.*.go/d' $(OUTPUT_DIR)/coverage.out
	go tool cover -html=$(OUTPUT_DIR)/coverage.out -o $(OUTPUT_DIR)/coverage.html

## build: Build source code for host platform.
.PHONY: build
build:
	go build -ldflags "$(GO_LDFLAGS)" -o $(OUTPUT_DIR)/ ./...

## update: Update all modules.
.PHONY: update
update:
	go get -u ./...
	go mod tidy -compat=1.17

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
	-rm -vrf $(OUTPUT_DIR)

## help: Show help info.
.PHONY: help
help: Makefile
	@echo "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
