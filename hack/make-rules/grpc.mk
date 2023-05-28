.PHONY: protoc.verify
protoc.verify:
	@echo "===========> verify protoc"
	@if ! which protoc &>/dev/null; then echo "Cannot found protoc compile tool. Please install protoc tool first."; exit 1; fi
	@if ! which protoc-gen-go &>/dev/null; then echo "Cannot found protoc-gen-go. Please install protoc-gen-go tool first."; exit 1; fi
	@if ! which protoc-gen-go-grpc &>/dev/null; then echo "Cannot found protoc-gen-go-grpc tool. Please install protoc-gen-go-grpc tool first."; exit 1; fi
	@if ! which protoc-gen-go-errors &>/dev/null; then echo "Cannot found protoc-gen-go-errors tool. Please install protoc-gen-go-errors tool first."; exit 1; fi

.PHONY: protoc.gen
protoc.gen: protoc.verify tools.verify
	@echo "===========> gen grpc code"
	@protoc --proto_path=. \
        --proto_path=./third_party \
        --go_out=. \
        --go_opt=paths=source_relative \
 		--go-grpc_out=. \
 		--go-grpc_opt=paths=source_relative \
 		--go-errors_out=paths=source_relative:. \
 		--go-http_out=paths=source_relative:. \
 		internal/context/workspace/interface/grpc/proto/*.proto \
 		internal/context/submission/interface/grpc/proto/*.proto \
 		internal/context/notebookserver/interface/grpc/proto/*.proto
