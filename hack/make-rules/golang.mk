GO := go
GO_SUPPORTED_VERSIONS ?= 1.13|1.14|1.15|1.16|1.17|1.18|1.19|1.20
GO_LDFLAGS += -X $(ROOT_PACKAGE)/pkg/version.Version=$(VERSION) \
	-X $(ROOT_PACKAGE)/pkg/version.Branch=$(GIT_BRANCH) \
	-X $(ROOT_PACKAGE)/pkg/version.GitCommit=$(GIT_COMMIT) \
	-X $(ROOT_PACKAGE)/pkg/version.GitTreeState=$(GIT_TREE_STATE) \
	-X $(ROOT_PACKAGE)/pkg/version.GitBranch=$(GIT_BRANCH) \
	-X $(ROOT_PACKAGE)/pkg/version.BuildTime=$(shell date +'%Y-%m-%dT%H:%M:%SZ')
ifneq ($(DLV),)
	GO_BUILD_FLAGS += -gcflags "all=-N -l"
	LDFLAGS = ""
endif
GO_BUILD_FLAGS += -ldflags "$(GO_LDFLAGS)"

COMMANDS ?= $(wildcard ${ROOT_DIR}/cmd/*)
BINS ?= $(foreach cmd,${COMMANDS},$(notdir ${cmd}))


.PHONY: go.tidy
go.tidy:
	@$(GO) mod tidy

.PHONY: go.build.verify
go.build.verify:
	@echo "===========> verify go"
	@if ! which go &>/dev/null; then echo "Cannot found go compile tool. Please install go tool first."; exit 1; fi


.PHONY: go.build
go.build: go.build.verify go.tidy $(addprefix go.build., $(addprefix $(PLATFORM)., $(BINS)))


.PHONY: go.build.%
go.build.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Building binary $(COMMAND) $(VERSION) for $(OS) $(ARCH)"
	@echo "GOOS=$(OS) GOARCH=$(ARCH) $(GO) build $(GO_BUILD_FLAGS) -o $(OUTPUT_DIR)/platforms/$(OS)/$(ARCH)/$(COMMAND)$(GO_OUT_EXT) $(ROOT_PACKAGE)/cmd/$(COMMAND)"
	@GOOS=$(OS) GOARCH=$(ARCH) $(GO) build $(GO_BUILD_FLAGS) -o $(OUTPUT_DIR)/platforms/$(OS)/$(ARCH)/$(COMMAND)$(GO_OUT_EXT) $(ROOT_PACKAGE)/cmd/$(COMMAND)

.PHONY: go.lint
go.lint: tools.verify.golangci-lint
	@echo "===========> Run golangci-lint"
	@golangci-lint run -c $(ROOT_DIR)/.golangci.yaml $(ROOT_DIR)/...

.PHONY: go.run
go.run: tools.install.womtool generate.certs $(addprefix go.run., $(addprefix $(PLATFORM)., apiserver))

.PHONY: generate.certs
generate.certs:
	@echo "===========> Generating certs"
	@bash ${ROOT_DIR}/conf/certs/certs.sh

.PHONY: go.run.%
go.run.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Running binary $(COMMAND) $(VERSION) for $(OS) $(ARCH)"
	@$(OUTPUT_DIR)/platforms/$(OS)/$(ARCH)/$(COMMAND)$(GO_OUT_EXT) --config conf/apiserver.yaml
