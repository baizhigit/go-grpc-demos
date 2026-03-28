
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR := $(dir $(MKFILE_PATH))
ROOT_DIR := $(abspath $(CURRENT_DIR)/..)
BIN_DIR := $(ROOT_DIR)/bin
PROTOC_GEN_GO := $(BIN_DIR)/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(BIN_DIR)/protoc-gen-go-grpc


.PHONY: proto-gen
proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=$(PROTOC_GEN_GO) \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(PROTOC_GEN_GO_GRPC) \
		proto/config.proto

.PHONY: run-server
run-server:
	go run cmd/server/main.go

.PHONY: run-client
run-client:
	go run cmd/client/main.go

.PHONY: run-servers
run-servers:
	PORT=50051 go run cmd/server/main.go &
	PORT=50052 go run cmd/server/main.go &
	PORT=50053 go run cmd/server/main.go
