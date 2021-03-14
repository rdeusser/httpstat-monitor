# Setup name variables for the package/tool.
NAME := httpstat-monitor
PKG := github.com/rdeusser/$(NAME)
BUILD_PATH := $(PKG)/cmd/$(NAME)
GIT_COMMIT := $(PKG)/version
VERSION := $(shell grep -oE "[0-9]+[.][0-9]+[.][0-9]+" version/version.go)

SEMVER := patch

OLDPWD := $(PWD)
export OLDPWD

OUT_DIR := $(PWD)/bin

FILES_TO_FMT ?= $(shell find . -path ./vendor -prune -o -name '*.go' -print)

DOCKER_IMAGE_REPO ?= rdeusser/$(NAME)

GOBIN		   ?= $(firstword $(subst :, ,${GOPATH}))/bin
GO111MODULE	   ?= on
export GO111MODULE

GOIMPORTS_VERSION	      ?= master
GOIMPORTS		      ?= $(GOBIN)/goimports

GOLANGCILINT_VERSION	      ?= v1.31.0
GOLANGCILINT		      ?= $(GOBIN)/golangci-lint

.DEFAULT_GOAL := help

define fetch_go_bin_version
	@cd /tmp
	@go get $(1)@$(2)
	@cd -
endef

.PHONY: help
help: ## Display this help text.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nAvailable targets:\n"} /^[\/0-9a-zA-Z_-]+:.*?##/ { printf "  \x1b[32;01m%-20s\x1b[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: fmt
fmt: $(GOIMPORTS) ## Formats Go code including imports and cleans up noise.
	@echo ">> formatting code"
	@$(GOIMPORTS) -local github.com/rdeusser/$(NAME) -w $(FILES_TO_FMT)
	@echo ">> cleaning up noise"
	@find . -type f \( -name "*.md" -o -name "*.go" \) | SED_BIN="$(SED)" xargs scripts/cleanup-noise.sh

.PHONY: lint
lint: $(GOLANGCILINT) ## Run various static analysis tools against our code.
	@echo ">> linting all of the Go files"
	@$(GOLANGCILINT) run

.PHONY: generate
generate: $(BUF) ## Generates code.
	@echo ">> generating code"
	@go generate ./...

.PHONY: test
test: ## Runs all httpstat-monitor's unit tests.
	@echo ">> running unit tests"
	@go test -coverprofile=coverage.out $(shell go list ./...);

.PHONY: build
build: ## Build httpstat-monitor.
	@echo ">> building"
	@-CGO_ENABLED=0 \
		go build \
		-o $(OUT_DIR)/httpstat-monitor \
		$(BUILD_PATH)

.PHONY: install
install: build ## Build and install httpstat-monitor.
	@echo ">> installing"
	 mv $(OUT_DIR)/httpstat-monitor $(GOBIN)

.PHONY: docker-build
docker-build: ## Build docker image.
	@echo ">> building docker image"
	@docker build -t $(DOCKER_IMAGE_REPO):$(VERSION) .

.PHONY: deploy
deploy: docker-build ## Builds and deploys httpstat-monitor.
	@kubectl apply -f deploy/

.PHONY: bump-version
bump-version: ## Bump the version in the version file. Set SEMVER to [ patch (default) | major | minor ].
	@./scripts/bump-version.sh $(SEMVER)

.PHONY: tag
tag: ## Create and push a new git tag (creates tag using version/version.go file).
	@./scripts/tag.sh

$(GOIMPORTS):
	$(call fetch_go_bin_version,golang.org/x/tools/cmd/goimports,$(GOIMPORTS_VERSION))

$(GOLANGCILINT):
	$(call fetch_go_bin_version,github.com/golangci/golangci-lint/cmd/golangci-lint,$(GOLANGCILINT_VERSION))
