.PHONY: proto build run clean

# Генерация protobuf кода
proto:
	@echo "🔧 Generating protobuf code..."
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       proto/auth.proto
	@echo "✅ Protobuf code generated"