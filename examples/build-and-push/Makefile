EXAMPLE=build-and-push

help: ## list allowed targets
	@echo 
	@echo " Executing example ${EXAMPLE}"
	@echo 
	@grep -E '^[a-zA-Z1-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf " \033[32m%-20s\033[0m %s\n", $$1, $$2}'
	@echo 

start: cleanup ## start docker registry
	docker-compose -f ../../test/docker-compose.yml up -d --build

cleanup: ## cleanup example environment
	docker-compose -f ../../test/docker-compose.yml down -v --remove-orphans

example: start ## executes the examples
	docker-compose -f ../../test/docker-compose.yml run -w /app/examples/${EXAMPLE} client go run ${EXAMPLE}.go

logs: ## show services logs
	docker-compose -f ../../test/docker-compose.yml logs
