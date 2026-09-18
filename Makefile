BASE_FUNCTIONAL_FOLDER=examples
GOLANG_VERSION ?= 1.26
GOLANG_IMAGE := golang:$(GOLANG_VERSION)-alpine


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

unit-test: ## Run unitary tests inside Go container
	@echo
	@echo " Run unit test (golang:$(GOLANG_VERSION)-alpine)"
	@echo
	docker run --rm -v $(CURDIR):/app -w /app $(GOLANG_IMAGE) go test ./pkg/... -cover -count=1

build: ## Build all packages inside Go container
	@echo
	@echo " Build (golang:$(GOLANG_VERSION)-alpine)"
	@echo
	docker run --rm -v $(CURDIR):/app -w /app $(GOLANG_IMAGE) go build ./...

vet: ## Run go vet inside Go container
	@echo
	@echo " Run go vet (golang:$(GOLANG_VERSION)-alpine)"
	@echo
	docker run --rm -v $(CURDIR):/app -w /app $(GOLANG_IMAGE) go vet ./...

fmt-check: ## Check gofmt formatting inside Go container
	@echo
	@echo " Check formatting (golang:$(GOLANG_VERSION)-alpine)"
	@echo
	docker run --rm -v $(CURDIR):/app -w /app $(GOLANG_IMAGE) sh -c 'gofmt -l . | grep . && echo "Formatting issues found" && exit 1 || echo "Formatting OK"'

tidy: ## Run go mod tidy inside Go container
	@echo
	@echo " Run go mod tidy (golang:$(GOLANG_VERSION)-alpine)"
	@echo
	docker run --rm -v $(CURDIR):/app -w /app $(GOLANG_IMAGE) go mod tidy

tidy-check: ## Verify go.mod/go.sum are tidy inside Go container
	@echo
	@echo " Check go mod tidy (golang:$(GOLANG_VERSION)-alpine)"
	@echo
	@cp go.mod /tmp/go.mod.before-check && cp go.sum /tmp/go.sum.before-check; \
	docker run --rm -v $(CURDIR):/app -w /app $(GOLANG_IMAGE) go mod tidy; \
	if diff -q /tmp/go.mod.before-check go.mod >/dev/null && diff -q /tmp/go.sum.before-check go.sum >/dev/null; then \
		echo "Tidy OK"; rm -f /tmp/go.mod.before-check /tmp/go.sum.before-check; \
	else \
		echo "go.mod/go.sum are not tidy"; diff /tmp/go.mod.before-check go.mod || true; diff /tmp/go.sum.before-check go.sum || true; rm -f /tmp/go.mod.before-check /tmp/go.sum.before-check; exit 1; \
	fi

static-analysis: vet fmt-check ## Run static analysis (vet + fmt) inside Go container

go-version: ## Show Go version inside container
	docker run --rm $(GOLANG_IMAGE) go version
