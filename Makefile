.PHONY: setup-grpc compose-up test-auth

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

compose_up:
	@if ! docker ps | grep -q "auth"; then \
		echo "Containers not running, starting them..."; \
		docker compose up -d; \
	else \
		echo "Containers are already running."; \
	fi

check_health: compose_up
	@until grpcurl -plaintext localhost:50051 list > /dev/null 2>&1; do \
		echo "Waiting for gRPC server..."; \
		sleep 2; \
	done
	@echo "gRPC server is ready."

test_auth: check_health
	cd services/auth && go test ./test/... -v