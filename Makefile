export GOFLAGS=-mod=mod
GOLANG_LINTER="./.bin/golangci-lint"
GOENV:=$(shell go env GOPATH)
export PATH:=$(PATH):$(GOENV)/bin

JAVA_BIN:=$(shell which java)
SONAR_VERSION="3.3.0.1492"


.PHONY: githook
githook:
	@echo "Creating pre-commit hook..."
	@rm -f ./.git/hooks/pre-commit
	@cp ./dev-setup/pre-commit.sh ./.git/hooks/pre-commit
.PHONY: setup
## setup: sets up git configurations
setup: gitconfig githook

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: checkfmt
checkfmt:
ifneq ($(shell gofmt -l . | grep -v vendor\/),)
	@gofmt -l . | grep -v vendor\/
	@exit 1
else
	@true
endif

.PHONY: checkgosum
checkgosum:
	@test -s go.sum || { echo "go.sum missing"; exit 1; }

.PHONY: installmockgen
installmockgen:
	@which mockgen >/dev/null 2>&1 || go get github.com/golang/mock/mockgen >/dev/null 2>&1

.PHONY: generate
## generate: auto-generates mock files
generate: installmockgen
	@echo "Generating files..."
	@mkdir mocks 2>/dev/null || true
	@go generate ./... > /dev/null 2>&1

.PHONY: test
## test: runs unit tests
test: generate
	@echo "Running tests..."
	@go test -race -tags=unit ./...

.PHONY: combinecoverage
combinecoverage:
	@{ echo "mode: set"; cat ./coverage/* | grep -v "mode:"; } > full_coverage.out

.PHONY: coverage
## coverage: generates coverage report
coverage: test
	@echo "Collecting coverage..."
	@mkdir ./coverage 2>/dev/null || true
	@go test -coverprofile=coverage/unit.out -tags=unit ./... > /dev/null

.PHONY: update_all_deps
## update_all_deps: update dependencies
update_all_deps:
	@echo "Updating modules..."
	@go get -u ./...
	@go mod tidy

.PHONY: fetch_deps
## fetch_deps: fetch dependencies
fetch_deps:
	@echo "Fetching dependencies..."
	@go get ./...

.PHONY: tidy
## tidy: tidy dependencies
tidy: fetch_deps
	@echo "Tidying modules..."
	@go mod tidy

.PHONY: lint
## lint: lints using golangci
lint:
	@ $(GOLANG_LINTER) run ./...

.PHONY: clean-vendor
clean-vendor:
	@echo "Cleaning up vendors folder..."
	@rm -rf vendor

.PHONY: vendor
vendor:
	@echo "Vendoring..."
	@go mod vendor

.PHONY: docker-prepare
## docker-prepare: prepares to run in local docker
docker-prepare: clean-vendor tidy vendor

.PHONY: mkdirs
mkdirs:
	@mkdir ./.bin 2>/dev/null || true

.PHONY: installlint
installlint: mkdirs
	@[ ! -f "./.bin/golangci-lint" ] && echo "Installing linter..." || true
	@[ ! -f "./.bin/golangci-lint" ] && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./.bin v1.32.2 >/dev/null 2>&1 || true


.PHONY: installsonar
installsonar: mkdirs
ifeq ($(JAVA_BIN),)
	@$(error "Java binary not found. Is JAVA installed?")
endif
	@echo "Installing sonar scanner..."
	@[ ! -f "./.bin/sonar-scanner-$(SONAR_VERSION)" ] && curl -Lso "./.bin/sonar-scanner.zip" "https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-$(SONAR_VERSION).zip"
	@[ ! -f "./.bin/sonar-scanner-$(SONAR_VERSION)" ] && unzip -qo "./.bin/sonar-scanner.zip" -d "./.bin" || true
	@[ ! -f "./.bin/sonar-scanner-$(SONAR_VERSION)" ] && rm -rf ./.bin/sonar-scanner.zip 2>/dev/null || true


.PHONY: install
## install: installs scripts like 'npm i'
install: mkdirs installlint installsonar


.PHONY: sonar
## sonar: reports project to sonar
sonar: installsonar coverage combinecoverage
	@echo "Reporting sonar..."
	@./.bin/sonar-scanner-$(SONAR_VERSION)/bin/sonar-scanner


.PHONY: precommit
precommit: generate checkgosum checkfmt lint test

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
