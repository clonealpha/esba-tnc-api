package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ResourceConfig는 변환할 리소스의 설정을 정의합니다.
type ResourceConfig struct {
	Name           string `yaml:"name"`            // 리소스 이름 (예: "interfaces")
	BinapiMessage  string `yaml:"binapi_message"`  // binapi 메시지 이름 (예: "SwInterfaceDetails")
	ProtoMessage   string `yaml:"proto_message"`    // proto 메시지 이름 (예: "Interface")
	ListMessage    string `yaml:"list_message"`     // proto 리스트 메시지 이름 (예: "InterfaceList")
	Fields         []FieldMapping `yaml:"fields,omitempty"` // 필드 매핑 (선택사항)
}

// FieldMapping은 binapi 필드를 proto 필드로 매핑합니다.
type FieldMapping struct {
	BinapiField string `yaml:"binapi_field"` // binapi 필드 이름
	ProtoField  string `yaml:"proto_field"`  // proto 필드 이름
	Converter   string `yaml:"converter,omitempty"` // 변환 함수 (선택사항)
}

// Config는 proto-generator의 전체 설정을 정의합니다.
type Config struct {
	Resources []ResourceConfig `yaml:"resources"`
}

// LoadConfig는 설정 파일을 로드합니다.
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("설정 파일 파싱 실패: %w", err)
	}

	return &config, nil
}

// DefaultConfig는 기본 설정을 반환합니다.
func DefaultConfig() *Config {
	return &Config{
		Resources: []ResourceConfig{
			{
				Name:          "interfaces",
				BinapiMessage: "SwInterfaceDetails",
				ProtoMessage:  "Interface",
				ListMessage:   "InterfaceList",
			},
			{
				Name:          "neighbors",
				BinapiMessage: "IPNeighborDetails",
				ProtoMessage:  "Neighbor",
				ListMessage:   "NeighborList",
			},
			{
				Name:          "fib",
				BinapiMessage: "IPRouteV2Details",
				ProtoMessage:  "FIBEntry",
				ListMessage:   "FIBList",
			},
			{
				Name:          "acl",
				BinapiMessage: "ACLDetails",
				ProtoMessage:  "ACLEntry",
				ListMessage:   "ACLList",
			},
			{
				Name:          "memif",
				BinapiMessage: "MemifDetails",
				ProtoMessage:  "MemifEntry",
				ListMessage:   "MemifList",
			},
			{
				Name:          "srv6",
				BinapiMessage: "SrLocalsidDetails",
				ProtoMessage:  "SRv6Entry",
				ListMessage:   "SRv6List",
			},
			{
				Name:          "version",
				BinapiMessage: "ShowVersionReply",
				ProtoMessage:  "VersionInfo",
				ListMessage:   "",
			},
			{
				Name:          "hardware",
				BinapiMessage: "",
				ProtoMessage:  "HardwareInfo",
				ListMessage:   "",
			},
			{
				Name:          "ip_addresses",
				BinapiMessage: "IPAddressDetails",
				ProtoMessage:  "IPAddressEntry",
				ListMessage:   "IPAddressList",
			},
			{
				Name:          "l2_fib",
				BinapiMessage: "L2FibTableDetails",
				ProtoMessage:  "L2FIBEntry",
				ListMessage:   "L2FIBList",
			},
			{
				Name:          "bridge_domains",
				BinapiMessage: "BridgeDomainDetails",
				ProtoMessage:  "BridgeDomainEntry",
				ListMessage:   "BridgeDomainList",
			},
			{
				Name:          "vxlan",
				BinapiMessage: "VxlanTunnelDetails",
				ProtoMessage:  "VXLANEntry",
				ListMessage:   "VXLANList",
			},
		},
	}
}

// FindResource는 리소스 이름으로 설정을 찾습니다.
func (c *Config) FindResource(name string) *ResourceConfig {
	for i := range c.Resources {
		if c.Resources[i].Name == name {
			return &c.Resources[i]
		}
	}
	return nil
}

// FindResourceByBinapi는 binapi 메시지 이름으로 설정을 찾습니다.
func (c *Config) FindResourceByBinapi(binapiMsg string) *ResourceConfig {
	for i := range c.Resources {
		if c.Resources[i].BinapiMessage == binapiMsg {
			return &c.Resources[i]
		}
	}
	return nil
}

