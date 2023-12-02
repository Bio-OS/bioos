TOOLS ?= swagger addlicense cfssl mockgen gotests git-chglog protoc-gen-go protoc-gen-go-grpc protoc-gen-go-http protoc-gen-go-errors go-gitlint golangci-lint

.PHONY: tools.install
tools.install: $(addprefix tools.install., $(TOOLS))

.PHONY: tools.verify
tools.verify: $(addprefix tools.verify., $(TOOLS))

.PHONY: tools.install.%
tools.install.%:
	@echo "===========> Installing $*"
	@$(MAKE) install.$*

.PHONY: tools.verify.%
tools.verify.%:
	@if ! which $* &>/dev/null; then $(MAKE) tools.install.$*; fi

.PHONY: install.swagger
install.swagger:
	@$(GO) install github.com/swaggo/swag/cmd/swag@v1.8.12

.PHONY: install.cfssl
install.cfssl:
	@$(GO) install github.com/cloudflare/cfssl/cmd/...@latest

.PHONY: install.go-gitlint
install.go-gitlint:
	@$(GO) install github.com/llorllale/go-gitlint@latest

.PHONY: install.git-chglog
install.git-chglog:
	@$(GO) install github.com/git-chglog/git-chglog/cmd/git-chglog@latest

.PHONY: install.mockgen
install.mockgen:
	@$(GO) install github.com/golang/mock/mockgen@latest

.PHONY: install.gotests
install.gotests:
	@$(GO) install github.com/cweill/gotests/gotests@latest

.PHONY: install.addlicense
install.addlicense:
	@$(GO) install github.com/google/addlicense@latest

.PHONY: install.protoc-gen-go
install.protoc-gen-go:
	@$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0

.PHONY: install.protoc-gen-go-grpc
install.protoc-gen-go-grpc:
	@$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

.PHONY: install.protoc-gen-go-http
install.protoc-gen-go-http:
	@$(GO) install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest

.PHONY: install.protoc-gen-go-errors
install.protoc-gen-go-errors:
	@$(GO) install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest

.PHONY: install.golangci-lint
install.golangci-lint:
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: install.womtool
install.womtool:
	@wget --progress=dot -cP ${ROOT_DIR} -O womtool.jar https://github.com/broadinstitute/cromwell/releases/download/85/womtool-85.jar

.PHONY: install.nextflow
install.nextflow:
	@wget --progress=dot -cP ${ROOT_DIR} -O nextflow https://github.com/nextflow-io/nextflow/releases/download/v23.10.0/nextflow-23.10.0-all
