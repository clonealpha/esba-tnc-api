package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	binapiDir    = flag.String("binapi-dir", "", "binapi 디렉토리 경로 (예: ../govpp/binapi)")
	outputDir    = flag.String("output", "proto", "출력 디렉토리")
	configFile   = flag.String("config", "config/proto.yaml", "설정 파일 경로")
	generateProto = flag.Bool("proto", true, "proto 파일 생성")
	generateConverters = flag.Bool("converters", false, "변환 함수 생성")
)

func main() {
	flag.Parse()

	if *binapiDir == "" {
		// 기본값: esba-tnc-api/govpp/binapi
		wd, _ := os.Getwd()
		*binapiDir = filepath.Join(wd, "..", "govpp", "binapi")
	}

	// binapi 디렉토리 확인
	if _, err := os.Stat(*binapiDir); os.IsNotExist(err) {
		log.Fatalf("binapi 디렉토리를 찾을 수 없습니다: %s", *binapiDir)
	}

	// 설정 파일 로드
	config, err := LoadConfig(*configFile)
	if err != nil {
		log.Printf("설정 파일을 찾을 수 없어 기본 설정을 사용합니다: %v", err)
		config = DefaultConfig()
	}

	// binapi 파싱
	log.Printf("binapi 파싱 중: %s", *binapiDir)
	binapiTypes, err := ParseBinapi(*binapiDir, config)
	if err != nil {
		log.Fatalf("binapi 파싱 실패: %v", err)
	}

	log.Printf("파싱된 타입 수: %d", len(binapiTypes))

	// proto 생성
	if *generateProto {
		log.Printf("proto 파일 생성 중: %s", *outputDir)
		if err := GenerateProto(binapiTypes, config, *outputDir); err != nil {
			log.Fatalf("proto 생성 실패: %v", err)
		}
		log.Println("proto 파일 생성 완료")
	}

	// 변환 함수 생성
	if *generateConverters {
		log.Printf("변환 함수 생성 중...")
		converterOutput := filepath.Join("..", "esba-tnc-agent", "agent", "grpc", "handler", "converters_gen.go")
		if err := GenerateConverters(binapiTypes, config, converterOutput); err != nil {
			log.Fatalf("변환 함수 생성 실패: %v", err)
		}
		log.Printf("변환 함수 생성 완료: %s", converterOutput)
	}

	fmt.Println("완료!")
}

