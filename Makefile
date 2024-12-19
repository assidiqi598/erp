PROTO_AUTH_SRC_DIR=services/auth/proto
PROTO_AUTH_OUT_DIR=services/auth/proto
PROTO_FLUTTER_OUT_DIR=flutter/lib/proto

# Command to generate Go and gRPC code
setup-grpc:
	protoc \
		--proto_path=$(PROTO_AUTH_SRC_DIR) \
		--go_out=$(PROTO_AUTH_SRC_DIR) \
		--go-grpc_out=$(PROTO_AUTH_SRC_DIR) \
		$(PROTO_AUTH_SRC_DIR)/auth.proto && \
	protoc \
		--proto_path=$(PROTO_AUTH_SRC_DIR) \
		--dart_out=grpc:$(PROTO_FLUTTER_OUT_DIR) \
		$(PROTO_AUTH_SRC_DIR)/auth.proto
