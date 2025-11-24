# Makefile for g0 load tester

# Binary name
BINARY_NAME=g0

# Build directory
BUILD_DIR=.
PKG_DIR=dist/pkg
DMG_DIR=dist/dmg

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Package info
PKG_NAME=g0
PKG_VERSION=1.0.0
PKG_ID=com.g0.loadtester
PKG_OUTPUT=dist/$(PKG_NAME)-$(PKG_VERSION).pkg
DMG_OUTPUT=dist/$(PKG_NAME)-$(PKG_VERSION).dmg

# Package metadata
PKG_TITLE=g0 Load Tester
PKG_DESCRIPTION=High-performance HTTP load testing tool
PKG_ORGANIZATION=g0
PKG_REPO_URL=https://github.com/calummacc/g0
PKG_AUTHOR=Calumma Team

.PHONY: all build clean run test help pkg dmg install-pkg

# Default target
all: clean build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "Build complete!"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@if [ -f $(BUILD_DIR)/$(BINARY_NAME) ]; then \
		rm -f $(BUILD_DIR)/$(BINARY_NAME); \
		echo "Removed $(BINARY_NAME)"; \
	else \
		echo "No binary to clean"; \
	fi

# Run the load test
run:
	@echo "Running load test..."
	./$(BINARY_NAME) run --url https://httpbin.org/get --concurrency 50 --duration 10s

# Clean, build, and run test
test: clean build run
	@echo "Test complete!"

# Build macOS installer package (.pkg)
pkg: build
	@echo "Building macOS installer package..."
	@mkdir -p $(PKG_DIR)/usr/local/bin
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(PKG_DIR)/usr/local/bin/
	@chmod +x $(PKG_DIR)/usr/local/bin/$(BINARY_NAME)
	@mkdir -p dist
	@echo "Creating component package..."
	@pkgbuild --root $(PKG_DIR) \
		--identifier $(PKG_ID) \
		--version $(PKG_VERSION) \
		--install-location / \
		--ownership recommended \
		dist/component.pkg
	@echo "Creating distribution package with metadata..."
	@echo '<?xml version="1.0" encoding="utf-8"?>' > dist/Distribution.xml
	@echo '<installer-gui-script minSpecVersion="1">' >> dist/Distribution.xml
	@echo '    <title>$(PKG_TITLE)</title>' >> dist/Distribution.xml
	@echo '    <organization>$(PKG_ID)</organization>' >> dist/Distribution.xml
	@echo '    <domains enable_localSystem="true"/>' >> dist/Distribution.xml
	@echo '    <options customize="never" require-scripts="false" rootVolumeOnly="true" />' >> dist/Distribution.xml
	@echo '    <pkg-ref id="$(PKG_ID)">' >> dist/Distribution.xml
	@echo '        <bundle-version>' >> dist/Distribution.xml
	@echo '            <string>$(PKG_VERSION)</string>' >> dist/Distribution.xml
	@echo '        </bundle-version>' >> dist/Distribution.xml
	@echo '    </pkg-ref>' >> dist/Distribution.xml
	@echo '    <choices-outline>' >> dist/Distribution.xml
	@echo '        <line choice="default">' >> dist/Distribution.xml
	@echo '            <line choice="$(PKG_ID)"/>' >> dist/Distribution.xml
	@echo '        </line>' >> dist/Distribution.xml
	@echo '    </choices-outline>' >> dist/Distribution.xml
	@echo '    <choice id="default"/>' >> dist/Distribution.xml
	@echo '    <choice id="$(PKG_ID)" visible="false">' >> dist/Distribution.xml
	@echo '        <pkg-ref id="$(PKG_ID)"/>' >> dist/Distribution.xml
	@echo '    </choice>' >> dist/Distribution.xml
	@echo '    <pkg-ref id="$(PKG_ID)" version="$(PKG_VERSION)" onConclusion="none">component.pkg</pkg-ref>' >> dist/Distribution.xml
	@echo '</installer-gui-script>' >> dist/Distribution.xml
	@mkdir -p dist/Resources
	@echo "Creating package info file..."
	@echo "Package: $(PKG_TITLE)" > dist/Resources/README.txt
	@echo "Version: $(PKG_VERSION)" >> dist/Resources/README.txt
	@echo "Organization: $(PKG_ORGANIZATION)" >> dist/Resources/README.txt
	@echo "Repository: $(PKG_REPO_URL)" >> dist/Resources/README.txt
	@echo "Author: $(PKG_AUTHOR)" >> dist/Resources/README.txt
	@echo "" >> dist/Resources/README.txt
	@echo "Description: $(PKG_DESCRIPTION)" >> dist/Resources/README.txt
	@productbuild --distribution dist/Distribution.xml \
		--package-path dist \
		--resources dist/Resources \
		$(PKG_OUTPUT)
	@echo "Package created: $(PKG_OUTPUT)"
	@echo "  Title: $(PKG_TITLE)"
	@echo "  Version: $(PKG_VERSION)"
	@echo "  Organization: $(PKG_ORGANIZATION)"
	@echo "  Repository: $(PKG_REPO_URL)"
	@rm -rf $(PKG_DIR) dist/component.pkg dist/Distribution.xml dist/Resources

# Build macOS disk image (.dmg)
dmg: pkg
	@echo "Building macOS disk image..."
	@mkdir -p $(DMG_DIR)
	@cp $(PKG_OUTPUT) $(DMG_DIR)/
	@mkdir -p dist
	@hdiutil create -volname "$(PKG_NAME) $(PKG_VERSION)" \
		-srcfolder $(DMG_DIR) \
		-ov -format UDZO \
		$(DMG_OUTPUT)
	@echo "Disk image created: $(DMG_OUTPUT)"
	@rm -rf $(DMG_DIR)

# Install package (requires sudo)
install-pkg: pkg
	@echo "Installing package (requires sudo)..."
	@sudo installer -pkg $(PKG_OUTPUT) -target /
	@echo "Installation complete! You can now run 'g0' from anywhere."

# Clean package artifacts
clean-pkg:
	@echo "Cleaning package artifacts..."
	@rm -rf dist $(PKG_DIR) $(DMG_DIR)
	@echo "Package artifacts cleaned!"

# Help target
help:
	@echo "Available targets:"
	@echo "  make build      - Build the application"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make run        - Run the load test"
	@echo "  make test       - Clean, build, and run test"
	@echo "  make all        - Clean and build"
	@echo "  make pkg        - Build macOS installer package (.pkg)"
	@echo "  make dmg        - Build macOS disk image (.dmg)"
	@echo "  make install-pkg - Install the package (requires sudo)"
	@echo "  make clean-pkg  - Clean package artifacts"
	@echo "  make help       - Show this help message"

