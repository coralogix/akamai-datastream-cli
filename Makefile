.PHONY: build deploy

export VERSION ?= v1.0.0
export GIT_COMMIT ?= $(shell git rev-list -1 HEAD 2> /dev/null || echo "unknown")
export GO_VERSION := $(shell go version)
export CGO_ENABLED=0

export PREFIX = coralogixrepo
export IMAGE = fluentd-coralogix-akamai
export TAG ?= 1.0.0

export COMPOSE_PROJECT_NAME=fluentd-coralogix-akamai
export COMPOSE_FILE=deploy/docker-compose.yml

vendor:
	@go mod vendor

fmt:
	@go fmt .

vet:
	@go vet .

build: clean vendor fmt vet
	@go build -a -tags netgo -ldflags "-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X 'main.GoVersion=$(GO_VERSION)' -X main.BuildDate=`date -u '+%Y-%m-%dT%H:%M:%SZ'`" -mod=vendor -o bin/akamai-datastream-cli .

clean:
	@rm -f bin/akamai-datastream-cli

image:
	@docker build \
		--tag $(PREFIX)/$(IMAGE):latest \
		--tag $(PREFIX)/$(IMAGE):$(TAG) \
		--build-arg VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--file build/Dockerfile \
		.

deploy:
	@docker-compose up -d

undeploy:
	@docker-compose down

all: build