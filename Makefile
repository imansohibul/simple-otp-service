# Load .env file
include .env

.PHONY: clean all init generate generate_mocks

GO_PACKAGES := $(shell go list ./... | grep -v 'cmd\|vendor\|tests\|/generated')
UNAME := $(shell uname)

all: build/main

build/main: cmd/main.go generated
	@echo "Building..."
	go build -o $@ $<

clean:
	rm -rf generated

init: clean generate
	go mod tidy
	go mod vendor

test:
	go clean -testcache
	go test -short -coverprofile coverage.out -short -v ./...
	grep -v -e "mock.gen.go" -e "api_test.go" coverage.out > coverage_clean.out
	go tool cover -func=coverage_clean.out

generate: generated generate_mocks

generated: api.yml install_openapi
	@echo "Generating files..."
	mkdir generated || true
	oapi-codegen --package generated -generate types,server,spec $< > generated/api.gen.go

migrate: install-go-migrate-tool
	@bin/migrate -source file://db/migrate -database "mysql://$(SERVICE_DB_USERNAME):$(SERVICE_DB_PASSWORD)@tcp($(SERVICE_DB_HOST):$(SERVICE_DB_PORT))/$(SERVICE_DB_NAME)?parseTime=true" $(MIGRATE_ARGS) $(N)

create-db-migration: install-go-migrate-tool
	@bin/migrate create -ext sql -dir db/migrate $(MIGRATE_NAME)

generate_mocks: install_mockgen
	go generate ./...

install_openapi:
	@if ! command -v oapi-codegen >/dev/null 2>&1; then \
		echo "Installing oapi-codegen..."; \
		go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest; \
	else \
		echo "oapi-codegen already installed."; \
	fi

install_mockgen:
	@if ! command -v mockgen &> /dev/null; then \
		echo "mockgen not found. Installing..."; \
		go install go.uber.org/mock/mockgen@latest; \
	else \
		echo "mockgen is already installed."; \
	fi

bin:
	@mkdir -p bin

install-go-migrate-tool: bin
ifneq (,$(wildcard bin/migrate))
    # do not download again
else ifeq ($(UNAME), Linux)
	@curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.linux-amd64.tar.gz | tar zxf - --directory /tmp \
	&& cp /tmp/migrate bin/
else ifeq ($(UNAME), Darwin)
	@curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.darwin-amd64.tar.gz | tar zxf - --directory /tmp \
	&& cp /tmp/migrate bin/
else
	@echo $(UNAME)
	@echo "Your OS is not supported."
endif

