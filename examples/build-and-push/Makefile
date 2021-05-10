EXAMPLE=build-and-push

help: ## list allowed targets
	@echo 
	@echo " Executing example ${EXAMPLE}"
	@echo 
	@grep -E '^[a-zA-Z1-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf " \033[32m%-20s\033[0m %s\n", $$1, $$2}'
	@echo 

start: cleanup generate-certs ## start docker registry
	docker-compose -f ../../test/docker-compose.yml up -d --build

cleanup: ## cleanup example environment
	docker-compose -f ../../test/docker-compose.yml down -v --remove-orphans

generate-certs: cleanup-certs ## generate certificate for go-docker-builder.test
	@docker-compose -f ../../test/docker-compose.yml run --rm openssl req -newkey rsa:2048 -nodes -keyout go-docker-builder.test.key -out go-docker-builder.test.csr -config /root/ssl/go-docker-builder.test.cnf
	@docker-compose -f ../../test/docker-compose.yml run --rm openssl  x509 -signkey go-docker-builder.test.key -in go-docker-builder.test.csr -req -days 365 -out go-docker-builder.test.crt -extensions req_ext -extfile /root/ssl/go-docker-builder.test.cnf

cleanup-certs: ## cleanup certificates
	@docker-compose -f ../../test/docker-compose.yml run --rm --entrypoint /bin/sh openssl -c 'rm -rf $$(ls)'

example: start ## executes the examples
	docker-compose -f ../../test/docker-compose.yml run -w /app/examples/${EXAMPLE} client go run ${EXAMPLE}.go

logs: ## show services logs
	docker-compose -f ../../test/docker-compose.yml logs