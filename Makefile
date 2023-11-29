BASE_FUNCTIONAL_FOLDER=examples


help: ## list allowed targets
	@grep -E '^[a-zA-Z1-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf " \033[32m%-20s\033[0m %s\n", $$1, $$2}'
	@echo 

test: unit-test functional-test ## Run all test

functional-test: build-and-push-test build-and-push-join-context-test build-git-context-test build-git-context-auth-test build-path-context-test copy-remote-test ## Run functional tests

build-and-push-test: ## Execute functional test build-and-push
	@echo
	@echo " Run functional test: build-and-push"
	@echo
	@RC=0; \
	cd ${BASE_FUNCTIONAL_FOLDER}/build-and-push && $(MAKE) test || RC=1; \
	cd -; \
	exit $$RC;

build-and-push-join-context-test: ## Execute functional test build-and-push-join-context
	@echo
	@echo " Run functional test: build-and-push-join-context"
	@echo
	@RC=0; \
	cd ${BASE_FUNCTIONAL_FOLDER}/build-and-push-join-context && $(MAKE) test || RC=1; \
	cd -; \
	exit $$RC;

build-git-context-test: ## Execute functional test build-git-context
	@echo
	@echo " Run functional test: build-git-context"
	@echo
	@RC=0; \
	cd ${BASE_FUNCTIONAL_FOLDER}/build-git-context && $(MAKE) test || RC=1; \
	cd -; \
	exit $$RC;

build-git-context-auth-test: ## Execute functional test build-git-context-auth
	@echo
	@echo " Run functional test: build-git-context-auth"
	@echo
	@RC=0; \
	cd ${BASE_FUNCTIONAL_FOLDER}/build-git-context-auth && $(MAKE) test || RC=1; \
	cd -; \
	exit $$RC;

build-path-context-test: ## Execute functional test build-path-context
	@echo
	@echo " Run functional test: build-path-context"
	@echo
	@RC=0; \
	cd ${BASE_FUNCTIONAL_FOLDER}/build-path-context && $(MAKE) test || RC=1; \
	cd -; \
	exit $$RC;

copy-remote-test: ## Execute functional test copy-remote
	@echo
	@echo " Run functional test: copy-remote"
	@echo
	@RC=0; \
	cd ${BASE_FUNCTIONAL_FOLDER}/copy-remote && $(MAKE) test || RC=1; \
	cd -; \
	exit $$RC;

unit-test: ## Run unitary tests
	@echo
	@echo " Run unit test"
	@echo
	go test ./pkg/... -cover -count=1
