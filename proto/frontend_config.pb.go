// Code generated by protoc-gen-go.
// source: proto/frontend_config.proto
// DO NOT EDIT!

package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type FrontendConfig struct {
	TcpProxyRoute      []*TCPProxyRoute        `protobuf:"bytes,1,rep,name=tcp_proxy_route" json:"tcp_proxy_route,omitempty"`
	Server             *FrontendConfig_Server  `protobuf:"bytes,2,opt,name=server" json:"server,omitempty"`
	Backend            *FrontendConfig_Backend `protobuf:"bytes,3,opt,name=backend" json:"backend,omitempty"`
	OauthConfig        *OauthConfig            `protobuf:"bytes,4,opt,name=oauth_config" json:"oauth_config,omitempty"`
	AuthCookieKey      []byte                  `protobuf:"bytes,5,opt,name=auth_cookie_key,proto3" json:"auth_cookie_key,omitempty"`
	AuthCookieInsecure bool                    `protobuf:"varint,6,opt,name=auth_cookie_insecure" json:"auth_cookie_insecure,omitempty"`
	EmailToUserId      map[string]string       `protobuf:"bytes,7,rep,name=email_to_user_id" json:"email_to_user_id,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *FrontendConfig) Reset()                    { *m = FrontendConfig{} }
func (m *FrontendConfig) String() string            { return proto1.CompactTextString(m) }
func (*FrontendConfig) ProtoMessage()               {}
func (*FrontendConfig) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{0} }

func (m *FrontendConfig) GetTcpProxyRoute() []*TCPProxyRoute {
	if m != nil {
		return m.TcpProxyRoute
	}
	return nil
}

func (m *FrontendConfig) GetServer() *FrontendConfig_Server {
	if m != nil {
		return m.Server
	}
	return nil
}

func (m *FrontendConfig) GetBackend() *FrontendConfig_Backend {
	if m != nil {
		return m.Backend
	}
	return nil
}

func (m *FrontendConfig) GetOauthConfig() *OauthConfig {
	if m != nil {
		return m.OauthConfig
	}
	return nil
}

func (m *FrontendConfig) GetEmailToUserId() map[string]string {
	if m != nil {
		return m.EmailToUserId
	}
	return nil
}

// TODO: this is repeated 3 different times between frontend and backend
// config. Consolidate them.
type FrontendConfig_Server struct {
	Addr string `protobuf:"bytes,1,opt,name=addr" json:"addr,omitempty"`
	// Types that are valid to be assigned to Security:
	//	*FrontendConfig_Server_Tls
	//	*FrontendConfig_Server_Insecure
	Security isFrontendConfig_Server_Security `protobuf_oneof:"security"`
}

func (m *FrontendConfig_Server) Reset()                    { *m = FrontendConfig_Server{} }
func (m *FrontendConfig_Server) String() string            { return proto1.CompactTextString(m) }
func (*FrontendConfig_Server) ProtoMessage()               {}
func (*FrontendConfig_Server) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{0, 0} }

type isFrontendConfig_Server_Security interface {
	isFrontendConfig_Server_Security()
}

type FrontendConfig_Server_Tls struct {
	Tls *TLSConfig `protobuf:"bytes,2,opt,name=tls,oneof"`
}
type FrontendConfig_Server_Insecure struct {
	Insecure bool `protobuf:"varint,3,opt,name=insecure,oneof"`
}

func (*FrontendConfig_Server_Tls) isFrontendConfig_Server_Security()      {}
func (*FrontendConfig_Server_Insecure) isFrontendConfig_Server_Security() {}

func (m *FrontendConfig_Server) GetSecurity() isFrontendConfig_Server_Security {
	if m != nil {
		return m.Security
	}
	return nil
}

func (m *FrontendConfig_Server) GetTls() *TLSConfig {
	if x, ok := m.GetSecurity().(*FrontendConfig_Server_Tls); ok {
		return x.Tls
	}
	return nil
}

func (m *FrontendConfig_Server) GetInsecure() bool {
	if x, ok := m.GetSecurity().(*FrontendConfig_Server_Insecure); ok {
		return x.Insecure
	}
	return false
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*FrontendConfig_Server) XXX_OneofFuncs() (func(msg proto1.Message, b *proto1.Buffer) error, func(msg proto1.Message, tag, wire int, b *proto1.Buffer) (bool, error), func(msg proto1.Message) (n int), []interface{}) {
	return _FrontendConfig_Server_OneofMarshaler, _FrontendConfig_Server_OneofUnmarshaler, _FrontendConfig_Server_OneofSizer, []interface{}{
		(*FrontendConfig_Server_Tls)(nil),
		(*FrontendConfig_Server_Insecure)(nil),
	}
}

func _FrontendConfig_Server_OneofMarshaler(msg proto1.Message, b *proto1.Buffer) error {
	m := msg.(*FrontendConfig_Server)
	// security
	switch x := m.Security.(type) {
	case *FrontendConfig_Server_Tls:
		b.EncodeVarint(2<<3 | proto1.WireBytes)
		if err := b.EncodeMessage(x.Tls); err != nil {
			return err
		}
	case *FrontendConfig_Server_Insecure:
		t := uint64(0)
		if x.Insecure {
			t = 1
		}
		b.EncodeVarint(3<<3 | proto1.WireVarint)
		b.EncodeVarint(t)
	case nil:
	default:
		return fmt.Errorf("FrontendConfig_Server.Security has unexpected type %T", x)
	}
	return nil
}

func _FrontendConfig_Server_OneofUnmarshaler(msg proto1.Message, tag, wire int, b *proto1.Buffer) (bool, error) {
	m := msg.(*FrontendConfig_Server)
	switch tag {
	case 2: // security.tls
		if wire != proto1.WireBytes {
			return true, proto1.ErrInternalBadWireType
		}
		msg := new(TLSConfig)
		err := b.DecodeMessage(msg)
		m.Security = &FrontendConfig_Server_Tls{msg}
		return true, err
	case 3: // security.insecure
		if wire != proto1.WireVarint {
			return true, proto1.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.Security = &FrontendConfig_Server_Insecure{x != 0}
		return true, err
	default:
		return false, nil
	}
}

func _FrontendConfig_Server_OneofSizer(msg proto1.Message) (n int) {
	m := msg.(*FrontendConfig_Server)
	// security
	switch x := m.Security.(type) {
	case *FrontendConfig_Server_Tls:
		s := proto1.Size(x.Tls)
		n += proto1.SizeVarint(2<<3 | proto1.WireBytes)
		n += proto1.SizeVarint(uint64(s))
		n += s
	case *FrontendConfig_Server_Insecure:
		n += proto1.SizeVarint(3<<3 | proto1.WireVarint)
		n += 1
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// TODO: make this repeated.
type FrontendConfig_Backend struct {
	Addr string `protobuf:"bytes,1,opt,name=addr" json:"addr,omitempty"`
	// Types that are valid to be assigned to Security:
	//	*FrontendConfig_Backend_Tls
	//	*FrontendConfig_Backend_Insecure
	Security isFrontendConfig_Backend_Security `protobuf_oneof:"security"`
}

func (m *FrontendConfig_Backend) Reset()                    { *m = FrontendConfig_Backend{} }
func (m *FrontendConfig_Backend) String() string            { return proto1.CompactTextString(m) }
func (*FrontendConfig_Backend) ProtoMessage()               {}
func (*FrontendConfig_Backend) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{0, 1} }

type isFrontendConfig_Backend_Security interface {
	isFrontendConfig_Backend_Security()
}

type FrontendConfig_Backend_Tls struct {
	Tls *TLSConfig `protobuf:"bytes,2,opt,name=tls,oneof"`
}
type FrontendConfig_Backend_Insecure struct {
	Insecure bool `protobuf:"varint,3,opt,name=insecure,oneof"`
}

func (*FrontendConfig_Backend_Tls) isFrontendConfig_Backend_Security()      {}
func (*FrontendConfig_Backend_Insecure) isFrontendConfig_Backend_Security() {}

func (m *FrontendConfig_Backend) GetSecurity() isFrontendConfig_Backend_Security {
	if m != nil {
		return m.Security
	}
	return nil
}

func (m *FrontendConfig_Backend) GetTls() *TLSConfig {
	if x, ok := m.GetSecurity().(*FrontendConfig_Backend_Tls); ok {
		return x.Tls
	}
	return nil
}

func (m *FrontendConfig_Backend) GetInsecure() bool {
	if x, ok := m.GetSecurity().(*FrontendConfig_Backend_Insecure); ok {
		return x.Insecure
	}
	return false
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*FrontendConfig_Backend) XXX_OneofFuncs() (func(msg proto1.Message, b *proto1.Buffer) error, func(msg proto1.Message, tag, wire int, b *proto1.Buffer) (bool, error), func(msg proto1.Message) (n int), []interface{}) {
	return _FrontendConfig_Backend_OneofMarshaler, _FrontendConfig_Backend_OneofUnmarshaler, _FrontendConfig_Backend_OneofSizer, []interface{}{
		(*FrontendConfig_Backend_Tls)(nil),
		(*FrontendConfig_Backend_Insecure)(nil),
	}
}

func _FrontendConfig_Backend_OneofMarshaler(msg proto1.Message, b *proto1.Buffer) error {
	m := msg.(*FrontendConfig_Backend)
	// security
	switch x := m.Security.(type) {
	case *FrontendConfig_Backend_Tls:
		b.EncodeVarint(2<<3 | proto1.WireBytes)
		if err := b.EncodeMessage(x.Tls); err != nil {
			return err
		}
	case *FrontendConfig_Backend_Insecure:
		t := uint64(0)
		if x.Insecure {
			t = 1
		}
		b.EncodeVarint(3<<3 | proto1.WireVarint)
		b.EncodeVarint(t)
	case nil:
	default:
		return fmt.Errorf("FrontendConfig_Backend.Security has unexpected type %T", x)
	}
	return nil
}

func _FrontendConfig_Backend_OneofUnmarshaler(msg proto1.Message, tag, wire int, b *proto1.Buffer) (bool, error) {
	m := msg.(*FrontendConfig_Backend)
	switch tag {
	case 2: // security.tls
		if wire != proto1.WireBytes {
			return true, proto1.ErrInternalBadWireType
		}
		msg := new(TLSConfig)
		err := b.DecodeMessage(msg)
		m.Security = &FrontendConfig_Backend_Tls{msg}
		return true, err
	case 3: // security.insecure
		if wire != proto1.WireVarint {
			return true, proto1.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.Security = &FrontendConfig_Backend_Insecure{x != 0}
		return true, err
	default:
		return false, nil
	}
}

func _FrontendConfig_Backend_OneofSizer(msg proto1.Message) (n int) {
	m := msg.(*FrontendConfig_Backend)
	// security
	switch x := m.Security.(type) {
	case *FrontendConfig_Backend_Tls:
		s := proto1.Size(x.Tls)
		n += proto1.SizeVarint(2<<3 | proto1.WireBytes)
		n += proto1.SizeVarint(uint64(s))
		n += s
	case *FrontendConfig_Backend_Insecure:
		n += proto1.SizeVarint(3<<3 | proto1.WireVarint)
		n += 1
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type TCPProxyRoute struct {
	Listen string `protobuf:"bytes,1,opt,name=listen" json:"listen,omitempty"`
	Dial   string `protobuf:"bytes,2,opt,name=dial" json:"dial,omitempty"`
}

func (m *TCPProxyRoute) Reset()                    { *m = TCPProxyRoute{} }
func (m *TCPProxyRoute) String() string            { return proto1.CompactTextString(m) }
func (*TCPProxyRoute) ProtoMessage()               {}
func (*TCPProxyRoute) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{1} }

type OauthConfig struct {
	ClientId     string `protobuf:"bytes,1,opt,name=client_id" json:"client_id,omitempty"`
	ClientSecret string `protobuf:"bytes,2,opt,name=client_secret" json:"client_secret,omitempty"`
	RedirectUrl  string `protobuf:"bytes,3,opt,name=redirect_url" json:"redirect_url,omitempty"`
}

func (m *OauthConfig) Reset()                    { *m = OauthConfig{} }
func (m *OauthConfig) String() string            { return proto1.CompactTextString(m) }
func (*OauthConfig) ProtoMessage()               {}
func (*OauthConfig) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{2} }

func init() {
	proto1.RegisterType((*FrontendConfig)(nil), "fproxy.FrontendConfig")
	proto1.RegisterType((*FrontendConfig_Server)(nil), "fproxy.FrontendConfig.Server")
	proto1.RegisterType((*FrontendConfig_Backend)(nil), "fproxy.FrontendConfig.Backend")
	proto1.RegisterType((*TCPProxyRoute)(nil), "fproxy.TCPProxyRoute")
	proto1.RegisterType((*OauthConfig)(nil), "fproxy.OauthConfig")
}

var fileDescriptor4 = []byte{
	// 481 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xbc, 0x53, 0x4b, 0x6f, 0xd3, 0x40,
	0x10, 0x26, 0x4d, 0xe3, 0x38, 0x93, 0x98, 0x94, 0xa1, 0x20, 0x2b, 0x05, 0x54, 0x82, 0x40, 0xe5,
	0x12, 0x10, 0x08, 0x54, 0x81, 0x90, 0x50, 0xaa, 0x22, 0x22, 0x90, 0xa8, 0x36, 0xed, 0x85, 0x8b,
	0xe5, 0xda, 0x13, 0x58, 0xc5, 0xf5, 0x56, 0x9b, 0x75, 0x85, 0x7f, 0x21, 0x7f, 0x8b, 0x7d, 0xe5,
	0x51, 0xa1, 0x5e, 0x7b, 0xf2, 0xcc, 0xf7, 0x98, 0xf1, 0xcc, 0xd8, 0xb0, 0x77, 0x29, 0x85, 0x12,
	0xaf, 0x66, 0x52, 0x94, 0x8a, 0xca, 0x3c, 0xc9, 0x44, 0x39, 0xe3, 0xbf, 0x46, 0x16, 0xc5, 0x60,
	0xa6, 0x9f, 0x7f, 0xea, 0x41, 0xdf, 0x89, 0x54, 0xb1, 0x70, 0xc4, 0xf0, 0x6f, 0x0b, 0xee, 0x7e,
	0xf1, 0x96, 0x23, 0xeb, 0xc0, 0x4f, 0xd0, 0x57, 0xd9, 0x65, 0x62, 0x0d, 0x89, 0x14, 0x95, 0xa2,
	0xb8, 0xb1, 0xdf, 0x3c, 0xe8, 0xbe, 0x79, 0x30, 0x72, 0x55, 0x46, 0xa7, 0x47, 0x27, 0x27, 0x26,
	0x60, 0x86, 0x64, 0x91, 0x56, 0xaf, 0x53, 0x7c, 0x07, 0xc1, 0x82, 0xe4, 0x15, 0xc9, 0x78, 0x6b,
	0xbf, 0xa1, 0x5d, 0x8f, 0x97, 0xae, 0xeb, 0x6d, 0x46, 0x53, 0x2b, 0x62, 0x5e, 0x8c, 0x87, 0xd0,
	0x3e, 0x4f, 0xb3, 0xb9, 0xe6, 0xe3, 0xa6, 0xf5, 0x3d, 0xb9, 0xc1, 0x37, 0x76, 0x2a, 0xb6, 0x94,
	0xe3, 0x7b, 0xe8, 0x89, 0xb4, 0x52, 0xbf, 0xfd, 0xc4, 0xf1, 0xb6, 0xb5, 0xdf, 0x5f, 0xda, 0x7f,
	0x18, 0xce, 0x79, 0x59, 0x57, 0xac, 0x13, 0x7c, 0x01, 0x7d, 0x6f, 0x13, 0x73, 0x4e, 0xc9, 0x9c,
	0xea, 0xb8, 0xa5, 0xad, 0x3d, 0x16, 0x39, 0x91, 0x41, 0xbf, 0x51, 0x8d, 0xaf, 0x61, 0x77, 0x53,
	0xc7, 0xcb, 0x05, 0x65, 0x95, 0xa4, 0x38, 0xd0, 0xe2, 0x90, 0xe1, 0x5a, 0x3c, 0xf1, 0x0c, 0x32,
	0xd8, 0xa1, 0x8b, 0x94, 0x17, 0x89, 0x12, 0x49, 0xa5, 0xe7, 0x4b, 0x78, 0x1e, 0xb7, 0xed, 0x0a,
	0x5f, 0xde, 0x30, 0xd4, 0xb1, 0x91, 0x9f, 0x8a, 0x33, 0x2d, 0x9e, 0xe4, 0xc7, 0xa5, 0x92, 0x35,
	0x8b, 0x68, 0x13, 0x1b, 0x5c, 0x40, 0xe0, 0x36, 0x86, 0x08, 0xdb, 0x69, 0x9e, 0x4b, 0x7d, 0x94,
	0xc6, 0x41, 0x87, 0xd9, 0x18, 0x9f, 0x43, 0x53, 0xdf, 0xd4, 0x6f, 0xfc, 0xde, 0xea, 0x4e, 0xdf,
	0xa7, 0xae, 0xfe, 0xd7, 0x3b, 0xcc, 0xf0, 0xf8, 0x08, 0xc2, 0xd5, 0xeb, 0x9b, 0x2d, 0x87, 0x9a,
	0x58, 0x21, 0x63, 0x80, 0xd0, 0x46, 0x5c, 0xd5, 0x83, 0x12, 0xda, 0x7e, 0xd1, 0xb7, 0xd3, 0xef,
	0x33, 0xe0, 0xff, 0x3b, 0xc0, 0x1d, 0x68, 0x9a, 0xb3, 0xb8, 0xce, 0x26, 0xc4, 0x5d, 0x68, 0x5d,
	0xa5, 0x45, 0x45, 0xb6, 0x75, 0x87, 0xb9, 0xe4, 0xc3, 0xd6, 0x61, 0x63, 0xf8, 0x11, 0xa2, 0x6b,
	0xdf, 0x25, 0x3e, 0x84, 0xa0, 0xe0, 0x0b, 0xbd, 0x64, 0xef, 0xf7, 0x99, 0x99, 0x27, 0xe7, 0x69,
	0xe1, 0x2b, 0xd8, 0x78, 0x28, 0xa1, 0xbb, 0xf1, 0x9d, 0xe0, 0x1e, 0x74, 0xb2, 0x82, 0x53, 0xa9,
	0xcc, 0xe5, 0x9c, 0x3b, 0x74, 0xc0, 0x24, 0xc7, 0x67, 0x10, 0x79, 0x52, 0xbf, 0xbd, 0x24, 0xe5,
	0x0b, 0xf5, 0x1c, 0x38, 0xb5, 0x18, 0x3e, 0x85, 0x9e, 0xa4, 0x9c, 0x4b, 0xca, 0x54, 0x52, 0xc9,
	0xc2, 0x4e, 0xdf, 0x61, 0xdd, 0x25, 0x76, 0x26, 0x8b, 0x71, 0xfb, 0x67, 0xcb, 0xfe, 0x83, 0xe7,
	0x81, 0x7d, 0xbc, 0xfd, 0x17, 0x00, 0x00, 0xff, 0xff, 0x04, 0x1d, 0x8f, 0x4e, 0xc2, 0x03, 0x00,
	0x00,
}
