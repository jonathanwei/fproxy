// Code generated by protoc-gen-go.
// source: proto/backend_config.proto
// DO NOT EDIT!

package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type BackendConfig struct {
	ServerAddr string `protobuf:"bytes,1,opt,name=server_addr" json:"server_addr,omitempty"`
}

func (m *BackendConfig) Reset()                    { *m = BackendConfig{} }
func (m *BackendConfig) String() string            { return proto1.CompactTextString(m) }
func (*BackendConfig) ProtoMessage()               {}
func (*BackendConfig) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func init() {
	proto1.RegisterType((*BackendConfig)(nil), "fproxy.proto.BackendConfig")
}

var fileDescriptor1 = []byte{
	// 106 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x92, 0x2a, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0x4a, 0x4c, 0xce, 0x4e, 0xcd, 0x4b, 0x89, 0x4f, 0xce, 0xcf, 0x4b, 0xcb, 0x4c,
	0xd7, 0x03, 0x0b, 0x0a, 0xf1, 0xa4, 0x01, 0xe9, 0x8a, 0x4a, 0x08, 0x4f, 0xc9, 0x80, 0x8b, 0xd7,
	0x09, 0xa2, 0xca, 0x19, 0xac, 0x48, 0x48, 0x9e, 0x8b, 0xbb, 0x38, 0xb5, 0xa8, 0x2c, 0xb5, 0x28,
	0x3e, 0x31, 0x25, 0xa5, 0x48, 0x82, 0x51, 0x81, 0x51, 0x83, 0x33, 0x88, 0x0b, 0x22, 0xe4, 0x08,
	0x14, 0x71, 0x62, 0x8f, 0x62, 0x05, 0x6b, 0x4d, 0x62, 0x03, 0x53, 0xc6, 0x80, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x95, 0xc7, 0x09, 0xb9, 0x6d, 0x00, 0x00, 0x00,
}
