# ESBA-TNC-API

공통 gRPC API 정의 및 생성된 코드를 제공하는 공유 모듈입니다.

## 📋 목차

- [목적](#목적)
- [구조](#구조)
- [사용 방법](#사용-방법)
- [API 정의](#api-정의)
- [업데이트 프로세스](#업데이트-프로세스)
- [의존성](#의존성)

## 목적

- `esba-tnc-agent`와 `esba-tnc-proxy` 간의 통신을 위한 공통 API 정의
- Protocol Buffers로 정의된 서비스 및 메시지 타입
- 두 프로젝트에서 공유하여 사용하는 생성된 gRPC 코드

## 구조

```
esba-tnc-api/
├── proto/
│   ├── agent.proto          # gRPC 서비스 정의
│   ├── agent.pb.go         # 생성된 메시지 코드
│   └── agent_grpc.pb.go    # 생성된 서비스 코드
├── scripts/
│   └── generate-proto.sh   # proto 파일 컴파일 스크립트
├── go.mod
├── go.sum
└── README.md
```

## 사용 방법

### 1. Protocol Buffers 컴파일

```bash
cd esba-tnc-api
./scripts/generate-proto.sh
```

또는 수동으로:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/agent.proto
```

### 2. 다른 프로젝트에서 사용

#### esba-tnc-agent/go.mod

```go
require esba-tnc-api v0.0.0

replace esba-tnc-api => ../esba-tnc-api
```

#### esba-tnc-proxy/go.mod

```go
require esba-tnc-api v0.0.0

replace esba-tnc-api => ../esba-tnc-api
```

#### 코드에서 import

```go
import (
    pb "esba-tnc-api/proto"
)

// 서버에서
proto.RegisterTncAgentServer(grpcServer, &server{})

// 클라이언트에서
client := pb.NewTncAgentClient(conn)
```

## API 정의

### 서비스

- `TncAgent`: VPP 데이터 수집 및 이벤트 스트리밍 서비스

### RPC 메서드

#### 조회 기능 (Unary RPC)

- `HealthCheck`: 서버 상태 확인
- `CollectInterfaces`: 인터페이스 정보 수집
- `CollectNeighbors`: Neighbor 정보 수집
- `CollectFIB`: FIB 정보 수집
- `CollectACL`: ACL 정보 수집
- `CollectMemif`: Memif 정보 수집
- `CollectSRv6`: SRv6 정보 수집

#### 이벤트 스트리밍 (Server Streaming)

- `WatchEvents`: VPP 이벤트 실시간 스트리밍

### 메시지 타입

주요 메시지 타입은 `proto/agent.proto`에 정의되어 있습니다:

- `HealthCheckRequest` / `HealthCheckResponse`: 헬스 체크
- `CollectRequest`: 리소스 수집 요청
- `InterfaceList`: 인터페이스 목록
- `NeighborList`: Neighbor 목록
- `FIBList`: FIB 목록
- `ACLList`: ACL 목록
- `MemifList`: Memif 목록
- `SRv6List`: SRv6 목록
- `WatchEventsRequest` / `Event`: 이벤트 스트리밍

## 업데이트 프로세스

1. `proto/agent.proto` 파일 수정
2. `./scripts/generate-proto.sh` 실행하여 코드 재생성
3. 변경사항 커밋
4. `esba-tnc-agent`와 `esba-tnc-proxy`에서 `go mod tidy` 실행

### 예시

```bash
# 1. proto 파일 수정
vim proto/agent.proto

# 2. 코드 재생성
./scripts/generate-proto.sh

# 3. 변경사항 확인
git diff proto/

# 4. 다른 프로젝트에서 업데이트
cd ../esba-tnc-agent
go mod tidy

cd ../esba-tnc-proxy
go mod tidy
```

## 의존성

- `google.golang.org/grpc`: gRPC 라이브러리
- `google.golang.org/protobuf`: Protocol Buffers 라이브러리

## 참고 자료

- [Protocol Buffers 문서](https://developers.google.com/protocol-buffers)
- [gRPC 문서](https://grpc.io/docs/)
- [esba-tnc-agent](../esba-tnc-agent/README.md): gRPC 서버 구현
- [esba-tnc-proxy](../esba-tnc-proxy/README.md): gRPC 클라이언트 구현
