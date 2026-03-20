# Log Analyzer Makefile

BINARY_NAME=log-analyzer
MAIN_FILE=main.go
BUILD_DIR=build

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "✅ Build complete: $(BINARY_NAME)"

# Build for multiple platforms
.PHONY: build-all
build-all: clean
	@echo "🔨 Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)
	@echo "✅ Cross-platform build complete"
	@ls -la $(BUILD_DIR)/

# Run tests
.PHONY: test
test: build
	@echo "🧪 Running tests..."
	@./$(BINARY_NAME) -file=examples/sample.log -verbose
	@./$(BINARY_NAME) -file=examples/sample.log -config=examples/custom-patterns.json
	@echo "✅ Tests completed"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "🧹 Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
	@rm -f *.json
	@echo "✅ Clean complete"

# Install locally
.PHONY: install
install: build
	@echo "📦 Installing $(BINARY_NAME)..."
	@sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "✅ Installed to /usr/local/bin/$(BINARY_NAME)"

# Show help
.PHONY: help
help:
	@echo "Log Analyzer - Makefile Commands:"
	@echo "  make build      - Build the application"
	@echo "  make build-all  - Build for all platforms"
	@echo "  make test       - Run basic tests"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make install    - Install locally"
	@echo "  make help       - Show this help"