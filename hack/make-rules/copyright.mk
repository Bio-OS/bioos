.PHONY: gen.add-copyright
gen.add-copyright:
	@addlicense -v -f $(ROOT_DIR)/hack/boilerplate/boilerplate.txt  \
		--ignore "build/**"              \
		--ignore "conf/**"               \
		--ignore "docker-compose.yaml"   \
		--ignore "docs/**"               \
		--ignore "internal/**"           \
		--ignore "third_party/**"        \
		--ignore "vendor/**"             \
		--ignore "web/node_modules/**"   \
		--ignore "web/*.js"              \
		--ignore "web/**/nbv.js"         \
		--ignore "web/**/prism.js"       \
		--ignore ".golangci.yaml"        \
		--ignore ".idea/**"              \
		--ignore ".git/**"  .