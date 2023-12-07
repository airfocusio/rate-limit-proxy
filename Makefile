.PHONY: *

run: redis
	go run . --config ./example.config.yaml --listen 127.0.0.1:8080 --target https://golang.org

test: redis
	go test -v ./internal

test-watch: redis
	watch -n1 go test -v

redis:
	docker-compose up -d
	sleep 1

build:
	goreleaser release --clean --skip=publish --snapshot

release:
	goreleaser release --clean

trivy: build
	trivy image ghcr.io/airfocusio/rate-limit-proxy:0.0.0-dev-amd64
