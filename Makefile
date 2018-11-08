PACKAGE  = goreactapp
DATE     = $(shell date +%s)
BIN      = $(GOPATH)/bin
BASE     = $(GOPATH)/src/github.com/vshiva/goreactapp
PKGS     = $(or $(PKG),$(shell cd $(BASE) && env GOPATH=$(GOPATH) $(GO) list ./... | grep -v "^$(PACKAGE)/vendor/"))
TESTPKGS = $(shell env GOPATH=$(GOPATH) $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))

COMPILED      ?= $(shell date +%FT%T%z)
PATCH_VERSION ?= $(shell if [ -z "$(BUILD_TAG)" ] && [ "$(BUILD_TAG)xxx" == "xxx" ]; then echo "dev"; else echo $(BUILD_TAG); fi)
GIT_COMMIT    ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo unknown)

GO       = go
GODOC    = godoc
GOFMT    = gofmt
DOCKER   = docker
TIMEOUT  = 15

LD_FLAGS = "-X github.com/vshiva/goreactapp.PatchVersion=$(PATCH_VERSION) -X github.com/vshiva/goreactapp.Compiled=$(DATE) -X github.com/vshiva/goreactapp.GitCommit=$(GIT_COMMIT)"
DOCKER_BUILD_ARGS = '--build-arg LD_FLAGS=${LD_FLAGS}'

V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

.PHONY: all
all: fmt lint build

# Tools

GOLINT = $(BIN)/golint
$(BIN)/golint: | $(BASE) ; $(info $(M) installing golint…)
	$Q $(GO) get -u golang.org/x/lint/golint

GODEP = $(BIN)/dep
$(BIN)/dep: | $(BASE) ; $(info $(M) installing dep…)
	$Q $(GO) get -u github.com/golang/dep/cmd/dep

PROTOC_GO = $(BIN)/protoc-gen-go
$(BIN)/protoc-gen-go: | $(BASE) ; $(info $(M) installing protoc-gen-go…)
	$Q $(GO) get -u github.com/golang/protobuf/protoc-gen-go

.PHONY: install-dependencies 
install-dependencies: $(BIN)/golint $(BIN)/dep $(BIN)/protoc-gen-go   ## Install dependent go tools

.PHONY: lint
lint: vendor | $(BASE) $(GOLINT) ; $(info $(M) running golint…) @ ## Run golint
	$Q cd $(BASE) && ret=0 && for pkg in $(PKGS); do \
		test -z "$$($(GOLINT) $$pkg | tee /dev/stderr)" || ret=1 ; \
	 done ; exit $$ret

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	@ret=0 && for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOFMT) -l -w $$d/*.go || ret=$$? ; \
	 done ; exit $$ret

# Dependency management

.PHONY: vendor
vendor: 
	$Q cd $(BASE) && $(GODEP) ensure -v; $(info $(M) retrieving dependencies…)

.PHONY: gen
gen: 
	$Q cd $(BASE)/web && yarn build && cd $(BASE) && $(GO) generate ./web; $(info $(M) generating web assets)

.PHONY: build
build: gen ## Build service
	$Q cd $(BASE) && $(GO) build -ldflags $(LD_FLAGS) \
		-o bin/$(PACKAGE) ./cmd ; $(info $(M) building executable…)

.PHONY: docker-build
build-image: ## Build docker image
	$Q cd $(BASE) && $(DOCKER) build '$(DOCKER_BUILD_ARGS)' -t goreactapp:$(PATCH_VERSION) . && \
		$(DOCKER) tag goreactapp:$(PATCH_VERSION) vshiva/goreactapp:$(PATCH_VERSION); $(info $(M) building docker image…)

.PHONY: docker-push
push-image: docker-build ## Build and publish docker image
	$Q $(DOCKER) push vshiva/goreactapp:$(PATCH_VERSION); $(info $(M) pushing docker image…)

# Misc

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf bin
	@rm -rf test/tests.* test/coverage.*

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'