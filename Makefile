
PROJECT_ROOT            := github.com/infobloxopen/atlas-cli/github.release.notes
BUILD_PATH              := bin
DOCKERFILE_PATH         := $(CURDIR)/docker

AWS_ACCESS_KEY_ID?=`aws configure get aws_access_key_id`
AWS_SECRET_ACCESS_KEY?=`aws configure get aws_secret_access_key`
AWS_REGION?=`aws configure get region`
DOCKER_ENV := -e AWS_REGION=$(AWS_REGION) -e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) -e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY)

# configuration for image names
USERNAME                := $(USER)
GIT_COMMIT              := $(shell git describe --dirty=-unsupported --always --tags || echo pre-commit)
IMAGE_VERSION           ?= $(GIT_COMMIT)
IMAGE_REGISTRY          ?= infoblox

IMAGE_NAME              ?= github.release.notes
# configuration for server binary and image
SERVER_BINARY           := $(BUILD_PATH)/server
SERVER_PATH             := $(PROJECT_ROOT)/cmd/server
SERVER_IMAGE            := $(IMAGE_REGISTRY)/$(IMAGE_NAME)
SERVER_DOCKERFILE       := $(DOCKERFILE_PATH)/Dockerfile

# configuration for the protobuf gentool
SRCROOT_ON_HOST         := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
SRCROOT_IN_CONTAINER    := /go/src/$(PROJECT_ROOT)
DOCKER_RUNNER           ?= docker run --rm -u `id -u`:`id -g` -e GOCACHE=/go -e CGO_ENABLED=0
DOCKER_RUNNER           += -v $(SRCROOT_ON_HOST):$(SRCROOT_IN_CONTAINER)
DOCKER_GENERATOR        := infoblox/atlas-gentool:latest
GENERATOR               := $(DOCKER_RUNNER) $(DOCKER_GENERATOR)

# configuration for building on host machine
export GOFLAGS          ?= -mod=vendor
GO_CACHE                := -pkgdir $(BUILD_PATH)/go-cache
GO_BUILD_FLAGS          ?= $(GO_CACHE) -i -v
GO_TEST_FLAGS           ?= -v -cover
GO_PACKAGES             := $(shell go list ./... | grep -v vendor)

# this file was generated by atlas-cli
# changes to this file will be overriden by atlas-cli update
#
PROJECT_ROOT            ?= $(PWD)
BUILD_PATH              ?= bin
DOCKERFILE_PATH         ?= $(CURDIR)/docker
GO_IMAGE                ?= golang:1.15-alpine
GO_RUNNER               ?= $(DOCKER_RUNNER) $(GO_IMAGE)

# configuration for image names
USERNAME                ?= $(USER)
GIT_COMMIT              ?= $(shell git describe --dirty=-unsupported --always --tags || echo pre-commit)
IMAGE_VERSION           ?= $(GIT_COMMIT)-$(USERNAME)

GO_MOD = go.mod

.PHONY all: all-atlas
all-atlas: vendor-atlas docker-atlas

.PHONY fmt: fmt-atlas
fmt-atlas:
	@$(GO_RUNNER) go fmt $(GO_PACKAGES)

.PHONY test: test-atlas
test-atlas: fmt-atlas
	$(GO_RUNNER) go test $(GO_TEST_FLAGS) $(GO_PACKAGES)

docker-atlas:
	@docker build -f $(SERVER_DOCKERFILE) -t $(SERVER_IMAGE):$(IMAGE_VERSION) .
	@docker image prune -f --filter label=stage=server-intermediate

.docker-$(IMAGE_NAME)-$(IMAGE_VERSION):
	$(MAKE) docker-atlas
	touch $@

.PHONY: docker
docker: .docker-$(IMAGE_NAME)-$(IMAGE_VERSION)

push-atlas: docker
ifndef IMAGE_REGISTRY
	@(echo "Please set IMAGE_REGISTRY variable in Makefile.vars to use push command"; exit 1)
else
	@docker push $(SERVER_IMAGE):$(IMAGE_VERSION)
endif

.push-$(IMAGE_NAME)-$(IMAGE_VERSION):
	$(MAKE) push-atlas
	touch $@

.PHONY: push
push: .push-$(IMAGE_NAME)-$(IMAGE_VERSION)

.PHONY vendor: vendor-atlas
vendor-atlas:
	@go mod tidy
	@go mod vendor
	@go mod download

.PHONY clean: clean-atlas
clean-atlas:
	@docker rmi -f $(shell docker images -q $(SERVER_IMAGE)) || true
	rm .push-* .docker-*
