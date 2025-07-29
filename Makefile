# Makefile pour Zelda Souls Game

# Variables
BINARY_NAME=zelda-souls-game
MAIN_PATH=./cmd/game
BUILD_DIR=./build
ASSETS_DIR=./assets
CONFIG_DIR=./configs

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Build flags
BUILD_FLAGS=-ldflags "-X main.Version=$(shell git describe --tags --always --dirty)"
DEBUG_FLAGS=-gcflags="all=-N -l"

# Platform specific
ifeq ($(OS),Windows_NT)
    BINARY_NAME := $(BINARY_NAME).exe
    RM = del /Q
    MKDIR = mkdir
else
    RM = rm -f
    MKDIR = mkdir -p
endif

.PHONY: all build clean test deps run dev release help setup assets

# Default target
all: deps test build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(MKDIR) $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build with debug symbols
debug:
	@echo "Building $(BINARY_NAME) with debug symbols..."
	$(MKDIR) $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) $(DEBUG_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Debug build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for release (optimized)
release:
	@echo "Building $(BINARY_NAME) for release..."
	$(MKDIR) $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -ldflags "-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Release build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	cd $(BUILD_DIR) && ./$(BINARY_NAME)

# Run in development mode (with auto-rebuild)
dev:
	@echo "Running in development mode..."
	$(GOBUILD) $(DEBUG_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	cd $(BUILD_DIR) && ./$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	$(RM) $(BUILD_DIR)/$(BINARY_NAME)
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	golangci-lint run

# Setup project structure and default files
setup:
	@echo "Setting up project structure..."
	$(MKDIR) $(BUILD_DIR)
	$(MKDIR) $(ASSETS_DIR)/textures/player
	$(MKDIR) $(ASSETS_DIR)/textures/enemies
	$(MKDIR) $(ASSETS_DIR)/textures/items
	$(MKDIR) $(ASSETS_DIR)/textures/ui
	$(MKDIR) $(ASSETS_DIR)/textures/environment
	$(MKDIR) $(ASSETS_DIR)/sounds/sfx
	$(MKDIR) $(ASSETS_DIR)/sounds/music
	$(MKDIR) $(ASSETS_DIR)/maps
	$(MKDIR) $(ASSETS_DIR)/tilesets
	$(MKDIR) $(ASSETS_DIR)/data
	$(MKDIR) $(CONFIG_DIR)
	$(MKDIR) ./saves
	$(MKDIR) ./logs
	$(MKDIR) ./screenshots
	@echo "Creating default config file..."
	@echo "# Configuration par défaut pour Zelda Souls Game" > $(CONFIG_DIR)/game_config.yaml
	@echo "window:" >> $(CONFIG_DIR)/game_config.yaml
	@echo "  width: 1280" >> $(CONFIG_DIR)/game_config.yaml
	@echo "  height: 720" >> $(CONFIG_DIR)/game_config.yaml
	@echo "  title: \"Zelda Souls Game\"" >> $(CONFIG_DIR)/game_config.yaml
	@echo "  fullscreen: false" >> $(CONFIG_DIR)/game_config.yaml
	@echo "  vsync: true" >> $(CONFIG_DIR)/game_config.yaml
	@echo "" >> $(CONFIG_DIR)/game_config.yaml
	@echo "rendering:" >> $(CONFIG_DIR)/game_config.yaml
	@echo "  target_fps: 60" >> $(CONFIG_DIR)/game_config.yaml
	@echo "  tile_size: 32" >> $(CONFIG_DIR)/game_config.yaml
	@echo "  enable_batching: true" >> $(CONFIG_DIR)/game_config.yaml
	@echo "  enable_culling: true" >> $(CONFIG_DIR)/game_config.yaml
	@echo "" >> $(CONFIG_DIR)/game_config.yaml
	@echo "audio:" >> $(CONFIG_DIR)/game_config.yaml
	@echo "  master_volume: 1.0" >> $(CONFIG_DIR)/game_config.yaml
	@echo "  music_volume: 0.7" >> $(CONFIG_DIR)/game_config.yaml
	@echo "  sfx_volume: 0.8" >> $(CONFIG_DIR)/game_config.yaml
	@echo "" >> $(CONFIG_DIR)/game_config.yaml
	@echo "debug:" >> $(CONFIG_DIR)/game_config.yaml
	@echo "  enable_debug: false" >> $(CONFIG_DIR)/game_config.yaml
	@echo "  show_fps: false" >> $(CONFIG_DIR)/game_config.yaml
	@echo "  show_colliders: false" >> $(CONFIG_DIR)/game_config.yaml
	@echo "Project structure created successfully!"

# Create example assets structure
assets:
	@echo "Creating example assets..."
	@echo "# Place your player sprites here" > $(ASSETS_DIR)/textures/player/README.md
	@echo "# Place your enemy sprites here" > $(ASSETS_DIR)/textures/enemies/README.md
	@echo "# Place your item sprites here" > $(ASSETS_DIR)/textures/items/README.md
	@echo "# Place your UI sprites here" > $(ASSETS_DIR)/textures/ui/README.md
	@echo "# Place your environment sprites here" > $(ASSETS_DIR)/textures/environment/README.md
	@echo "# Place your sound effects here (.wav, .ogg)" > $(ASSETS_DIR)/sounds/sfx/README.md
	@echo "# Place your music files here (.mp3, .ogg)" > $(ASSETS_DIR)/sounds/music/README.md
	@echo "# Place your Tiled maps here (.tmx)" > $(ASSETS_DIR)/maps/README.md
	@echo "# Place your Tiled tilesets here (.tsx)" > $(ASSETS_DIR)/tilesets/README.md
	@echo "Assets structure created!"

# Build for different platforms
build-windows:
	@echo "Building for Windows..."
	$(MKDIR) $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME:.exe=)-windows.exe $(MAIN_PATH)

build-linux:
	@echo "Building for Linux..."
	$(MKDIR) $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)

build-mac:
	@echo "Building for macOS..."
	$(MKDIR) $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-macos $(MAIN_PATH)

# Build for all platforms
build-all: build-windows build-linux build-mac

# Create a release package
package: release
	@echo "Creating release package..."
	$(MKDIR) $(BUILD_DIR)/release
	cp $(BUILD_DIR)/$(BINARY_NAME) $(BUILD_DIR)/release/
	cp -r $(ASSETS_DIR) $(BUILD_DIR)/release/
	cp -r $(CONFIG_DIR) $(BUILD_DIR)/release/
	@echo "Release package created in $(BUILD_DIR)/release/"

# Install the game (copy to system)
install: build
	@echo "Installing $(BINARY_NAME)..."
ifeq ($(OS),Windows_NT)
	@echo "Manual installation required on Windows"
else
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "$(BINARY_NAME) installed to /usr/local/bin/"
endif

# Uninstall the game
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
ifeq ($(OS),Windows_NT)
	@echo "Manual uninstallation required on Windows"
else
	sudo $(RM) /usr/local/bin/$(BINARY_NAME)
	@echo "$(BINARY_NAME) uninstalled"
endif

# Generate documentation
docs:
	@echo "Generating documentation..."
	$(GOCMD) doc -all ./... > docs/api.md
	@echo "Documentation generated in docs/api.md"

# Benchmarks
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Profile the application
profile: build
	@echo "Running with profiling..."
	cd $(BUILD_DIR) && ./$(BINARY_NAME) -cpuprofile=cpu.prof -memprofile=mem.prof

# Check for security issues (requires gosec)
security:
	@echo "Checking for security issues..."
	gosec ./...

# Check for dependency vulnerabilities (requires nancy)
vuln-check:
	@echo "Checking for vulnerable dependencies..."
	$(GOMOD) list -json -m all | nancy sleuth

# Full check (format, lint, test, security)
check: fmt lint test security

# Help
help:
	@echo "Available targets:"
	@echo "  all          - Download deps, run tests, and build"
	@echo "  build        - Build the application"
	@echo "  debug        - Build with debug symbols"
	@echo "  release      - Build optimized for release"
	@echo "  run          - Build and run the application"
	@echo "  dev          - Run in development mode"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  deps         - Install dependencies"
	@echo "  deps-update  - Update dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  setup        - Setup project structure"
	@echo "  assets       - Create example assets structure"
	@echo "  build-all    - Build for all platforms"
	@echo "  package      - Create release package"
	@echo "  install      - Install the game system-wide"
	@echo "  uninstall    - Uninstall the game"
	@echo "  docs         - Generate documentation"
	@echo "  bench        - Run benchmarks"
	@echo "  profile      - Run with profiling enabled"
	@echo "  security     - Check for security issues"
	@echo "  vuln-check   - Check for vulnerable dependencies"
	@echo "  check        - Run all checks (format, lint, test, security)"
	@echo "  help         - Show this help message"