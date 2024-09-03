all: mubeng

APP_NAME = mubeng
VERSION  = $(shell git describe --always --tags)
LINT_CMD = "golangci-lint run ./... -v --timeout 5m"
LINT_BIN = "https://install.goreleaser.com/github.com/golangci/golangci-lint.sh"

mubeng: test build

test:
	@echo "Testing ${APP_NAME} package ${VERSION}"
	@go test -short github.com/kitabisa/mubeng/pkg/mubeng
	@go test -short github.com/kitabisa/mubeng/pkg/helper

test-extra: golangci-lint test

build:
	@echo "Building ${APP_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	@mkdir -p bin/
	@go build -ldflags "-s -w -X github.com/kitabisa/mubeng/common.Version=${VERSION}" -o ./bin/${APP_NAME} ./cmd/${APP_NAME}


golangci-lint:
	@echo "Run GolangCI-Lint"
	@if [ -x $(command -v golangci-lint) ]; then \
		eval "${LINT_CMD}"; \
	else \
		echo "Download GolangCI-Lint..."; \
		curl -sfL "${LINT_BIN}" | sh; \
		eval "./bin/${LINT_CMD}"; \
	fi;

clean:
	@echo "Removing binaries"
	@rm -rf bin/