# Makefile for local tests, lint
# (release is now goreleaser)

PACKAGES ?= $(shell go list ./...)

# Local targets:
go-install:
	go install $(PACKAGES)

TEST_TIMEOUT:=90s

OS:=$(shell go env GOOS)

# Local test
test:
ifeq ($(OS),windows)
	@echo "Skipping tests on Windows until we can figure out what's wrong with testscript on windows."
else
	go test -timeout $(TEST_TIMEOUT) -race $(PACKAGES)
endif

# To debug strange linter errors, uncomment
# DEBUG_LINTERS="--debug"

.golangci.yml: Makefile
	curl -fsS -o .golangci.yml https://raw.githubusercontent.com/fortio/workflows/main/golangci.yml

lint: .golangci.yml
	golangci-lint $(DEBUG_LINTERS) run $(LINT_PACKAGES)

coverage:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...

.PHONY: lint coverage test
