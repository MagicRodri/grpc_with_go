PROTO_DIR = protos
PROTO_FILE = $(shell find $(PROTO_DIR) -name "*.proto")
OUT_DIR = pkg/generated

dev:
	air .

proto:
	protoc --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) $(PROTO_FILE)

fmt:
	go fmt ./...