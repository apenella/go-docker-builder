BASE_FUNCTIONAL_FOLDER=examples


help: ## list allowed targets
	@grep -E '^[a-zA-Z1-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf " \033[32m%-20s\033[0m %s\n", $$1, $$2}'
	@echo 

test: unit-test functional-test ## Run all test

functional-test: ## Run functional test
	@echo
	@echo " Run functional test: build-and-push"
	@echo 	  
	cd ${BASE_FUNCTIONAL_FOLDER}/build-and-push && $(MAKE) test ; cd -

	@echo
	@echo " Run functional test: build-and-push-join-context"
	@echo 	  
	cd ${BASE_FUNCTIONAL_FOLDER}/build-and-push-join-context && $(MAKE) test ; cd -

	@echo
	@echo " Run functional test: build-git-context"
	@echo 	  
	cd ${BASE_FUNCTIONAL_FOLDER}/build-git-context && $(MAKE) test ; cd -

	@echo
	@echo " Run functional test: build-git-context-auth"
	@echo 	  
	cd ${BASE_FUNCTIONAL_FOLDER}/build-git-context-auth && $(MAKE) test ; cd -

	@echo
	@echo " Run functional test: build-git-path"
	@echo 	  
	cd ${BASE_FUNCTIONAL_FOLDER}/build-git-path && $(MAKE) test ; cd -

	@echo
	@echo " Run functional test: copy-remote"
	@echo 	  
	cd ${BASE_FUNCTIONAL_FOLDER}/copy-remote && $(MAKE) test ; cd -


unit-test: ## Run unitary test
	@echo
	@echo " Run unit test"
	@echo
	go test ./pkg/... -v -cover -count=1