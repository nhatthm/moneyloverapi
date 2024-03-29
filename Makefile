JSON_CLI_VERSION = 1.8.3
SWAC_VERSION = 0.1.19
GOLANGCI_LINT_VERSION ?= v1.55.2

BIN_DIR = bin
VENDOR_DIR = vendor
OPENAPI = resources/internal/openapi.yaml

GITHUB_OUTPUT ?= /dev/null

GO ?= go
GOLANGCI_LINT ?= $(shell go env GOPATH)/bin/golangci-lint-$(GOLANGCI_LINT_VERSION)
JSON_CLI ?= ${BIN_DIR}/json-cli
SWAC ?= ${BIN_DIR}/swac

.PHONY: $(VENDOR_DIR)
$(VENDOR_DIR):
	@mkdir -p $(VENDOR_DIR)
	@$(GO) mod vendor
	@$(GO) mod tidy

.PHONY: $(GITHUB_OUTPUT)
$(GITHUB_OUTPUT):
	@echo "MODULE_NAME=$(MODULE_NAME)" >> "$@"
	@echo "GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION)" >> "$@"

$(BIN_DIR):
	@mkdir -p $@

$(GOLANGCI_LINT):
	@echo "$(OK_COLOR)==> Installing golangci-lint $(GOLANGCI_LINT_VERSION)$(NO_COLOR)"; \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin "$(GOLANGCI_LINT_VERSION)"
	@mv ./bin/golangci-lint $(GOLANGCI_LINT)

$(JSON_CLI): $(BIN_DIR)
	@curl -s -L 'https://github.com/swaggest/json-cli/releases/download/v$(JSON_CLI_VERSION)/json-cli' > $@
	@chmod +x $@

$(SWAC): $(BIN_DIR)
	@curl -s -L 'https://github.com/swaggest/swac/releases/download/v$(SWAC_VERSION)/swac' > $@
	@chmod +x $@

.PHONY: lint
lint:
	@$(GOLANGCI_LINT) run

.PHONY: test
test: test-unit

## Run unit tests
.PHONY: test-unit
test-unit:
	@echo ">> unit test"
	@$(GO) test -gcflags=-l -coverprofile=unit.coverprofile -covermode=atomic -race ./...

#.PHONY: test-integration
#test-integration:
#	@echo ">> integration test"
#	@$(GO) test ./features/... -gcflags=-l -coverprofile=features.coverprofile -coverpkg ./... -godog -race

.PHONY: generate-category
generate-category: $(JSON_CLI)
	@$(JSON_CLI) gen-go $(OPENAPI) \
		--ptr-in-schema \
			'#/components/schemas/Categories' \
			'#/components/schemas/FlatCategories' \
		--def-ptr '#/components/schemas' \
		--package-name category \
		--output ./pkg/category/entity.generated.go && \
		gofmt -w ./pkg/category/entity.generated.go

.PHONY: generate-transaction
generate-transaction: $(JSON_CLI)
	@$(JSON_CLI) gen-go $(OPENAPI) \
		--patches resources/internal/patch-transaction.json \
		--ptr-in-schema \
			'#/components/schemas/Transactions' \
		--def-ptr '#/components/schemas' \
		--package-name transaction \
		--output ./pkg/transaction/entity.generated.go && \
		gofmt -w ./pkg/transaction/entity.generated.go

.PHONY: generate-user
generate-user: $(JSON_CLI)
	@$(JSON_CLI) gen-go $(OPENAPI) \
		--ptr-in-schema \
			'#/components/schemas/Users' \
		--def-ptr '#/components/schemas' \
		--package-name user \
		--output ./pkg/user/entity.generated.go && \
		gofmt -w ./pkg/user/entity.generated.go

.PHONY: generate-wallet
generate-wallet: $(JSON_CLI)
	@$(JSON_CLI) gen-go $(OPENAPI) \
		--patches resources/internal/patch-wallet.json \
		--ptr-in-schema \
			'#/components/schemas/Wallets' \
		--def-ptr '#/components/schemas' \
		--package-name wallet \
		--output ./pkg/wallet/entity.generated.go && \
		gofmt -w ./pkg/wallet/entity.generated.go

.PHONY: generate-client
generate-client: $(JSON_CLI) $(SWAC)
	@rm -rf ./internal/api && \
		mkdir -p ./internal/api

	@$(SWAC) go-client $(OPENAPI) \
		--patches resources/internal/patch-client.json \
		--operations post/user/login-url,post/token,post/category/list-all,post/transaction/list \
		--skip-default-additional-properties \
		--out ./internal/api \
		--pkg-name api && \
		gofmt -w ./internal/api

.PHONY: generate
generate: generate-client generate-category generate-transaction generate-user generate-wallet
