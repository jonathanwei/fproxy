// Code generated by protoc-gen-go.
// source: proto/cookie.proto
// DO NOT EDIT!

package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type AuthCookie struct {
	User string `protobuf:"bytes,1,opt,name=user" json:"user,omitempty"`
}

func (m *AuthCookie) Reset()                    { *m = AuthCookie{} }
func (m *AuthCookie) String() string            { return proto1.CompactTextString(m) }
func (*AuthCookie) ProtoMessage()               {}
func (*AuthCookie) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

type EncryptedMessage struct {
	Nonce      []byte `protobuf:"bytes,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Ciphertext []byte `protobuf:"bytes,2,opt,name=ciphertext,proto3" json:"ciphertext,omitempty"`
}

func (m *EncryptedMessage) Reset()                    { *m = EncryptedMessage{} }
func (m *EncryptedMessage) String() string            { return proto1.CompactTextString(m) }
func (*EncryptedMessage) ProtoMessage()               {}
func (*EncryptedMessage) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func init() {
	proto1.RegisterType((*AuthCookie)(nil), "fproxy.AuthCookie")
	proto1.RegisterType((*EncryptedMessage)(nil), "fproxy.EncryptedMessage")
}

var fileDescriptor2 = []byte{
	// 140 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0x2a, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0xce, 0xcf, 0xcf, 0xce, 0x4c, 0xd5, 0x03, 0x73, 0x84, 0xd8, 0xd2, 0x80, 0x74,
	0x45, 0xa5, 0x92, 0x02, 0x17, 0x97, 0x63, 0x69, 0x49, 0x86, 0x33, 0x58, 0x4e, 0x48, 0x88, 0x8b,
	0xa5, 0xb4, 0x38, 0xb5, 0x48, 0x82, 0x51, 0x81, 0x51, 0x83, 0x33, 0x08, 0xcc, 0x56, 0xf2, 0xe0,
	0x12, 0x70, 0xcd, 0x4b, 0x2e, 0xaa, 0x2c, 0x28, 0x49, 0x4d, 0xf1, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c,
	0x4f, 0x15, 0x12, 0xe1, 0x62, 0xcd, 0xcb, 0xcf, 0x4b, 0x4e, 0x05, 0x2b, 0xe4, 0x09, 0x82, 0x70,
	0x84, 0xe4, 0xb8, 0xb8, 0x92, 0x33, 0x0b, 0x32, 0x52, 0x8b, 0x4a, 0x52, 0x2b, 0x4a, 0x24, 0x98,
	0xc0, 0x52, 0x48, 0x22, 0x4e, 0xec, 0x51, 0xac, 0x60, 0xcb, 0x93, 0xd8, 0xc0, 0x94, 0x31, 0x20,
	0x00, 0x00, 0xff, 0xff, 0x0b, 0x5a, 0x50, 0x4b, 0x99, 0x00, 0x00, 0x00,
}
