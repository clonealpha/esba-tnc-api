#!/bin/bash

# VPP API에서 binapi 생성 스크립트
# Usage: ./scripts/generate-binapi.sh [VPP_API_DIR]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
GOVPP_DIR="${PROJECT_DIR}/govpp"
BINAPI_DIR="${GOVPP_DIR}/binapi"

# VPP API 디렉토리 (기본값: /usr/share/vpp/api)
VPP_API_DIR="${1:-/usr/share/vpp/api}"

if [ ! -d "$VPP_API_DIR" ]; then
    echo "오류: VPP API 디렉토리를 찾을 수 없습니다: $VPP_API_DIR"
    echo "사용법: $0 [VPP_API_DIR]"
    exit 1
fi

echo "VPP API에서 binapi 생성 중..."
echo "  VPP API 디렉토리: $VPP_API_DIR"
echo "  출력 디렉토리: $BINAPI_DIR"
echo ""

# binapi-generator 빌드
echo "binapi-generator 빌드 중..."
cd "$GOVPP_DIR"
if ! go build -o /tmp/binapi-generator ./cmd/binapi-generator; then
    echo "오류: binapi-generator 빌드 실패"
    exit 1
fi

# 기존 binapi 디렉토리 백업 (선택사항)
if [ -d "$BINAPI_DIR" ] && [ "$(ls -A $BINAPI_DIR)" ]; then
    echo "기존 binapi 디렉토리 발견, 백업 생성 중..."
    BACKUP_DIR="${BINAPI_DIR}.backup.$(date +%Y%m%d_%H%M%S)"
    mv "$BINAPI_DIR" "$BACKUP_DIR"
    echo "백업 위치: $BACKUP_DIR"
fi

# binapi 생성
echo "binapi 생성 중..."
mkdir -p "$BINAPI_DIR"

# binapi-generator 실행
/tmp/binapi-generator \
    --input "$VPP_API_DIR" \
    --output "$BINAPI_DIR" \
    --import-prefix "github.com/clonealpha/esba-tnc-api/govpp"

if [ $? -eq 0 ]; then
    echo "binapi 생성 완료!"
    echo "생성된 파일: $BINAPI_DIR"
else
    echo "오류: binapi 생성 실패"
    exit 1
fi

# 임시 파일 정리
rm -f /tmp/binapi-generator

echo ""
echo "다음 단계:"
echo "  1. proto 생성: ./scripts/generate-proto.sh"
echo "  2. proto 컴파일: ./scripts/generate-proto.sh --compile"

