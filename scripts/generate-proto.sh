#!/bin/bash

echo "🚀 Generating protobuf code..."

# Проверяем наличие protoc
if ! command -v protoc &> /dev/null; then
    echo "❌ protoc not found. Please install protobuf compiler"
    echo "  macOS: brew install protobuf"
    echo "  Linux: apt install protobuf-compiler"
    echo "  Windows: choco install protoc"
    exit 1
fi

# Генерируем для каждого сервиса
for service in auth task execution grading analytics ai; do
    if [ -f "proto/$service/$service.proto" ]; then
        echo "  📝 Generating $service..."
        protoc \
            --proto_path=proto \
            --go_out=. \
            --go-grpc_out=. \
            "proto/$service/$service.proto"
        
        if [ $? -eq 0 ]; then
            echo "    ✅ $service generated"
        else
            echo "    ❌ Failed to generate $service"
        fi
    fi
done

echo "✨ Done!"