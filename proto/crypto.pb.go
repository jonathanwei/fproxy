// Code generated by protoc-gen-go.
// source: proto/crypto.proto
// DO NOT EDIT!

package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type EncryptedMessage struct {
	Nonce      []byte `protobuf:"bytes,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Ciphertext []byte `protobuf:"bytes,2,opt,name=ciphertext,proto3" json:"ciphertext,omitempty"`
}

func (m *EncryptedMessage) Reset()                    { *m = EncryptedMessage{} }
func (m *EncryptedMessage) String() string            { return proto1.CompactTextString(m) }
func (*EncryptedMessage) ProtoMessage()               {}
func (*EncryptedMessage) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func init() {
	proto1.RegisterType((*EncryptedMessage)(nil), "fproxy.EncryptedMessage")
}

var fileDescriptor2 = []byte{
	// 113 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0x2a, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0x2e, 0xaa, 0x2c, 0x28, 0xc9, 0xd7, 0x03, 0x73, 0x84, 0xd8, 0xd2, 0x80, 0x74,
	0x45, 0xa5, 0x92, 0x07, 0x97, 0x80, 0x6b, 0x1e, 0x58, 0x26, 0x35, 0xc5, 0x37, 0xb5, 0xb8, 0x38,
	0x31, 0x3d, 0x55, 0x48, 0x84, 0x8b, 0x35, 0x2f, 0x3f, 0x2f, 0x39, 0x55, 0x82, 0x51, 0x81, 0x51,
	0x83, 0x27, 0x08, 0xc2, 0x11, 0x92, 0xe3, 0xe2, 0x4a, 0xce, 0x2c, 0xc8, 0x48, 0x2d, 0x2a, 0x49,
	0xad, 0x28, 0x91, 0x60, 0x02, 0x4b, 0x21, 0x89, 0x38, 0xb1, 0x47, 0xb1, 0x82, 0x8d, 0x4e, 0x62,
	0x03, 0x53, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb5, 0x52, 0x89, 0x0c, 0x77, 0x00, 0x00,
	0x00,
}
