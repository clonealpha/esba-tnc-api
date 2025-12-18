# ESBA-TNC-API

공통 gRPC API 정의 및 생성된 코드를 제공하는 공유 모듈입니다.

## 📋 목차

- [목적](#목적)
- [구조](#구조)
- [사용 방법](#사용-방법)
- [API 정의](#api-정의)
- [업데이트 프로세스](#업데이트-프로세스)
- [새로운 리소스 추가 절차](#새로운-리소스-추가-절차)
- [의존성](#의존성)

## 목적

- `esba-tnc-agent`와 `esba-tnc-proxy` 간의 통신을 위한 공통 API 정의
- Protocol Buffers로 정의된 서비스 및 메시지 타입
- 두 프로젝트에서 공유하여 사용하는 생성된 gRPC 코드
- GoVPP 라이브러리 및 binapi 코드 제공
- binapi에서 proto 자동 생성 도구 (proto-generator)

## 구조

```
esba-tnc-api/
├── govpp/                  # GoVPP 라이브러리 (로컬)
│   ├── binapi/            # 생성된 binapi 코드
│   ├── binapigen/         # binapi generator
│   └── ...
├── proto/
│   ├── agent.proto        # gRPC 서비스 정의
│   ├── agent.pb.go        # 생성된 메시지 코드
│   └── agent_grpc.pb.go   # 생성된 서비스 코드
├── tools/
│   └── proto-generator/   # binapi → proto 변환 도구
│       ├── config/
│       │   └── proto.yaml # 변환 설정 파일
│       └── ...
├── scripts/
│   ├── generate-binapi.sh # binapi 생성 스크립트
│   └── generate-proto.sh  # proto 파일 컴파일 스크립트
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
- `CollectVersion`: VPP 버전 정보 수집
- `CollectHardware`: 하드웨어 정보 수집
- `CollectIPAddresses`: IP 주소 정보 수집
- `CollectL2FIB`: L2 FIB 정보 수집
- `CollectBridgeDomains`: Bridge Domain 정보 수집
- `CollectVXLAN`: VXLAN 터널 정보 수집
- `CollectUPFApplications`: UPF 애플리케이션 정보 수집
- `CollectUPFNWI`: UPF Network Instance 정보 수집
- `CollectUPFPFCPEndpoints`: UPF PFCP 엔드포인트 정보 수집
- `CollectUPFPolicies`: UPF 정책 정보 수집
- `CollectUPFNATPools`: UPF NAT 풀 정보 수집

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
- `VersionInfo`: VPP 버전 정보
- `HardwareInfo`: 하드웨어 정보
- `IPAddressList`: IP 주소 목록
- `L2FIBList`: L2 FIB 목록
- `BridgeDomainList`: Bridge Domain 목록
- `VXLANList`: VXLAN 터널 목록
- `UPFApplicationList`: UPF 애플리케이션 목록
- `UPFNWIList`: UPF Network Instance 목록
- `UPFPFCPEndpointList`: UPF PFCP 엔드포인트 목록
- `UPFPolicyList`: UPF 정책 목록
- `UPFNATPoolList`: UPF NAT 풀 목록
- `WatchEventsRequest` / `Event`: 이벤트 스트리밍

## 업데이트 프로세스

### 방법 1: 자동 생성 (권장)

proto-generator를 사용하여 binapi에서 자동으로 proto와 변환 함수를 생성:

```bash
# 전체 빌드 (binapi → proto → 컴파일 → 변환 함수)
./scripts/generate-all.sh --with-binapi

# 또는 단계별 실행
# 1. binapi 생성 (VPP API에서)
./scripts/generate-binapi.sh /usr/share/vpp/api

# 2. proto 생성
cd tools/proto-generator
go run . --proto

# 3. proto 컴파일
cd ../..
./scripts/generate-proto.sh

# 4. 변환 함수 생성
cd tools/proto-generator
go run . --converters
```

### 방법 2: 수동 수정

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

## proto-generator 사용법

proto-generator는 binapi Go 코드를 파싱하여 proto 정의와 변환 함수를 자동 생성합니다.

### 설정 파일

`tools/proto-generator/config/proto.yaml`에서 변환할 리소스를 설정합니다:

```yaml
resources:
  - name: interfaces
    binapi_message: SwInterfaceDetails
    proto_message: Interface
    list_message: InterfaceList
    fields:
      - binapi_field: SwIfIndex
        proto_field: sw_if_index
```

### 실행

```bash
cd tools/proto-generator

# proto만 생성
go run . --proto

# 변환 함수만 생성
go run . --converters

# 둘 다 생성
go run . --proto --converters
```

### 출력

- `proto/agent_generated.proto`: 생성된 proto 정의
- `../esba-tnc-agent/agent/grpc/handler/converters_gen.go`: 생성된 변환 함수

## 새로운 리소스 추가 절차

새로운 VPP 리소스를 `esba-tnc-api`에 추가할 때는 다음 단계를 따라야 합니다:

### 1. proto-generator 설정 파일 업데이트

`tools/proto-generator/config/proto.yaml`에 새로운 리소스를 추가합니다.

**예시**: `tools/proto-generator/config/proto.yaml`

```yaml
resources:
  # ... 기존 리소스들 ...
  
  # 새로운 리소스 추가
  - name: new_resource
    binapi_message: NewResourceDetails
    proto_message: NewResource
    list_message: NewResourceList
    fields:
      - binapi_field: ID
        proto_field: id
        converter: "uint32({{field}})"
      - binapi_field: Name
        proto_field: name
        converter: "string(bytes.Trim({{field}}, \"\\x00\"))"
      # ... 추가 필드들 ...
```

### 2. Proto 정의 추가 (수동 또는 자동 생성)

#### 방법 A: proto-generator 사용 (권장)

```bash
cd tools/proto-generator
go run . --binapi-dir=../govpp/binapi --output=../proto --config=config/proto.yaml --proto
```

이 명령은 `proto/agent_generated.proto`에 새로운 리소스 정의를 생성합니다.

#### 방법 B: 수동 추가

`proto/agent.proto`에 직접 추가:

```protobuf
// RPC 메서드 추가
rpc CollectNewResource(CollectRequest) returns (NewResourceList);

// 메시지 타입 추가
message NewResourceList {
  repeated NewResourceEntry resources = 1;
}

message NewResourceEntry {
  uint32 id = 1;
  string name = 2;
  // ... 추가 필드들 ...
}
```

### 3. Proto 파일 컴파일

```bash
./scripts/generate-proto.sh
```

또는 수동으로:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/agent.proto
```

이 명령은 다음 파일들을 생성합니다:
- `proto/agent.pb.go`: 메시지 타입 코드
- `proto/agent_grpc.pb.go`: 서비스 코드

### 4. 버전 업데이트 및 태그 생성

변경사항을 커밋하고 새 버전 태그를 생성합니다:

```bash
# 변경사항 커밋
git add .
git commit -m "feat: Add NewResource collection support"

# 새 버전 태그 생성 (예: v0.1.5)
git tag -a v0.1.5 -m "feat: Add NewResource collection support"

# GitHub에 push
git push origin main
git push origin v0.1.5
```

### 5. 의존성 업데이트

다른 프로젝트(`esba-tnc-agent`, `esba-tnc-proxy`)에서 새 버전을 사용하도록 업데이트:

```bash
# esba-tnc-agent/go.mod
require github.com/clonealpha/esba-tnc-api v0.1.5

# esba-tnc-proxy/go.mod
require github.com/clonealpha/esba-tnc-api v0.1.5
```

### 6. 변환 함수 생성 (선택사항)

proto-generator를 사용하여 변환 함수를 자동 생성할 수 있습니다:

```bash
cd tools/proto-generator
go run . --binapi-dir=../govpp/binapi --output=../esba-tnc-agent/agent/grpc/handler --config=config/proto.yaml --converters
```

이 명령은 `esba-tnc-agent/agent/grpc/handler/converters_gen.go`에 변환 함수를 생성합니다.

### 체크리스트

새 리소스 추가 시 다음 항목을 확인하세요:

- [ ] `tools/proto-generator/config/proto.yaml`에 리소스 설정 추가
- [ ] `proto/agent.proto`에 RPC 메서드 및 메시지 타입 추가 (수동 또는 자동)
- [ ] Proto 파일 컴파일 완료 (`agent.pb.go`, `agent_grpc.pb.go` 생성 확인)
- [ ] 빌드 테스트 (`go build ./...`)
- [ ] 버전 태그 생성 및 GitHub push
- [ ] `esba-tnc-agent`에서 새 버전 사용하도록 업데이트
- [ ] `esba-tnc-proxy`에서 새 버전 사용하도록 업데이트

### 주의사항

1. **버전 관리**: 새로운 리소스를 추가할 때마다 버전을 올려야 합니다 (예: v0.1.4 → v0.1.5)
2. **하위 호환성**: 기존 RPC 메서드나 메시지 타입을 변경할 때는 주의해야 합니다
3. **proto-generator 사용**: 가능하면 proto-generator를 사용하여 일관성을 유지하세요
4. **문서 업데이트**: README의 RPC 메서드 목록과 메시지 타입 목록을 업데이트하세요

## 의존성

- `github.com/clonealpha/esba-tnc-api/govpp`: GoVPP 라이브러리 (서브모듈)
- `google.golang.org/grpc`: gRPC 라이브러리
- `google.golang.org/protobuf`: Protocol Buffers 라이브러리

### GoVPP 사용

`esba-tnc-api`는 GoVPP를 **서브모듈**(`govpp/`)로 제공합니다. 모듈 경로는 `github.com/clonealpha/esba-tnc-api/govpp`입니다.

다른 프로젝트에서는 `esba-tnc-api/govpp`를 직접 import합니다:

```go
// esba-tnc-agent/go.mod
require (
    github.com/clonealpha/esba-tnc-api v0.1.5
    github.com/clonealpha/esba-tnc-api/govpp v0.1.5
)
```

코드에서는 다음과 같이 import합니다:

```go
import (
    "github.com/clonealpha/esba-tnc-api/govpp/api"
    "github.com/clonealpha/esba-tnc-api/govpp/binapi/interface"
)
```

## 참고 자료

- [Protocol Buffers 문서](https://developers.google.com/protocol-buffers)
- [gRPC 문서](https://grpc.io/docs/)
- [esba-tnc-agent](../esba-tnc-agent/README.md): gRPC 서버 구현
- [esba-tnc-proxy](../esba-tnc-proxy/README.md): gRPC 클라이언트 구현
