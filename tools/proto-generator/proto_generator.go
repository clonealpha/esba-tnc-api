package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GenerateProto는 proto 파일을 생성합니다.
func GenerateProto(binapiTypes map[string]*BinapiType, config *Config, outputDir string) error {
	// 출력 디렉토리 생성
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("출력 디렉토리 생성 실패: %w", err)
	}

	protoFile := filepath.Join(outputDir, "agent_generated.proto")
	
	// 기존 proto 파일 확인 (수동으로 작성된 부분 유지)
	existingProto := filepath.Join(outputDir, "agent.proto")
	hasExisting := false
	if _, err := os.Stat(existingProto); err == nil {
		hasExisting = true
	}

	// proto 파일 생성
	f, err := os.Create(protoFile)
	if err != nil {
		return fmt.Errorf("proto 파일 생성 실패: %w", err)
	}
	defer f.Close()

	// 헤더 작성
	writeProtoHeader(f)

	// 각 리소스에 대해 proto 메시지 생성
	for _, resource := range config.Resources {
		if resource.BinapiMessage == "" {
			continue
		}

		binapiType, ok := binapiTypes[resource.BinapiMessage]
		if !ok {
			// binapi 타입을 찾을 수 없으면 스킵
			continue
		}

		// 개별 메시지 생성
		if err := writeProtoMessage(f, binapiType, &resource); err != nil {
			return fmt.Errorf("proto 메시지 생성 실패 (%s): %w", resource.Name, err)
		}

		// 리스트 메시지 생성
		if resource.ListMessage != "" {
			writeProtoListMessage(f, &resource)
		}
	}

	// 주석 추가
	fmt.Fprintf(f, "\n// 이 파일은 자동 생성되었습니다.\n")
	fmt.Fprintf(f, "// 수정하려면 tools/proto-generator를 사용하세요.\n")

	if hasExisting {
		fmt.Fprintf(os.Stderr, "\n주의: 기존 proto 파일(%s)이 있습니다.\n", existingProto)
		fmt.Fprintf(os.Stderr, "생성된 파일: %s\n", protoFile)
		fmt.Fprintf(os.Stderr, "필요시 수동으로 병합하거나 agent.proto를 업데이트하세요.\n")
	}

	return nil
}

// writeProtoHeader는 proto 파일 헤더를 작성합니다.
func writeProtoHeader(f *os.File) {
	fmt.Fprintf(f, "syntax = \"proto3\";\n\n")
	fmt.Fprintf(f, "package tnc.agent;\n\n")
	fmt.Fprintf(f, "option go_package = \"esba-tnc-api/proto\";\n\n")
	fmt.Fprintf(f, "// 자동 생성된 proto 메시지 정의\n")
	fmt.Fprintf(f, "// 원본: binapi Details 메시지\n\n")
}

// writeProtoMessage는 proto 메시지를 작성합니다.
func writeProtoMessage(f *os.File, binapiType *BinapiType, resource *ResourceConfig) error {
	fmt.Fprintf(f, "// %s 관련 메시지\n", resource.Name)
	fmt.Fprintf(f, "message %s {\n", resource.ProtoMessage)

	fieldNum := 1
	for _, field := range binapiType.Fields {
		// 필드 매핑 확인
		protoFieldName := toProtoFieldName(field.Name)
		protoType := toProtoType(field.Type)

		// 설정에 필드 매핑이 있으면 사용
		if resource.Fields != nil {
			for _, mapping := range resource.Fields {
				if mapping.BinapiField == field.Name {
					protoFieldName = mapping.ProtoField
					break
				}
			}
		}

		// proto 필드 작성
		fmt.Fprintf(f, "  %s %s = %d;\n", protoType, protoFieldName, fieldNum)
		fieldNum++
	}

	fmt.Fprintf(f, "}\n\n")
	return nil
}

// writeProtoListMessage는 proto 리스트 메시지를 작성합니다.
func writeProtoListMessage(f *os.File, resource *ResourceConfig) {
	// 리스트 필드 이름 생성 (복수형)
	listFieldName := toPlural(resource.Name)
	
	fmt.Fprintf(f, "message %s {\n", resource.ListMessage)
	fmt.Fprintf(f, "  repeated %s %s = 1;\n", resource.ProtoMessage, listFieldName)
	fmt.Fprintf(f, "}\n\n")
}

// toProtoFieldName은 Go 필드 이름을 proto 필드 이름으로 변환합니다.
func toProtoFieldName(goName string) string {
	// SwIfIndex -> sw_if_index
	var result strings.Builder
	for i, r := range goName {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// toProtoType는 Go 타입을 proto 타입으로 변환합니다.
func toProtoType(goType string) string {
	// 타입 매핑
	typeMap := map[string]string{
		"uint32":   "uint32",
		"uint16":   "uint32", // proto는 uint16 없음
		"uint8":    "uint32",
		"int32":    "int32",
		"int16":    "int32",
		"int8":     "int32",
		"bool":     "bool",
		"string":   "string",
		"float64":  "double",
		"float32":  "float",
		"[]uint32": "repeated uint32",
		"[]string": "repeated string",
	}

	// 직접 매핑
	if protoType, ok := typeMap[goType]; ok {
		return protoType
	}

	// 특수 타입 처리
	if strings.Contains(goType, "interface_types.InterfaceIndex") {
		return "uint32"
	}
	if strings.Contains(goType, "interface_types.IfStatusFlags") {
		return "uint32"
	}
	if strings.Contains(goType, "ethernet_types.MacAddress") {
		return "string" // .String() 사용
	}
	if strings.Contains(goType, "ip_types.Address") {
		return "string" // .String() 사용
	}
	if strings.Contains(goType, "[]") {
		// 배열 타입
		elemType := strings.TrimPrefix(goType, "[]")
		protoElemType := toProtoType(elemType)
		return fmt.Sprintf("repeated %s", protoElemType)
	}

	// 기본값: string
	return "string"
}

// toPlural는 단수형을 복수형으로 변환합니다.
func toPlural(singular string) string {
	if strings.HasSuffix(singular, "y") {
		return strings.TrimSuffix(singular, "y") + "ies"
	}
	if strings.HasSuffix(singular, "s") || strings.HasSuffix(singular, "x") {
		return singular + "es"
	}
	return singular + "s"
}

