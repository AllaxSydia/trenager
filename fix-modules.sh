#!/bin/bash

echo "🔧 Настройка go.mod для всех сервисов..."

# AuthService
cat > backend/AuthService/go.mod << 'MOD'
module auth-service

go 1.25

require (
    github.com/AllaxSydia/trenager/proto v0.0.0-00010101000000-000000000000
    github.com/golang-jwt/jwt/v5 v5.3.1
    github.com/google/uuid v1.6.0
    github.com/lib/pq v1.10.9
    golang.org/x/crypto v0.42.0
    google.golang.org/grpc v1.81.1
    google.golang.org/protobuf v1.36.11
)

replace github.com/AllaxSydia/trenager/proto => ../../proto
MOD

# TaskService
cat > backend/TaskService/go.mod << 'MOD'
module task-service

go 1.25

require (
    github.com/AllaxSydia/trenager/proto v0.0.0-00010101000000-000000000000
    github.com/google/uuid v1.6.0
    github.com/lib/pq v1.10.9
    google.golang.org/grpc v1.81.1
    google.golang.org/protobuf v1.36.11
)

replace github.com/AllaxSydia/trenager/proto => ../../proto
MOD

# GradingService
cat > backend/GradingService/go.mod << 'MOD'
module grading-service

go 1.25

require (
    github.com/AllaxSydia/trenager/proto v0.0.0-00010101000000-000000000000
    github.com/google/uuid v1.6.0
    github.com/lib/pq v1.10.9
    google.golang.org/grpc v1.81.1
    google.golang.org/protobuf v1.36.11
)

replace github.com/AllaxSydia/trenager/proto => ../../proto
MOD

# ExecutionService
cat > backend/ExecutionService/go.mod << 'MOD'
module execution-service

go 1.25

require (
    github.com/AllaxSydia/trenager/proto v0.0.0-00010101000000-000000000000
    github.com/google/uuid v1.6.0
    google.golang.org/grpc v1.81.1
    google.golang.org/protobuf v1.36.11
)

replace github.com/AllaxSydia/trenager/proto => ../../proto
MOD

# AIService
cat > backend/AIService/go.mod << 'MOD'
module ai-service

go 1.25

require (
    github.com/AllaxSydia/trenager/proto v0.0.0-00010101000000-000000000000
    github.com/google/uuid v1.6.0
    google.golang.org/grpc v1.81.1
    google.golang.org/protobuf v1.36.11
)

replace github.com/AllaxSydia/trenager/proto => ../../proto
MOD

# AnalyticsService
cat > backend/AnalyticsService/go.mod << 'MOD'
module analytics-service

go 1.25

require (
    github.com/AllaxSydia/trenager/proto v0.0.0-00010101000000-000000000000
    github.com/google/uuid v1.6.0
    github.com/lib/pq v1.10.9
    google.golang.org/grpc v1.81.1
    google.golang.org/protobuf v1.36.11
)

replace github.com/AllaxSydia/trenager/proto => ../../proto
MOD

echo ""
echo "✅ Все go.mod файлы обновлены!"
echo ""
echo "Теперь выполните:"
echo "  cd backend/AuthService && go mod tidy && go run cmd/main.go"