EXAMPLE=build-and-push

.PHONY: help start example cleanup

help: ## list allowed targets
	@grep -E '^[a-zA-Z1-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-10s\033[0m %s\n", $$1, $$2}'

start: ## start docker registry
	docker-compose -f ../../test/docker-compose.yml up -d --build registry

example: start ## executes the examples
	docker-compose -f ../../test/docker-compose.yml run -w /app/examples/${EXAMPLE} client go run ${EXAMPLE}.go

logs: ## show services logs
	docker-compose -f ../../test/docker-compose.yml logs

cleanup: ## cleanup example environment
	docker-compose -f ../../test/docker-compose.yml down -v --remove-orphans