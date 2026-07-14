
all: lint codegen generate test coverage

codegen: clean-mocks
	mockery
.PHONY: codegen

generate:
	@go generate ./cmd/birb
.PHONY: generate

CONTROLLER_GEN ?= $(shell which controller-gen 2>/dev/null || echo $(GOPATH)/bin/controller-gen)

controller-gen:
	@command -v controller-gen >/dev/null 2>&1 || go install sigs.k8s.io/controller-tools/cmd/controller-gen@latest
.PHONY: controller-gen

deepcopy: controller-gen
	$(CONTROLLER_GEN) object paths="./internal/deepcopy/..."
.PHONY: deepcopy

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
	@go test -v ./...
	@internal/failing/validate.sh
.PHONY: test-go

buildDir:
	@mkdir -p build/coverage
.PHONY: buildDir

coverage: buildDir
	@go test -coverpkg=./... -coverprofile=c.out ./...
	@go tool cover -html=c.out -o build/coverage/index.html
	@echo ""
	@echo "Coverage summary:"
	@go tool cover -func=c.out | grep total
.PHONY: coverage

failing:
	@go test -v -tags failing ./internal/failing || true
.PHONY: failing

lint:
	golangci-lint run
.PHONY: lint

lint-fix:
	golangci-lint run --fix
.PHONY: lint-fix
