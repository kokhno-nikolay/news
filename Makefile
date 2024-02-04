help: ## List all available targets with help
  @grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
    | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build project
  GOOS=linux CGO_ENABLED=0 go build -o=server ./cmd/server/main.go

mock: ## Generate mocks
  go generate ./...