EXAMPLE=build-and-push

help: ## list allowed targets
	@echo 
	@echo " Executing example ${EXAMPLE}"
	@echo 
	@grep -E '^[a-zA-Z1-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf " \033[32m%-20s\033[0m %s\n", $$1, $$2}'
	@echo 

start: cleanup generate-certs ## start docker registry
	@echo
	@echo " Start environment"
	@echo
	@docker compose -f ../../test/docker-compose.yml up -d --build

cleanup: ## cleanup example environment
	@echo
	@echo " Cleanup environment"
	@echo
	@docker compose -f ../../test/docker-compose.yml down -v --remove-orphans --timeout 3

generate-certs: cleanup-certs ## generate certificate for go-docker-builder.test
	@echo
	@echo " Generate certificates"
	@echo
	@docker compose -f ../../test/docker-compose.yml run --rm openssl req -newkey rsa:2048 -nodes -keyout go-docker-builder.test.key -out go-docker-builder.test.csr -config /root/ssl/go-docker-builder.test.cnf
	@docker compose -f ../../test/docker-compose.yml run --rm openssl  x509 -signkey go-docker-builder.test.key -in go-docker-builder.test.csr -req -days 365 -out go-docker-builder.test.crt -extensions req_ext -extfile /root/ssl/go-docker-builder.test.cnf

cleanup-certs: ## cleanup certificates
	@echo
	@echo " Cleanup certificates"
	@echo
	@docker compose -f ../../test/docker-compose.yml run --rm --entrypoint /bin/sh openssl -c 'rm -rf $$(ls)'

prepare: ## prepare docker images required to run the example or test
	@echo
	@echo " Prepare"
	@echo
	@docker compose -f ../../test/docker-compose.yml run client sh -c 'docker pull base-registry.go-docker-builder.test:5000/alpine:3.13 \
		|| ( docker pull alpine:3.13 \
		&& docker tag alpine:3.13 base-registry.go-docker-builder.test:5000/alpine:3.13 \
		&& docker push base-registry.go-docker-builder.test:5000/alpine:3.13)'

example: start prepare ## executes the examples
	@echo
	@echo " Example"
	@echo
	@docker compose -f ../../test/docker-compose.yml run -w /app/examples/${EXAMPLE} client go run ${EXAMPLE}.go -test
	@docker compose -f ../../test/docker-compose.yml down -v --remove-orphans --timeout 3

test: start prepare ## executes functional test
	@echo
	@echo " Test"
	@echo
	@docker compose -f ../../test/docker-compose.yml run -w /app/examples/${EXAMPLE} client go test . -count=1
	@docker compose -f ../../test/docker-compose.yml down -v --remove-orphans --timeout 3

logs: ## show services logs
	@echo
	@echo " Logs"
	@echo
	@docker compose -f ../../test/docker-compose.yml logs
