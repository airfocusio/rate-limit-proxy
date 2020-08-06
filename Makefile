.PHONY: run test test-watch build push redis

# The binary to build (just the basename).
BIN := rate-limit-proxy

# Where to push the docker image.
REGISTRY ?= choffmeister

IMAGE := $(REGISTRY)/$(BIN)-amd64

# This version-strategy uses git tags to set the version string
VERSION := $(shell git describe --tags --always --dirty)

run: redis
	cd src && go run . --config ../example.config.yaml --target https://golang.org

test: redis
	cd src && go test -v

test-watch: redis
	cd src && watch -n1 go test -v

build:
	docker build -t $(IMAGE):$(VERSION) .

push: build
	docker push $(IMAGE):$(VERSION)
	docker tag $(IMAGE):$(VERSION) $(IMAGE):latest
	docker push $(IMAGE):latest

redis:
	docker-compose up -d
	sleep 1