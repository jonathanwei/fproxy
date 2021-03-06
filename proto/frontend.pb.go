// Code generated by protoc-gen-go.
// source: proto/frontend.proto
// DO NOT EDIT!

package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type PortUpdate struct {
	Port int32 `protobuf:"varint,1,opt,name=port" json:"port,omitempty"`
}

func (m *PortUpdate) Reset()                    { *m = PortUpdate{} }
func (m *PortUpdate) String() string            { return proto1.CompactTextString(m) }
func (*PortUpdate) ProtoMessage()               {}
func (*PortUpdate) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{0} }

func init() {
	proto1.RegisterType((*PortUpdate)(nil), "fproxy.PortUpdate")
}

var fileDescriptor4 = []byte{
	// 90 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0x29, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0x2b, 0xca, 0xcf, 0x2b, 0x49, 0xcd, 0x4b, 0xd1, 0x03, 0x73, 0x85, 0xd8, 0xd2,
	0x80, 0x74, 0x45, 0xa5, 0x92, 0x02, 0x17, 0x57, 0x40, 0x7e, 0x51, 0x49, 0x68, 0x41, 0x4a, 0x62,
	0x49, 0xaa, 0x90, 0x10, 0x17, 0x4b, 0x01, 0x90, 0x27, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x1a, 0x04,
	0x66, 0x3b, 0xb1, 0x47, 0xb1, 0x82, 0xb5, 0x24, 0xb1, 0x81, 0x29, 0x63, 0x40, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x1a, 0xc1, 0xac, 0x5a, 0x51, 0x00, 0x00, 0x00,
}
