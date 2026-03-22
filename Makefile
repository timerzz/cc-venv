APP_NAME := ccv
CMD_PATH := ./cmd/ccv
FRONTEND_DIR := web
FRONTEND_DIST_DIR := $(FRONTEND_DIR)/dist
EMBED_STATIC_DIR := internal/web/static
BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin
NPM_DIR := npm
NPM_MAIN_DIR := $(NPM_DIR)/ccv
NPM_PLATFORM_DIR := $(NPM_DIR)/platforms

.PHONY: frontend linux-amd64 linux-arm64 windows-amd64 windows-arm64 macos-amd64 macos-arm64 linux-x86 windows-x86 macos-x86 all npm-packages clean

frontend:
	@echo "==> Building frontend"
	@if [ ! -d "$(FRONTEND_DIR)/node_modules" ]; then \
		cd "$(FRONTEND_DIR)" && npm install; \
	fi
	cd "$(FRONTEND_DIR)" && npm run build
	@echo "==> Syncing frontend dist to $(EMBED_STATIC_DIR)"
	mkdir -p "$(EMBED_STATIC_DIR)"
	rm -rf "$(EMBED_STATIC_DIR)"/*
	cp -R "$(FRONTEND_DIST_DIR)"/. "$(EMBED_STATIC_DIR)"/

linux-amd64: frontend
	@echo "==> Building linux/amd64"
	mkdir -p "$(BIN_DIR)"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "$(BIN_DIR)/$(APP_NAME)-linux-amd64" $(CMD_PATH)

linux-arm64: frontend
	@echo "==> Building linux/arm64"
	mkdir -p "$(BIN_DIR)"
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o "$(BIN_DIR)/$(APP_NAME)-linux-arm64" $(CMD_PATH)

windows-amd64: frontend
	@echo "==> Building windows/amd64"
	mkdir -p "$(BIN_DIR)"
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o "$(BIN_DIR)/$(APP_NAME)-windows-amd64.exe" $(CMD_PATH)

windows-arm64: frontend
	@echo "==> Building windows/arm64"
	mkdir -p "$(BIN_DIR)"
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o "$(BIN_DIR)/$(APP_NAME)-windows-arm64.exe" $(CMD_PATH)

macos-amd64: frontend
	@echo "==> Building darwin/amd64"
	mkdir -p "$(BIN_DIR)"
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o "$(BIN_DIR)/$(APP_NAME)-darwin-amd64" $(CMD_PATH)

macos-arm64: frontend
	@echo "==> Building darwin/arm64"
	mkdir -p "$(BIN_DIR)"
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o "$(BIN_DIR)/$(APP_NAME)-darwin-arm64" $(CMD_PATH)

# User-facing aliases: interpret "x86" as 64-bit x86 (amd64).
linux-x86: linux-amd64

windows-x86: windows-amd64

macos-x86: macos-amd64

all: linux-amd64 linux-arm64 windows-amd64 windows-arm64 macos-amd64 macos-arm64

npm-packages: all
	@echo "==> Syncing binaries into npm package directories"
	mkdir -p "$(NPM_PLATFORM_DIR)/ccv-linux-x64/bin"
	cp "$(BIN_DIR)/$(APP_NAME)-linux-amd64" "$(NPM_PLATFORM_DIR)/ccv-linux-x64/bin/ccv"
	chmod +x "$(NPM_PLATFORM_DIR)/ccv-linux-x64/bin/ccv"
	mkdir -p "$(NPM_PLATFORM_DIR)/ccv-linux-arm64/bin"
	cp "$(BIN_DIR)/$(APP_NAME)-linux-arm64" "$(NPM_PLATFORM_DIR)/ccv-linux-arm64/bin/ccv"
	chmod +x "$(NPM_PLATFORM_DIR)/ccv-linux-arm64/bin/ccv"
	mkdir -p "$(NPM_PLATFORM_DIR)/ccv-win32-x64/bin"
	cp "$(BIN_DIR)/$(APP_NAME)-windows-amd64.exe" "$(NPM_PLATFORM_DIR)/ccv-win32-x64/bin/ccv.exe"
	mkdir -p "$(NPM_PLATFORM_DIR)/ccv-win32-arm64/bin"
	cp "$(BIN_DIR)/$(APP_NAME)-windows-arm64.exe" "$(NPM_PLATFORM_DIR)/ccv-win32-arm64/bin/ccv.exe"
	mkdir -p "$(NPM_PLATFORM_DIR)/ccv-darwin-x64/bin"
	cp "$(BIN_DIR)/$(APP_NAME)-darwin-amd64" "$(NPM_PLATFORM_DIR)/ccv-darwin-x64/bin/ccv"
	chmod +x "$(NPM_PLATFORM_DIR)/ccv-darwin-x64/bin/ccv"
	mkdir -p "$(NPM_PLATFORM_DIR)/ccv-darwin-arm64/bin"
	cp "$(BIN_DIR)/$(APP_NAME)-darwin-arm64" "$(NPM_PLATFORM_DIR)/ccv-darwin-arm64/bin/ccv"
	chmod +x "$(NPM_PLATFORM_DIR)/ccv-darwin-arm64/bin/ccv"
	@echo "==> npm packages are ready under $(NPM_DIR)/"

clean:
	rm -rf "$(BUILD_DIR)"
