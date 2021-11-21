.PHONY: *

run: redis
	go run . --config ./example.config.yaml --target https://golang.org

test: redis
	go test -v ./internal

test-watch: redis
	watch -n1 go test -v

redis:
	docker-compose up -d
	sleep 1
