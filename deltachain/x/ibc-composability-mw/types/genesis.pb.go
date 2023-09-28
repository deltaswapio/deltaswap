// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ibc-composability-mw/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// GenesisState defines the ibc-composability-mw genesis state
type GenesisState struct {
	// key - information about modified packet: src_channel
	// (parsedReceiver.Channel), src_port (parsedReceiver.Port), sequence value -
	// bytes are the packet data bytes as they came in
	TransposedDataInFlight map[string][]byte `protobuf:"bytes,1,rep,name=transposed_data_in_flight,json=transposedDataInFlight,proto3" json:"transposed_data_in_flight" yaml:"transposed_data_in_flight" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_17a5790dbeee9b79, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetTransposedDataInFlight() map[string][]byte {
	if m != nil {
		return m.TransposedDataInFlight
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "wormhole_foundation.deltachain.ibc_composability_mw.v1.GenesisState")
	proto.RegisterMapType((map[string][]byte)(nil), "wormhole_foundation.deltachain.ibc_composability_mw.v1.GenesisState.TransposedDataInFlightEntry")
}

func init() {
	proto.RegisterFile("ibc-composability-mw/genesis.proto", fileDescriptor_17a5790dbeee9b79)
}

var fileDescriptor_17a5790dbeee9b79 = []byte{
	// 330 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x91, 0xcd, 0x4a, 0xf3, 0x40,
	0x14, 0x86, 0x33, 0x2d, 0xdf, 0x07, 0xc6, 0x2e, 0x24, 0x14, 0xa9, 0x15, 0xd2, 0x90, 0x55, 0x36,
	0x9d, 0xa0, 0x82, 0x48, 0x17, 0x2e, 0x8a, 0x3f, 0x74, 0x5b, 0x5d, 0x88, 0x9b, 0x38, 0x49, 0xa6,
	0xc9, 0x60, 0x32, 0x27, 0x64, 0x4e, 0x1b, 0x73, 0x17, 0xde, 0x82, 0x17, 0x23, 0x74, 0xd9, 0xa5,
	0xab, 0x22, 0xed, 0x1d, 0x78, 0x05, 0xd2, 0x14, 0xfc, 0x81, 0xea, 0xc2, 0xdd, 0x3b, 0x9b, 0xe7,
	0x79, 0xe7, 0xbc, 0xba, 0x2d, 0xfc, 0xa0, 0x1b, 0x40, 0x9a, 0x81, 0x62, 0xbe, 0x48, 0x04, 0x96,
	0xdd, 0xb4, 0x70, 0x23, 0x2e, 0xb9, 0x12, 0x8a, 0x66, 0x39, 0x20, 0x18, 0xc7, 0x05, 0xe4, 0x69,
	0x0c, 0x09, 0xf7, 0x46, 0x30, 0x96, 0x21, 0x43, 0x01, 0x92, 0x86, 0x3c, 0x41, 0x16, 0xc4, 0x4c,
	0x48, 0x2a, 0xfc, 0xc0, 0xfb, 0x86, 0xf0, 0xd2, 0x82, 0x4e, 0x0e, 0xda, 0xcd, 0x08, 0x22, 0xa8,
	0x10, 0xee, 0x2a, 0xad, 0x69, 0xf6, 0x53, 0x4d, 0x6f, 0x5c, 0xae, 0xf9, 0x57, 0xc8, 0x90, 0x1b,
	0xcf, 0x44, 0xdf, 0xc3, 0x9c, 0x49, 0x95, 0x81, 0xe2, 0xa1, 0x17, 0x32, 0x64, 0x9e, 0x90, 0xde,
	0x28, 0x11, 0x51, 0x8c, 0x2d, 0x62, 0xd5, 0x9d, 0xed, 0xc3, 0x3b, 0xfa, 0xb7, 0x0e, 0xf4, 0xab,
	0x89, 0x5e, 0x7f, 0x58, 0xce, 0x18, 0xb2, 0x81, 0xbc, 0xa8, 0x14, 0xe7, 0x12, 0xf3, 0xb2, 0xef,
	0x4c, 0xe7, 0x1d, 0xed, 0x6d, 0xde, 0xb1, 0x4a, 0x96, 0x26, 0x3d, 0xfb, 0xc7, 0x42, 0xf6, 0x70,
	0x17, 0x37, 0x62, 0xda, 0x03, 0x7d, 0xff, 0x17, 0x81, 0xb1, 0xa3, 0xd7, 0xef, 0x79, 0xd9, 0x22,
	0x16, 0x71, 0xb6, 0x86, 0xab, 0x68, 0x34, 0xf5, 0x7f, 0x13, 0x96, 0x8c, 0x79, 0xab, 0x66, 0x11,
	0xa7, 0x31, 0x5c, 0x3f, 0x7a, 0xb5, 0x13, 0xd2, 0xbf, 0x99, 0x2e, 0x4c, 0x32, 0x5b, 0x98, 0xe4,
	0x75, 0x61, 0x92, 0xc7, 0xa5, 0xa9, 0xcd, 0x96, 0xa6, 0xf6, 0xb2, 0x34, 0xb5, 0xdb, 0xd3, 0x48,
	0x60, 0x3c, 0xf6, 0x69, 0x00, 0xa9, 0x5b, 0x7d, 0x5f, 0x15, 0x2c, 0x13, 0xe0, 0x7e, 0x9e, 0xc2,
	0x7d, 0x70, 0x37, 0x6e, 0x8a, 0x65, 0xc6, 0x95, 0xff, 0xbf, 0x1a, 0xe1, 0xe8, 0x3d, 0x00, 0x00,
	0xff, 0xff, 0x3f, 0xe8, 0x0d, 0xd5, 0xf8, 0x01, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TransposedDataInFlight) > 0 {
		for k := range m.TransposedDataInFlight {
			v := m.TransposedDataInFlight[k]
			baseI := i
			if len(v) > 0 {
				i -= len(v)
				copy(dAtA[i:], v)
				i = encodeVarintGenesis(dAtA, i, uint64(len(v)))
				i--
				dAtA[i] = 0x12
			}
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintGenesis(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintGenesis(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.TransposedDataInFlight) > 0 {
		for k, v := range m.TransposedDataInFlight {
			_ = k
			_ = v
			l = 0
			if len(v) > 0 {
				l = 1 + len(v) + sovGenesis(uint64(len(v)))
			}
			mapEntrySize := 1 + len(k) + sovGenesis(uint64(len(k))) + l
			n += mapEntrySize + 1 + sovGenesis(uint64(mapEntrySize))
		}
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TransposedDataInFlight", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.TransposedDataInFlight == nil {
				m.TransposedDataInFlight = make(map[string][]byte)
			}
			var mapkey string
			mapvalue := []byte{}
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGenesis
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGenesis
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthGenesis
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthGenesis
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var mapbyteLen uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGenesis
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapbyteLen |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intMapbyteLen := int(mapbyteLen)
					if intMapbyteLen < 0 {
						return ErrInvalidLengthGenesis
					}
					postbytesIndex := iNdEx + intMapbyteLen
					if postbytesIndex < 0 {
						return ErrInvalidLengthGenesis
					}
					if postbytesIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = make([]byte, mapbyteLen)
					copy(mapvalue, dAtA[iNdEx:postbytesIndex])
					iNdEx = postbytesIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipGenesis(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthGenesis
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.TransposedDataInFlight[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
