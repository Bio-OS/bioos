DOCKER := docker
IMAGES_DIR ?= $(wildcard ${ROOT_DIR}/build/*)
IMAGES ?= $(foreach image,${IMAGES_DIR},$(notdir ${image}))
IMAGE_REGISTRY ?= quay.io/bioos

.PHONY: image.verify
image.verify:
	@echo "===========> verify docker"
	@if ! which docker &>/dev/null; then echo "Cannot found docker. Please install docker first."; exit 1; fi
	@if ! docker ps &>/dev/null; then echo "Docker daemon not running. Please start it first.";exit 1;fi

.PHONY: buildx.verify
buildx.verify:
	@echo "===========> verify docker buildx"
	@if ! docker buildx version &>/dev/null; then echo "Cannot found docker buildx. Please install docker buildx first. Link: https://docs.docker.com/build/architecture/"; exit 1; fi

.PHONY: image.build
image.build: image.verify  $(addprefix image.build., $(addprefix $(IMAGE_PLAT)., $(IMAGES)))

.PHONY: image.buildx.multiarch.%
image.buildx.multiarch.%: image.verify buildx.verify
	$(eval IMAGE := $(*))
	@echo "===========> Building multiple arch docker image $(IMAGE)"
	$(DOCKER) buildx build \
		--platform linux/amd64,linux/arm64 \
		-t $(IMAGE_REGISTRY)/$(IMAGE):$(VERSION) \
		-f $(ROOT_DIR)/build/$(IMAGE)/Dockerfile . --push

.PHONY: image.buildx
image.buildx: image.verify  $(foreach i,$(IMAGES), $(addprefix image.buildx.multiarch.,$(i)))

.PHONY: image.build.%
image.build.%:
	$(eval IMAGE := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Building docker image $(IMAGE) for $(OS) $(ARCH)"
	DOCKER_BUILDKIT=1 $(DOCKER) build \
		--platform $(OS)/$(ARCH) \
		-t $(IMAGE_REGISTRY)/$(IMAGE):$(VERSION) \
		-f $(ROOT_DIR)/build/$(IMAGE)/Dockerfile .

.PHONY: image.push
image.push: image.build  $(addprefix image.push., $(addprefix $(IMAGE_PLAT)., $(IMAGES)))

.PHONY: image.push.%
image.push.%:
	$(eval IMAGE := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Pushing docker image $(IMAGE) for $(OS) $(ARCH)"
	@echo "$(DOCKER) push $(IMAGE_REGISTRY)/$(IMAGE):$(VERSION)"
	$(DOCKER) push $(IMAGE_REGISTRY)/$(IMAGE):$(VERSION)