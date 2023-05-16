# ########################################################## #
# Makefile for Golang Project
# Includes cross-compiling, installation, cleanup
# ########################################################## #

# Default Goal of the makefile is to show the help
.DEFAULT_GOAL := help

# Sets default shell to Bash
#SHELL := /bin/bash

# Check for required command tools to build or stop immediately
EXECUTABLES = go find pwd awk docker
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))
# Setting up variables for the make targets
COMMIT_USER=$(shell git log -1 --pretty=format:'%an')
BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
BUILD=$(shell git rev-parse HEAD)
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BINARY=polygon-client
SRC_DIR:=$(ROOT_DIR)
OUT_DIR:=$(ROOT_DIR)/bin
LDFLAGS=-ldflags "-X main.Build=$(BUILD) -s -w"
ARCHS = linux/amd64 linux/arm64 darwin/amd64
# Indicates that the following targets have no physical files
.PHONY: all clean build docker build-multi-arch help test vet tooling gotest-junit ginkgo

all: proto build ## Runs everything

clean: ## Clean generated sources
	rm -rf $(SRC_DIR)/bin

build: ## Build the application
	cd $(SRC_DIR) && go mod tidy
	cd $(SRC_DIR) && CGO_ENABLED=0 go build $(LDFLAGS) -o $(OUT_DIR)/$(BINARY)

vet: # Vet the application code
	cd $(SRC_DIR) && go vet ./...
test: ## Test the application
	cd $(SRC_DIR) && go test -v -cover ./...

docker: ## Docker build
	docker build --build-arg BUILD_REVISION=$(BUILD) -t rafaribe/$(BINARY) .
upgrade: ## Upgrade go dependencies, please build after to make sure it still works.
	cd $(ROOT_DIR) && go get -u  && go mod tidy

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)