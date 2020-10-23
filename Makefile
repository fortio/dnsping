# Makefile to build dnsping's docker images as well as short cut
# for local test/install
#
# See also release/release.sh and fortio's release/Readme.md
#

DOCKER_PREFIX := docker.io/fortio/dnsping
BUILD_IMAGE_TAG := v27
BUILD_IMAGE :=  docker.io/fortio/fortio.build:$(BUILD_IMAGE_TAG)

TAG:=$(USER)$(shell date +%y%m%d_%H%M%S)

DOCKER_TAG = $(DOCKER_PREFIX)$(IMAGE):$(TAG)

PACKAGES ?= $(shell go list ./...)

# Local targets:
go-install:
	go install $(PACKAGES)

# Run/test dependencies
dependencies:

TEST_TIMEOUT:=90s

# Local test
test: dependencies
	go test -timeout $(TEST_TIMEOUT) -race $(PACKAGES)

# To debug strange linter errors, uncomment
# DEBUG_LINTERS="--debug"

local-lint:
	golangci-lint $(DEBUG_LINTERS) run $(LINT_PACKAGES)

# Lint everything by default but ok to "make lint LINT_PACKAGES=./fhttp"
LINT_PACKAGES:=./...
lint:
	docker run -v $(CURDIR):/go/src/fortio.org/dnsping $(BUILD_IMAGE) bash -c \
		"cd /go/src/fortio.org/dnsping \
		&& time make local-lint DEBUG_LINTERS=\"$(DEBUG_LINTERS)\" LINT_PACKAGES=\"$(LINT_PACKAGES)\""

coverage: dependencies
	./.circleci/coverage.sh
	curl -s https://codecov.io/bash | bash

# Docker: Pushes the combo image and the smaller image(s)
all: test go-install lint docker-version docker-push-internal
	@for img in $(IMAGES); do \
		$(MAKE) docker-push-internal IMAGE=.$$img TAG=$(TAG); \
	done

# When changing the build image, this Makefile should be edited first
# (bump BUILD_IMAGE_TAG), also change this list if the image is used in
# more places.
FILES_WITH_IMAGE:= .circleci/config.yml Dockerfile Dockerfile.echosrv \
	Dockerfile.test Dockerfile.fcurl release/Dockerfile.in Webtest.sh

docker-internal: dependencies
	@echo "### Now building $(DOCKER_TAG)"
	docker build -f Dockerfile$(IMAGE) -t $(DOCKER_TAG) .

docker-push-internal: docker-internal
	@echo "### Now pushing $(DOCKER_TAG)"
	docker push $(DOCKER_TAG)

release:
	release/release.sh

.PHONY: all docker-internal docker-push-internal docker-version test dependencies

.PHONY: go-install lint install-linters coverage webtest release-test update-build-image

.PHONY: local-lint update-build-image-tag release pull certs certs-clean

# Targets used for official builds (initially from Dockerfile)
OFFICIAL_BIN := ../fortio.bin
GOOS := 
GO_BIN := go
VERSION ?= $(shell git describe --tags --match 'v*' --dirty)
# Main/default binary to build: (can be changed to build fcurl or echosrv instead)
OFFICIAL_TARGET := fortio.org/dnsping

LINK_FLAGS="-s -X main.Version=$(VERSION)"

.PHONY: official-build official-build-version official-build-clean

official-build:
	$(GO_BIN) version
	CGO_ENABLED=0 GOOS=$(GOOS) $(GO_BIN) build -a -ldflags $(LINK_FLAGS) -o $(OFFICIAL_BIN) $(OFFICIAL_TARGET)
	
official-build-version: official-build
	$(OFFICIAL_BIN) version

official-build-clean:
	-$(RM) $(OFFICIAL_BIN)

# Create a complete source tree with naming matching debian package conventions
TAR ?= tar # on macos need gtar to get --owner
DIST_VERSION ?= $(shell echo $(GIT_TAG) | sed -e "s/^v//")
DIST_PATH:=release/dnsping_$(DIST_VERSION).orig.tar

# Install target more compatible with standard gnu/debian practices. Uses DESTDIR as staging prefix

install: official-install

.PHONY: install official-install

BIN_INSTALL_DIR = $(DESTDIR)/usr/bin
LIB_INSTALL_DIR = $(DESTDIR)$(LIB_DIR)
MAN_INSTALL_DIR = $(DESTDIR)/usr/share/man/man1
BIN_INSTALL_EXEC = dnsping

official-install: official-build-clean official-build-version
	-mkdir -p $(BIN_INSTALL_DIR) $(LIB_INSTALL_DIR) $(MAN_INSTALL_DIR)
	cp $(OFFICIAL_BIN) $(BIN_INSTALL_DIR)/$(BIN_INSTALL_EXEC)
	cp dnsping.1 $(MAN_INSTALL_DIR)
