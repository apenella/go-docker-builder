COMPOSE_FILE=test/docker-compose.yml
PROJECT=go-docker-builder

help:
	@grep -E '^[a-zA-Z1-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

up: ## Set up project
	docker-compose -p ${PROJECT} -f ${COMPOSE_FILE} up -d --build

down: ## Tear down project
	docker-compose -p ${PROJECT} -f ${COMPOSE_FILE} down --remove-orphans

client-sh: up ## attach to client service
	docker-compose -p ${PROJECT} -f ${COMPOSE_FILE} exec client sh

test: up ## attach to client service
	docker-compose -p ${PROJECT} -f ${COMPOSE_FILE} exec client go test ./pkg/...1