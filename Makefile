PROTO_DIR = internal/proto
PROTO_FILE = $(PROTO_DIR)/*.proto
OUT_DIR = pkg

dev:
	air .

proto:
	protoc --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) $(PROTO_FILE)

fmt:
	go fmt ./...