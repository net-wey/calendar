.PHONY: build run run-main docker-up docker-down clean clean-bin lint test check

build:
	go build -o bin/main ./cmd/main

run-main:
	go run cmd/main/main.go

run: build
	./bin/main

docker-up:
	docker-compose -f docker/docker-compose.yml up -d

docker-down:
	docker-compose -f docker/docker-compose.yml down

clean: docker-down
	cd docker && docker-compose down -v

clean-bin:
	rm -rf bin/


lint:
	golangci-lint run ./...

test:
	go test -v -race ./...


check: lint test build