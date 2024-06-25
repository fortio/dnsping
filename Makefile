# Makefile for local tests, lint
# (release is now goreleaser)

PACKAGES ?= $(shell go list ./...)

# Local targets:
go-install:
	go install $(PACKAGES)

TEST_TIMEOUT:=90s

# Local test
test:
	go test -timeout $(TEST_TIMEOUT) -race $(PACKAGES)

# To debug strange linter errors, uncomment
# DEBUG_LINTERS="--debug"

.golangci.yml: Makefile
	curl -fsS -o .golangci.yml https://raw.githubusercontent.com/fortio/workflows/main/golangci.yml

lint: .golangci.yml
	golangci-lint $(DEBUG_LINTERS) run $(LINT_PACKAGES)

coverage:
	./.circleci/coverage.sh
	curl -s https://codecov.io/bash | bash


.PHONY: lint coverage test

