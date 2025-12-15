# Protocol Buffers 생성 파일 관리

## 현재 상태

생성된 파일 (`*.pb.go`)이 **GitHub 저장소에 포함**되어 있습니다.

## 이유

1. **사용 편의성**: 다른 프로젝트에서 모듈을 사용할 때 별도로 컴파일할 필요가 없습니다.
2. **일관성**: 모든 사용자가 동일한 생성된 코드를 사용합니다.
3. **빌드 단순화**: CI/CD에서 추가 빌드 단계가 필요 없습니다.

## .gitignore 설정

원래 `.gitignore`에 `*.pb.go`가 있었지만, 공유 모듈의 특성상 생성된 파일을 포함하는 것이 더 유용하므로 제외 목록에서 제거했습니다.

## 업데이트 프로세스

proto 파일을 수정한 후:

```bash
# 1. proto 파일 수정
vim proto/agent.proto

# 2. 생성된 파일 재생성
./scripts/generate-proto.sh

# 3. 생성된 파일도 함께 커밋
git add proto/agent.proto proto/*.pb.go
git commit -m "Update proto definition"
git push origin main

# 4. 새 버전 태그 생성
git tag v0.1.3
git push origin v0.1.3
```

## 대안: 생성된 파일 제외하기

만약 생성된 파일을 제외하고 싶다면:

1. `.gitignore`에 `*.pb.go` 다시 추가
2. GitHub에서 생성된 파일 제거
3. 사용하는 쪽에서 `go generate` 또는 빌드 스크립트로 컴파일

하지만 이 경우 사용자가 매번 컴파일해야 하므로 불편합니다.

