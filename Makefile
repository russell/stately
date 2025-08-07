.PHONY: build test test-no-lint clean run deps help

# Default target
help:
	@echo "Available targets:"
	@echo "  build      - Build all targets"
	@echo "  test       - Run all tests"
	@echo "  test-no-lint - Run tests excluding lint checks"
	@echo "  clean      - Clean build artifacts"
	@echo "  run        - Build and run stately"
	@echo "  deps       - Update dependencies (go mod tidy + gazelle)"
	@echo "  help       - Show this help message"

build:
	bazelisk build //...

test:
	bazelisk test //...

test-no-lint:
	bazelisk test //... --test_tag_filters=-lint

clean:
	bazelisk clean

run: build
	bazel-bin/stately

deps:
	go mod tidy
	bazelisk run //:gazelle-update-repos

# For development - watch and rebuild on changes
dev:
	@echo "Use 'ibazel build //...' for continuous builds"
	@echo "Use 'ibazel test //...' for continuous testing"