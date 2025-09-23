.PHONY: all build test clean run

APP_NAME=syntropy
BUILD_DIR=bin
MAIN_GO_FILE=cmd/main.go

all: build

build:
	@echo "Building $(APP_NAME)..."
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_GO_FILE)
	@echo "Build complete. Executable: $(BUILD_DIR)/$(APP_NAME)"

test:
	@echo "Running tests..."
	go test ./...

clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)
	@echo "Cleanup complete."

run: build
	@echo "Running $(APP_NAME)..."
	./$(BUILD_DIR)/$(APP_NAME) $(ARGS)

# Example usage: make run ARGS="usb list"
# Example usage: make run ARGS="usb create --help"

install: build
	@echo "Installing $(APP_NAME)..."
	sudo cp $(BUILD_DIR)/$(APP_NAME) /usr/local/bin/
	@echo "Installation complete."

uninstall:
	@echo "Uninstalling $(APP_NAME)..."
	sudo rm -f /usr/local/bin/$(APP_NAME)
	@echo "Uninstallation complete."

deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies updated."

fmt:
	@echo "Formatting code..."
	go fmt ./...

lint:
	@echo "Running linter..."
	golangci-lint run

vet:
	@echo "Running go vet..."
	go vet ./...
