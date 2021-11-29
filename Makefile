# Looks ugly in your editor, but the formatting directives used are:
# * \33[1m -- start bold
# * \33[4m -- start underline
# * \33[1;4m -- start bold and underline
# * \33[0m -- reset formatting
# See: https://misc.flogisoft.com/bash/tip_colors_and_formatting
help:
	@printf "\33[1mAvailable make commands\33[0m:\n \
\33[1;4m\33[0m\n \
## \33[1mRunning!\33[0m\n \
 * \33[1;4mhelp\33[0m: display this message\n \
 * \33[1;4mnetsetup\33[0m: create docker network\n \
 * \33[1;4mall\33[0m: build and starts containers\n \
 * \33[1;4mbuild\33[0m: builds service image\n \
 * \33[1;4mup\33[0m: starts containers\n \
 * \33[1;4mdown\33[0m: stops containers\n \
 * \33[1;4mlogs\33[0m: show containers logs\n \
 * \33[1;4mtest\33[0m: run unit tests\n \
\33[1;4m\33[0m\n"

default: help

setup:
	docker network create rarecircles

build:
	docker-compose -f docker-compose.yml build

up:
	docker-compose -f docker-compose.yml up

down:
	docker-compose -f docker-compose.yml down

restart: down up

all: build up

logs:
	docker-compose -f docker-compose.yml logs -f

test:
	go test -v ./...
