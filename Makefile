#!make
.SILENT:
.DEFAULT_GOAL := help

help: ## Show this help
	@echo "Usage:\n  make <target>\n"
	@echo "Targets:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Build commands
build-base: ## Build base Docker image
	docker build --platform linux/amd64 -t bz-findata-base -f build/docker/base.Dockerfile .

build-all: build-base ## Build all applications
	docker compose build

build-cex-collector: build-base ## Build only main app
	docker build --platform linux/amd64 -t cex-collector -f cmd/cex-collector/Dockerfile .

build-analysis: build-base ## Build only analysis app
	docker build --platform linux/amd64 -t analysis-app -f cmd/analysis/Dockerfile .

build-dex: build-base ## Build only dex app
	docker build --platform linux/amd64 -t dex-app -f cmd/dex/Dockerfile .

build-liquidator: build-base ## Build only liquidator app
	docker build --platform linux/amd64 -t liquidator-app -f cmd/liquidator/Dockerfile .

# Docker commands
docker-compose-up: ## Run all services
	docker compose up

cex-collector-up: build-cex-collector ## Run only main app
	@echo "Running main app and mysql services..."
	docker compose up cex_collector mysql

analysis-up: build-analysis ## Run only analysis app
	@echo "Running analysis app and mysql services..."
	docker compose up analysis_app mysql 

dex-up: build-dex ## Run only dex app
	@echo "Running dex app and mysql services..."
	docker compose up dex_app mysql

liquidator-up: build-liquidator ## Run only liquidator app
	@echo "Running liquidator app"
	docker compose up liquidator_app

docker-compose-down: ## Stop all services
	docker compose down

clean: ## Clean all built images
	docker compose down --rmi all

mysql-up: ## Run mysql service
	docker compose up mysql

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