
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR := $(dir $(MKFILE_PATH))
ROOT_DIR := $(abspath $(CURRENT_DIR)/..)
BIN_DIR := $(ROOT_DIR)/bin
PROTOC_GEN_GO := $(BIN_DIR)/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(BIN_DIR)/protoc-gen-go-grpc
MOCKERY := $(BIN_DIR)/mockery


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


.PHONY: mockery-init
mockery-init:
	$(MOCKERY) init github.com/baizhigit/go-grpc-demos/module7/internal/todo

.PHONY: mockery-run
mockery-run:
	$(MOCKERY)

.PHONY: test-unit
test-unit:
	go test -v ./...