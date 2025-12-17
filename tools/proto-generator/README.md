# proto-generator

binapi Go 코드에서 proto 정의와 변환 함수를 자동 생성하는 도구입니다.

## 목적

- binapi Details 메시지를 proto 메시지로 자동 변환
- binapi → proto 변환 함수 자동 생성
- binapi 업데이트 시 자동 반영

## 사용법

### 기본 사용

```bash
cd tools/proto-generator

# proto 생성
go run . --proto

# 변환 함수 생성
go run . --converters

# 둘 다 생성
go run . --proto --converters
```

### 옵션

- `--binapi-dir`: binapi 디렉토리 경로 (기본값: ../govpp/binapi)
- `--output`: proto 출력 디렉토리 (기본값: ../proto)
- `--config`: 설정 파일 경로 (기본값: config/proto.yaml)
- `--proto`: proto 파일 생성
- `--converters`: 변환 함수 생성

## 설정 파일

`config/proto.yaml`에서 변환할 리소스를 정의합니다:

```yaml
resources:
  - name: interfaces
    binapi_message: SwInterfaceDetails
    proto_message: Interface
    list_message: InterfaceList
    fields:
      - binapi_field: SwIfIndex
        proto_field: sw_if_index
        converter: "uint32({{field}})"
```

### 필드

- `name`: 리소스 이름
- `binapi_message`: binapi 메시지 이름 (예: "SwInterfaceDetails")
- `proto_message`: proto 메시지 이름 (예: "Interface")
- `list_message`: proto 리스트 메시지 이름 (예: "InterfaceList")
- `fields`: 필드 매핑 (선택사항)

## 출력

### proto 파일

`proto/agent_generated.proto`에 생성됩니다:

```proto
message Interface {
  uint32 sw_if_index = 1;
  string interface_name = 2;
  // ...
}
```

### 변환 함수

`../esba-tnc-agent/agent/grpc/handler/converters_gen.go`에 생성됩니다:

```go
func convertInterfaces(data interface{}) (interface{}, error) {
    // ...
}
```

## 타입 매핑

자동 타입 변환:

- `uint32`, `uint16`, `uint8` → `uint32`
- `int32`, `int16`, `int8` → `int32`
- `bool` → `bool`
- `string` → `string`
- `interface_types.InterfaceIndex` → `uint32`
- `ethernet_types.MacAddress` → `string` (`.String()` 사용)
- `ip_types.Address` → `string` (`.String()` 사용)

## 제한사항

- 복잡한 타입 변환은 수동으로 설정 파일에 추가 필요
- 중첩된 구조체는 현재 지원하지 않음
- 커스텀 변환 로직은 설정 파일의 `converter` 필드 사용

