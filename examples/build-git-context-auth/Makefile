EXAMPLE=build-git-context-auth

.PHONY: help start example cleanup

help: ## list allowed targets
	@echo 
	@echo " Executing example ${EXAMPLE}"
	@echo 
	@grep -E '^[a-zA-Z1-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-10s\033[0m %s\n", $$1, $$2}'

start: cleanup generate-keys ## start docker registry
	docker-compose -f ../../test/docker-compose.yml up -d --build

cleanup: ## cleanup example environment
	docker-compose -f ../../test/docker-compose.yml down -v --remove-orphans

generate-keys: cleanup-keys ## generate an ssh key pair required to autheneticate to git server
	docker-compose -f ../../test/docker-compose.yml run --rm openssh -t rsa -q -N "password" -f id_rsa -C "apenella@go-docker-builder.test"

cleanup-keys: ## cleanup the ssh key pair
	docker-compose -f ../../test/docker-compose.yml run --rm --entrypoint /bin/sh openssh -c 'rm -rf $$(ls)'
#	docker-compose -f ../../test/docker-compose.yml down

example: start ## executes the examples
	docker-compose -f ../../test/docker-compose.yml run -w /app/examples/${EXAMPLE} client go run ${EXAMPLE}.go

logs: ## show services logs
	docker-compose -f ../../test/docker-compose.yml logs
