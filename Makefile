include hack/make-rules/common.mk
include hack/make-rules/copyright.mk
include hack/make-rules/image.mk
include hack/make-rules/golang.mk
include hack/make-rules/grpc.mk
include hack/make-rules/swagger.mk
include hack/make-rules/tools.mk

ROOT_PACKAGE=github.com/Bio-OS/bioos
VERSION        ?= $(shell git describe --tags --always --dirty)
GIT_BRANCH     ?= $(shell git branch | grep \* | cut -d ' ' -f2)
GIT_COMMIT     ?= $(shell git rev-parse HEAD)
GIT_TREE_STATE ?= $(if $(shell git status --porcelain),dirty,clean)
GIT_BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_VERSION     ?= $(shell go version)


.PHONY: build
build: go.build web

.PHONY: image
image: image.build

.PHONY: push
push: image.push

.PHONY: tools
tools:
	@$(MAKE) tools.install

.PHONY: swagger
swagger: swagger.run

## web install
.PHONY: web.install
web.install: 
	npm --prefix=web install

## web dev
.PHONY: web.run
web.run: 
	npm --prefix=web run dev

## web build
.PHONY: web
web:
	@echo "===========> Building web"
	npm --prefix=web run build

.PHONY: clean
clean:
	@echo "===========> Cleaning all build binary"
	@-rm -vrf $(OUTPUT_DIR)
