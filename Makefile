GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOPATH=$(shell go env GOPATH)

.PHONY: all build clean run test lint fmt deps docs

all: build

deps:
	@echo " > Installing dependencies..."
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	go install golang.org/x/tools/cmd/goimports@latest

dev:
	@echo " > Starting development server..."
	go run ./main.go

# docs:
# 	@echo " > Generating Swagger/OpenAPI docs..."
# 	$(GOPATH)/bin/swag init -g cmd/api/main.go --output docs --parseDependency --parseInternal
# 	@echo " > Docs generated at docs/"

format:
	@echo " > Formatting code..."
	gofumpt -w .
	goimports -w -local github.com/zenkiet/edge-gateway .

lint:
	@echo " > Linting..."
	golangci-lint run ./...