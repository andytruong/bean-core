BUILD          ?= $(shell git rev-parse --short HEAD)
BUILD_CODENAME  = unnamed
BUILD_DATE     ?= $(shell git log -1 --format=%ci)
BUILD_BRANCH   ?= $(shell git rev-parse --abbrev-ref HEAD)
BUILD_VERSION  ?= $(shell git describe --always --tags)
MODIFIED        = $(shell git diff-index --quiet HEAD || echo "-mod")

version:
	@echo Bean-Core   ${BUILD_VERSION}
	@echo Build:      ${BUILD}
	@echo Codename:   ${BUILD_CODENAME}${MODIFIED}
	@echo Build date: ${BUILD_DATE}
	@echo Branch:     ${BUILD_BRANCH}
	@echo Go version: $(shell go version)

server:
	@CONFIG="config.yaml" go run -mod=vendor cmd/main.go http-server

migrate:
	@CONFIG="config.yaml" go run -mod=vendor cmd/main.go migrate

# ---------------------
# Dev commands
# ---------------------
dev.test:
	@go test -mod=vendor -race -count=1 ./... -v

dev.gql:
	@git checkout pkg/infra/gql/schema.go
	@go run cmd/gqlgen/main.go
	@rm -rf pkg/infra/gql.resolver/*
	@rm pkg/infra/__tmp__resolvers.go
	@go fmt ./...

dev.clean:
	@go mod vendor
	@go fmt ./pkg/...
	@go fmt ./components/...
	@go fmt ./cmd/...
	@go mod tidy
	@git fetch --prune origin
	@rm -rf ./pkg/infra/gql/__tmp__*

dev.lint:
	@go vet ./...
	@staticcheck ./...

dev.check-size:
	@go build -mod=vendor -o /tmp/go-bean cmd/main.go
	@du -h /tmp/go-bean
	@rm /tmp/go-bean

tools.generate.key:
	CONFIG=config.yaml go run cmd/main.go gen-key

tools.generate.ulid:
	@go run cmd/tools/ulid/main.go

# ---------------------
# Command to test on local environment
# ---------------------
local.config:
	export ENV=dev \
	 && export CONFIG=config.yaml \
	 && export CONFIG=config.yaml \
	 && export DB_DRIVER='postgres' \
	 && export DB_MASTER_URL='postgres://postgres:and1bean@127.0.0.1/bean-core?sslmode=disable'

# go run cmd/main.go migrate
local.migrate: dev.config
	env

local.docker.db.init:
	docker run --rm -d --name=hi-pg -p 5432:5432 -e "POSTGRES_PASSWORD=and1bean" postgres:12-alpine

local.docker.db.start:
	docker start hi-pg

local.docker.db.stop:
	@docker stop hi-pg
