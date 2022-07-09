# MatticNote Makefile
PACKAGE_NAME  = github.com/MatticNote/MatticNote
BUILD_SUFFIX := $(or $(BUILD_SUFFIX), unknown)

.PHONY: build

build-frontend:
	cd client && npm install
	cd client && npm run css:production
	cd client && npm run build -- --no-progress
fetch-meta:
	$(eval MN_VERSION=$(or $(shell git describe --tags --abbrev=0), unknown))
	$(eval MN_REVISION=$(shell git rev-parse --short HEAD))
	$(eval MN_BUILDDATE=$(shell date '+%Y/%m/%d-%H:%M:%S%z'))
	@echo Version: $(MN_VERSION)-$(MN_REVISION)
build: build-frontend fetch-meta
	go build \
	-o build/matticnote-$(MN_VERSION)-$(MN_REVISION)-$(BUILD_SUFFIX) \
	-ldflags "-X ${PACKAGE_NAME}/internal.version=$(MN_VERSION) \
	-X ${PACKAGE_NAME}/internal.revision=$(MN_REVISION) \
	-X ${PACKAGE_NAME}/internal.buildDate=$(MN_BUILDDATE)"
