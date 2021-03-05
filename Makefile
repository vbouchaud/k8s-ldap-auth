# COLORS
TARGET_MAX_CHAR_NUM := 10
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

# The binary to build (just the basename).
PWD := $(shell pwd)
NOW := $(shell date +%s)
BIN := k8s-ldap-auth

ORG ?= registry.aegir.bouchaud.org
NAMESPACE := legion/kubernetes
PKG := bouchaud.org/${NAMESPACE}/${BIN}
PLATFORM ?= "linux/arm/v7,linux/arm64/v8,linux/amd64"
GO ?= go
GOFMT ?= gofmt -s
GOFILES := $(shell find . -name "*.go" -type f)
GOVERSION := $(shell go version | sed -r 's/go version go(.+)\s.+/\1/')
PACKAGES ?= $(shell $(GO) list ./...)

# This version-strategy uses git tags to set the version string
GIT_TAG := $(shell git describe --tags --always --dirty || echo unsupported)
GIT_COMMIT := $(shell git rev-parse --short HEAD || echo unsupported)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
GIT_BRANCH_CLEAN := $(shell echo $(GIT_BRANCH) | sed -e "s/[^[:alnum:]]/-/g")
BUILDTIME := $(shell date -u +"%FT%TZ%:z")
ARCH := $(shell uname -m)
TAG ?= $(GIT_TAG)

.PHONY: fmt fmt-check vet test test-coverage cover install hooks docker tag push dev help
default: help

## Format go source code
fmt:
	$(GOFMT) -w $(GOFILES)

## Check if source code is formatted correctly
fmt-check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

## Check source code for common errors
vet:
	$(GO) vet ${PACKAGES}

## Execute unit tests
test:
	$(GO) test ${PACKAGES}

## Execute unit tests & compute coverage
test-coverage:
	$(GO) test -coverprofile=coverage.out ${PACKAGES}

## Compute coverage
cover: test-coverage
	$(GO) tool cover -html=coverage.out

## Install dependencies used for development
install: hooks
	$(GO) mod download

## Install git hooks for post-checkout & pre-commit
hooks:
	@cp -f ./scripts/post-checkout .git/hooks
	@cp -f ./scripts/pre-commit .git/hooks
	@chmod +x .git/hooks/post-checkout
	@chmod +x .git/hooks/pre-commit

## Build the docker images
docker:
	@docker buildx build \
		--push \
		--build-arg COMMITHASH="$(GIT_COMMIT)" \
		--build-arg BUILDTIME="$(BUILDTIME)" \
		--build-arg VERSION="$(GIT_TAG)" \
		--build-arg PKG="$(PKG)" \
		--build-arg APPNAME="$(BIN)" \
		--platform $(PLATFORM) \
		--tag $(ORG)/$(BIN):$(TAG) \
		--tag $(ORG)/$(BIN):latest \
		.

## Dev build outside of docker, not stripped
dev:
	$(GO) build \
		-o $(BIN) -ldflags "\
				-X $(PKG)/version.APPNAME=$(BIN) \
				-X $(PKG)/version.VERSION=$(GIT_TAG) \
				-X $(PKG)/version.GOVERSION=$(GOVERSION) \
				-X $(PKG)/version.BUILDTIME=$(BUILDTIME) \
				-X $(PKG)/version.COMMITHASH=$(GIT_COMMIT)"

## Print this help message
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
