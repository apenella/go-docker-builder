help: ## list allowed targets
	@grep -E '^[a-zA-Z1-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf " \033[32m%-20s\033[0m %s\n", $$1, $$2}'
	@echo 

start: cleanup generate-keys generate-certs ## start docker registry
	docker compose up --detach --build registry base-registry gitserver

cleanup: cleanup-certs ## cleanup example environment
	docker compose down --volumes --remove-orphans --timeout 3

generate-keys: cleanup-keys ## generate an ssh key pair required to autheneticate to git server
	@docker compose run --rm openssh -t rsa -q -N "password" -f id_rsa -C "apenella@go-docker-builder.test"

cleanup-keys: ## cleanup the ssh key pair
	@docker compose run --rm --entrypoint /bin/sh openssh -c 'rm -rf $$(ls)'

generate-certs: cleanup-certs ## generate certificate for go-docker-builder.test
	@docker compose run --rm openssl req -newkey rsa:2048 -nodes -keyout go-docker-builder.test.key -out go-docker-builder.test.csr -config /root/ssl/go-docker-builder.test.cnf
	@docker compose run --rm openssl  x509 -signkey go-docker-builder.test.key -in go-docker-builder.test.csr -req -days 365 -out go-docker-builder.test.crt -extensions req_ext -extfile /root/ssl/go-docker-builder.test.cnf

cleanup-certs: ## cleanup certificates
	@docker compose run --rm --entrypoint /bin/sh openssl -c 'rm -rf $$(ls)'

attach-client: start ## attach to client container
	docker compose run --build --workdir /app/examples client sh