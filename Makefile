.PHONY: help
help:
	@echo ${print_help}

.PHONY: deploy-challenge
deploy-challenge:
	docker-compose up -d

.PHONY: challenge-down
challenge-down:
	docker-compose down

.PHONY: challenge-logs
challenge-logs:
	docker logs -f backend-challenge-go_challenge-service_1

.PHONY: dal-logs
dal-logs:
	docker logs -f backend-challenge-go_postgres_1

.PHONY: cleanup
cleanup:
	docker-compose down
	docker-compose rm
	docker rmi backend-challenge-go_challenge-service
	docker volume rm backend-challenge-go_postgres

# Only have one unit test in utils
.PHONY: unit-test
unit-test:
	@go test -failfast ./utils/



define print_help
"Commands: \n\n\
	make help				: Display this message \n\
	make deploy-challenge			: Deploys challenge service \n\
	make challenge-down			: Downs all challenge related containers \n\
	make challenge-logs			: Displays logs specific to challenge service \n\
	make dal-logs				: Displays logs specific to Postgres \n\
	make cleanup				: Removes docker containers, image, volumes \n\
	make unit-test				: Runs unit tests \n\
"
endef