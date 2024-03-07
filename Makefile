help: ## List all available targets with help
  @grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
    | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build project
  GOOS=linux CGO_ENABLED=0 go build -o=server ./cmd/server/main.go

grpc: ## Generate proto files
	protoc --proto_path api/proto --proto_path vendor.protogen \
	--go_out=api/proto --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=api/proto --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--grpc-gateway_out=api/proto --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
  api/proto/posts.proto

mock: ## Generate mocks
  go generate ./...

swag-format: ## Formatting swagger comments  
  swag fmt - format swag comments

swag: ## Generate swagger documentation
  swag init -g ./cmd/service/main.go

test: ## Runs all tests
  go test -v ./...

local-run: ## Runs docker-compose with postgres, migrations and service 
  docker-compose -f ./deploy/local/docker-compose.yml up