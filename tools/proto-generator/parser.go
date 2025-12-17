package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// BinapiType은 파싱된 binapi 타입 정보를 저장합니다.
type BinapiType struct {
	Name        string              // 타입 이름 (예: "SwInterfaceDetails")
	Package     string              // 패키지 이름 (예: "interface")
	Fields      []BinapiField       // 필드 목록
	IsDetails   bool                // Details 메시지인지 여부
	IsReply     bool                // Reply 메시지인지 여부
}

// BinapiField는 binapi 구조체의 필드 정보를 저장합니다.
type BinapiField struct {
	Name     string // 필드 이름 (예: "SwIfIndex")
	Type     string // Go 타입 (예: "interface_types.InterfaceIndex")
	Tag      string // binapi 태그
	JSONName string // JSON 필드 이름
}

// ParseBinapi는 binapi 디렉토리를 파싱하여 타입 정보를 추출합니다.
func ParseBinapi(binapiDir string, config *Config) (map[string]*BinapiType, error) {
	types := make(map[string]*BinapiType)

	// binapi 디렉토리 순회
	err := filepath.Walk(binapiDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// .ba.go 파일만 처리
		if !strings.HasSuffix(path, ".ba.go") {
			return nil
		}

		// 파일 파싱
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			// 파싱 오류는 무시 (일부 파일은 파싱 실패할 수 있음)
			return nil
		}

		// 패키지 이름 추출
		pkgName := node.Name.Name

		// 구조체 타입 찾기
		ast.Inspect(node, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.TypeSpec:
				if st, ok := x.Type.(*ast.StructType); ok {
					typeName := x.Name.Name
					
					// Details로 끝나는 타입만 처리
					if !strings.HasSuffix(typeName, "Details") && 
					   !strings.HasSuffix(typeName, "Reply") {
						return true
					}

					// 설정에 있는 메시지만 처리
					resource := config.FindResourceByBinapi(typeName)
					if resource == nil {
						return true
					}

					// 필드 추출
					fields := extractFields(st, pkgName)
					
					types[typeName] = &BinapiType{
						Name:      typeName,
						Package:   pkgName,
						Fields:    fields,
						IsDetails: strings.HasSuffix(typeName, "Details"),
						IsReply:   strings.HasSuffix(typeName, "Reply"),
					}
				}
			}
			return true
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("binapi 파싱 실패: %w", err)
	}

	return types, nil
}

// extractFields는 구조체에서 필드를 추출합니다.
func extractFields(st *ast.StructType, pkgName string) []BinapiField {
	var fields []BinapiField

	for _, field := range st.Fields.List {
		if len(field.Names) == 0 {
			continue
		}

		fieldName := field.Names[0].Name
		
		// 타입 추출
		typeStr := extractTypeString(field.Type)
		
		// 태그 추출
		var tag, jsonName string
		if field.Tag != nil {
			tag = field.Tag.Value
			jsonName = extractJSONName(tag)
		}

		fields = append(fields, BinapiField{
			Name:     fieldName,
			Type:     typeStr,
			Tag:      tag,
			JSONName: jsonName,
		})
	}

	return fields
}

// extractTypeString은 ast 타입을 문자열로 변환합니다.
func extractTypeString(expr ast.Expr) string {
	switch x := expr.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.SelectorExpr:
		pkg := extractTypeString(x.X)
		sel := x.Sel.Name
		return fmt.Sprintf("%s.%s", pkg, sel)
	case *ast.ArrayType:
		elem := extractTypeString(x.Elt)
		return fmt.Sprintf("[]%s", elem)
	case *ast.MapType:
		key := extractTypeString(x.Key)
		val := extractTypeString(x.Value)
		return fmt.Sprintf("map[%s]%s", key, val)
	default:
		return "unknown"
	}
}

// extractJSONName는 태그에서 JSON 필드 이름을 추출합니다.
func extractJSONName(tag string) string {
	// `json:"field_name,omitempty"` 형식에서 추출
	tag = strings.Trim(tag, "`")
	parts := strings.Split(tag, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, "json:") {
			jsonTag := strings.TrimPrefix(part, "json:")
			jsonTag = strings.Trim(jsonTag, "\"")
			jsonParts := strings.Split(jsonTag, ",")
			if len(jsonParts) > 0 {
				return jsonParts[0]
			}
		}
	}
	return ""
}

