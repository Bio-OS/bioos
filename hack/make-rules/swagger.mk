.PHONY: swagger.gen
swagger.gen: install.swagger
	@echo "===========> Generating swag API docs"
	@swag init --parseDependency --dir ./cmd/apiserver,./internal -g apiserver.go
.PHONY: swagger.fmt
swagger.fmt: install.swagger
	@echo "===========> Format swag comments"
	@swag fmt --dir ./cmd/apiserver,./internal

.PHONY: swagger.run
swagger.run: swagger.gen swagger.fmt
