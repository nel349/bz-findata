#!make
.SILENT:
.DEFAULT_GOAL := help

help: ## Show this help
	@echo "Usage:\n  make <target>\n"
	@echo "Targets:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Build commands
build-base: ## Build base Docker image
	docker build -t bz-findata-base -f build/docker/base.Dockerfile .

build-all: build-base ## Build all applications
	docker-compose build

build-app: build-base ## Build only main app
	docker-compose build app

build-analysis: build-base ## Build only analysis app
	docker-compose build analysis_app

# Docker commands
docker-compose-up: ## Run all services
	docker-compose up

app-up: ## Run only main app
	docker-compose up --build app mysql

analysis-up: ## Run only analysis app
	docker-compose up --build analysis_app mysql

docker-compose-down: ## Stop all services
	docker-compose down

clean: ## Clean all built images
	docker-compose down --rmi all

# Development commands
deps: ## Download dependencies
	go mod download && go mod tidy

lint: ## Check code (used golangci-lint)
	GO111MODULE=on golangci-lint run

test: ## Run tests
	go clean --testcache
	go test ./...

release: ## Git tag create and push
	git tag -s -a v${tag} -m 'chore(release): v$(tag) [skip ci]'
	git push origin v${tag}

release.revert: ## Revert git release tag
	git tag -d v${tag}
	git push --delete origin v${tag}

run: ## Run application local
	go run cmd/app/main.go