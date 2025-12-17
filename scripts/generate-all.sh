#!/bin/bash

# 전체 빌드 스크립트
# binapi 생성 → proto 생성 → proto 컴파일 → 변환 함수 생성

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_DIR"

echo "=== ESBA-TNC-API 전체 빌드 ==="
echo ""

# 1. binapi 생성 (선택사항)
if [ "$1" == "--with-binapi" ]; then
    echo "1. binapi 생성 중..."
    if [ -n "$2" ]; then
        ./scripts/generate-binapi.sh "$2"
    else
        ./scripts/generate-binapi.sh
    fi
    echo ""
fi

# 2. proto 생성
echo "2. proto 생성 중..."
cd tools/proto-generator
go run . \
    --binapi-dir="$PROJECT_DIR/govpp/binapi" \
    --output="$PROJECT_DIR/proto" \
    --config="config/proto.yaml" \
    --proto
cd "$PROJECT_DIR"
echo ""

# 3. proto 컴파일
echo "3. proto 컴파일 중..."
./scripts/generate-proto.sh
echo ""

# 4. 변환 함수 생성 (선택사항)
if [ "$1" == "--with-converters" ] || [ "$1" == "--with-binapi" ]; then
    echo "4. 변환 함수 생성 중..."
    cd tools/proto-generator
    go run . \
        --binapi-dir="$PROJECT_DIR/govpp/binapi" \
        --config="config/proto.yaml" \
        --converters
    cd "$PROJECT_DIR"
    echo ""
fi

echo "=== 빌드 완료 ==="
echo ""
echo "생성된 파일:"
echo "  - proto/agent_generated.proto"
if [ "$1" == "--with-converters" ] || [ "$1" == "--with-binapi" ]; then
    echo "  - ../esba-tnc-agent/agent/grpc/handler/converters_gen.go"
fi

