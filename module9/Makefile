
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR := $(dir $(MKFILE_PATH))
ROOT_DIR := $(abspath $(CURRENT_DIR)/..)
BIN_DIR := $(ROOT_DIR)/bin
PROTOC_GEN_GO := $(BIN_DIR)/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(BIN_DIR)/protoc-gen-go-grpc
MOCKERY := $(BIN_DIR)/mockery
BUF := $(BIN_DIR)/buf


.PHONY: proto-gen
proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=$(PROTOC_GEN_GO) \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(PROTOC_GEN_GO_GRPC) \
		proto/todo.proto

.PHONY: run-server
run-server:
	go run cmd/server/main.go

.PHONY: run-client
run-client:
	go run cmd/client/main.go

# Init buf (buf.yaml file)
.PHONY: buf-init
buf-init:
	$(BUF) config init

# Add deps like protovalidate
.PHONY: buf-dep
buf-dep:
	$(BUF) dep update

# Verify buf module (exits with no output)
.PHONY: buf-build
buf-build:
	$(BUF) build

# Generate code
.PHONY: buf-gen
buf-gen:
	$(BUF) generate

.PHONY: run-grpc-server
run-grpc-server:
	go run cmd/grpc-server/main.go

.PHONY: run-grpc-client
run-grpc-client:
	go run cmd/grpc-client/main.go

create-client:
	npx create-next-app client
