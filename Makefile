GIT_TAG := $(shell git describe --tags --abbrev=0)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
VERSION := "$(GIT_TAG)-$(GIT_COMMIT)"

WINDOWS_PACKAGE_NAME := "qa-$(GIT_TAG)-windows-x64.zip"
LINUX_PACKAGE_NAME := "qa-$(GIT_TAG)-linux-x64.zip"

default: check_deps build-windows build-linux

check_deps:
	@zip --version > /dev/null
	@go version > /dev/null

go-get: check_deps
	@go get

build-windows: check_deps go-get
	@echo "Building $(VERSION) for Windows/x64"
	@GOOS=windows GOARC=arm64 go build -ldflags="-X main.version=$(APPVERSION)"
	@zip ./artifacts/$(WINDOWS_PACKAGE_NAME) qa.exe > /dev/null
	@echo "Generated $(WINDOWS_PACKAGE_NAME)"

build-linux: check_deps go-get
	@echo "Building $(VERSION) for Linux/x64"
	@GOOS=linux GOARC=arm64 go build -ldflags="-X main.version=$(APPVERSION)"
	@zip ./artifacts/$(LINUX_PACKAGE_NAME) qa > /dev/null
	@echo "Generated $(LINUX_PACKAGE_NAME)"
