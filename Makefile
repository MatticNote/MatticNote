# MatticNote Makefile
CLI_UI_PATH   = ./cli-ui
CLIENT_PATH   = ./client
PACKAGE_NAME  = github.com/MatticNote/MatticNote
BUILD_SUFFIX := $(or $(BUILD_SUFFIX), unknown)

.PHONY: build

build-css:
	npm --prefix ${CLI_UI_PATH} install ${CLI_UI_PATH} --no-bin-links
	npm --prefix ${CLI_UI_PATH} run build:production
build-client: build-css
	npm --prefix ${CLIENT_PATH} install ${CLIENT_PATH} --no-bin-links
	npm --prefix ${CLIENT_PATH} run build
fetch-meta:
	$(eval MN_VERSION=$(or $(shell git describe --tags --abbrev=0), unknown))
	$(eval MN_REVISION=$(shell git rev-parse --short HEAD))
	@echo Version: $(MN_VERSION)-$(MN_REVISION)
build: build-client fetch-meta
	go build \
	-o build/matticnote-$(MN_VERSION)-$(MN_REVISION)-$(BUILD_SUFFIX) \
	-ldflags "-X ${PACKAGE_NAME}/internal.Version=$(MN_VERSION) \
	-X ${PACKAGE_NAME}/internal.Revision=$(MN_REVISION)"
