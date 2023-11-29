EXAMPLE := $(shell basename ${PWD})
GREEN := "\\033[32m"
NC := "\\033[0m"
DOCKER_COMPOSE_VERSION := $(shell docker compose version > /dev/null 2>&1 && echo "docker compose" || (which docker-compose > /dev/null 2>&1 && echo "docker-compose" || (echo "docker compose not found. Aborting." >&2; exit 1)))

help: ## list allowed targets
	@echo 
	@echo "$(GREEN) Executing example $(EXAMPLE)$(NC)"
	@echo 
	@grep -E '^[a-zA-Z1-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf " \033[32m%-20s\033[0m %s\n", $$1, $$2}'
	@echo 

start: cleanup generate-keys generate-certs ## Starts the stack
	@echo
	@echo "$(GREEN) Starting the stack to run $(EXAMPLE)$(NC)"
	@echo
	@$(DOCKER_COMPOSE_VERSION) --project-name $(EXAMPLE) --file ../../test/docker-compose.yml up --detach --build

cleanup: ## Stops the stack
	@echo
	@echo "$(GREEN) Stopping the stack to run $(EXAMPLE)$(NC)"
	@echo
	@$(DOCKER_COMPOSE_VERSION) --project-name $(EXAMPLE) --file ../../test/docker-compose.yml down --volumes --remove-orphans --timeout 3

generate-keys: cleanup-keys ## Generates an ssh key pair required to autheneticate to git server
	@echo
	@echo "$(GREEN) Generate ssh key pair for $(EXAMPLE)$(NC)"
	@echo
	@docker compose --project-name ${EXAMPLE}  --file ../../test/docker-compose.yml run --rm openssh -t rsa -q -N "password" -f id_rsa -C "apenella@go-docker-builder.test"

cleanup-keys: ## Cleans up the ssh key pair
	@echo
	@echo "$(GREEN) Cleanup ssh key pair for $(EXAMPLE)$(NC)"
	@echo
	@docker compose --project-name ${EXAMPLE}  --file ../../test/docker-compose.yml run --rm --entrypoint /bin/sh openssh -c 'rm -rf $$(ls)'

generate-certs: cleanup-certs ## Generates certificate for go-docker-builder.test
	@echo
	@echo "$(GREEN) Generate certificates for $(EXAMPLE)$(NC)"
	@echo
	@$(DOCKER_COMPOSE_VERSION) --project-name $(EXAMPLE) --file ../../test/docker-compose.yml run --rm openssl req -newkey rsa:2048 -nodes -keyout go-docker-builder.test.key -out go-docker-builder.test.csr -config /root/ssl/go-docker-builder.test.cnf
	@$(DOCKER_COMPOSE_VERSION) --project-name $(EXAMPLE) --file ../../test/docker-compose.yml run --rm openssl  x509 -signkey go-docker-builder.test.key -in go-docker-builder.test.csr -req -days 365 -out go-docker-builder.test.crt -extensions req_ext -extfile /root/ssl/go-docker-builder.test.cnf

cleanup-certs: ## Cleans up certificates
	@echo
	@echo "$(GREEN) Cleanup certificates $(EXAMPLE)$(NC)"
	@echo
	@$(DOCKER_COMPOSE_VERSION) --project-name $(EXAMPLE) --file ../../test/docker-compose.yml run --rm --entrypoint /bin/sh openssl -c 'rm -rf $$(ls)'

prepare: start ## Prepares Docker images required to run the example or test
	@echo
	@echo "$(GREEN) Prepare Docker images for $(EXAMPLE)$(NC)"
	@echo

	@echo " [$(EXAMPLE)]$(GREEN) Waiting for dockerd$(NC)"
	@for i in $$(seq 1 5); \
	do $(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec client /usr/local/bin/wait-for-dockerd.sh && RC=0 && break || RC=$$? && $(DOCKER_COMPOSE_VERSION) --project-name $(EXAMPLE) --file ../../test/docker-compose.yml restart client --timeout 3; \
	done || { echo " Error synchronizing $(EXAMPLE)"; exit $$RC; }
	@echo

	@echo " [$(EXAMPLE)]$(GREEN) Pulling Docker images$(NC)"
	@$(DOCKER_COMPOSE_VERSION) --project-name $(EXAMPLE) --file ../../test/docker-compose.yml exec client sh -c 'docker pull base-registry.go-docker-builder.test:5000/alpine:3.13 \
		|| ( docker pull alpine:3.13 \
		&& docker tag alpine:3.13 base-registry.go-docker-builder.test:5000/alpine:3.13 \
		&& docker push base-registry.go-docker-builder.test:5000/alpine:3.13)'
	@$(DOCKER_COMPOSE_VERSION) --project-name $(EXAMPLE) --file ../../test/docker-compose.yml exec client sh -c 'docker pull base-registry.go-docker-builder.test:5000/golang:1.15-alpine \
		|| ( docker pull golang:1.15-alpine \
		&& docker tag golang:1.15-alpine base-registry.go-docker-builder.test:5000/golang:1.15-alpine \
		&& docker push base-registry.go-docker-builder.test:5000/golang:1.15-alpine)'
	@echo


example: prepare ## Executes the examples
	@echo
	@echo "$(GREEN) Run example $(EXAMPLE)$(NC)"
	@echo

	@echo " [$(EXAMPLE)]$(GREEN) Waiting for dockerd$(NC)"
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec client /usr/local/bin/wait-for-dockerd.sh || { echo " Error synchronizing $(EXAMPLE)"; exit 1; }
	@echo

	@echo " [$(EXAMPLE)]$(GREEN) Executing the example$(NC)"
	@$(DOCKER_COMPOSE_VERSION) --project-name $(EXAMPLE) --file ../../test/docker-compose.yml exec --workdir /app/examples/$(EXAMPLE) client go run $(EXAMPLE).go -test
	@echo

	@echo " [$(EXAMPLE)]$(GREEN) Cleanup the stack$(NC)"
	@$(DOCKER_COMPOSE_VERSION) --project-name $(EXAMPLE) --file ../../test/docker-compose.yml down --volumes --remove-orphans --timeout 3
	@echo

test: prepare ## Executes functional test
	@echo
	@echo "$(GREEN) Testing $(EXAMPLE)$(NC)"
	@echo

	@echo " [$(EXAMPLE)]$(GREEN) Waiting for dockerd$(NC)"
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec client /usr/local/bin/wait-for-dockerd.sh || { echo " Error synchronizing $(EXAMPLE)"; exit 1; }
	@echo

	@echo " [$(EXAMPLE)]$(GREEN) Running tests$(NC)"
	@$(DOCKER_COMPOSE_VERSION) --project-name $(EXAMPLE) --file ../../test/docker-compose.yml exec --workdir /app/examples/$(EXAMPLE) client go test . -count=1 || RC=1; \
	echo; echo " [$(EXAMPLE)]$(GREEN) Cleanup the stack$(NC)"; echo && \
	$(DOCKER_COMPOSE_VERSION) --project-name $(EXAMPLE) --file ../../test/docker-compose.yml down --volumes --remove-orphans --timeout 3; \
	exit $$RC;
	@echo

logs: ## Show services logs
	@echo
	@echo "$(GREEN) Logs for $(EXAMPLE)$(NC)"
	@echo
	@$(DOCKER_COMPOSE_VERSION) --project-name $(EXAMPLE) --file ../../test/docker-compose.yml logs

follow-logs: ## Follow services logs
	@echo
	@echo "$(GREEN) Follow logs for $(EXAMPLE)$(NC)"
	@echo
	@$(DOCKER_COMPOSE_VERSION) --project-name $(EXAMPLE) --file ../../test/docker-compose.yml logs --follow
