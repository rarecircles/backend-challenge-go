dev:
	go run ./cmd/challenge

build:
	docker-compose -f docker-compose.yml build

up:
	docker-compose -f docker-compose.yml up

down:
	docker-compose -f docker-compose.yml down

all: build up

test:
	go test -v ./...
	staticcheck ./...