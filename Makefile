network:
	docker network create rarecircles

build:
	docker-compose -f docker-compose.yml build

up:
	docker-compose -f docker-compose.yml up

down:
	docker-compose -f docker-compose.yml down

logs:
	docker-compose -f docker-compose.yml logs -f

test:
	go test -v ./...