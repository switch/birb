
all: lint codegen generate test coverage

codegen: clean-mocks
	mockery
.PHONY: codegen

generate:
	@go generate ./cmd/birb
.PHONY: generate

clean: clean-mocks clean-buildDir
.PHONY: clean

clean-mocks:
	@rm -f test/e2e/mocks_test.go
.PHONY: clean-mocks

clean-buildDir:
	@rm -rf build
.PHONY: clean-buildDir

test: test-go
.PHONY: test

test-go:
	@go test -v $$( go list ./... | grep -v /failing )
	@internal/failing/validate.sh
.PHONY: test-go

buildDir:
	@mkdir -p build/coverage
.PHONY: buildDir

coverage: buildDir
	@go test ./... --cover --coverprofile=c.out $$(go list ./...)
	@go tool cover -html=c.out -o build/coverage/index.html
.PHONY: coverage

lint:
	golangci-lint run
.PHONY: lint

lint-fix:
	golangci-lint run --fix
.PHONY: lint-fix
