// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: deltaswap/governance.proto

package types

import (
	bytes "bytes"
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

// PhylaxSetUpdateProposal defines a phylax set update governance proposal
type PhylaxSetUpdateProposal struct {
	Title        string    `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description  string    `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	NewPhylaxSet PhylaxSet `protobuf:"bytes,3,opt,name=newPhylaxSet,proto3" json:"newPhylaxSet"`
}

func (m *PhylaxSetUpdateProposal) Reset()      { *m = PhylaxSetUpdateProposal{} }
func (*PhylaxSetUpdateProposal) ProtoMessage() {}
func (*PhylaxSetUpdateProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_69d0ecfe102360ce, []int{0}
}
func (m *PhylaxSetUpdateProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PhylaxSetUpdateProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PhylaxSetUpdateProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PhylaxSetUpdateProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PhylaxSetUpdateProposal.Merge(m, src)
}
func (m *PhylaxSetUpdateProposal) XXX_Size() int {
	return m.Size()
}
func (m *PhylaxSetUpdateProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_PhylaxSetUpdateProposal.DiscardUnknown(m)
}

var xxx_messageInfo_PhylaxSetUpdateProposal proto.InternalMessageInfo

func (m *PhylaxSetUpdateProposal) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *PhylaxSetUpdateProposal) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *PhylaxSetUpdateProposal) GetNewPhylaxSet() PhylaxSet {
	if m != nil {
		return m.NewPhylaxSet
	}
	return PhylaxSet{}
}

// GovernanceDeltaswapMessageProposal defines a governance proposal to emit a generic message in the governance message
// format.
type GovernanceDeltaswapMessageProposal struct {
	Title       string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Action      uint32 `protobuf:"varint,3,opt,name=action,proto3" json:"action,omitempty"`
	Module      []byte `protobuf:"bytes,4,opt,name=module,proto3" json:"module,omitempty"`
	TargetChain uint32 `protobuf:"varint,5,opt,name=targetChain,proto3" json:"targetChain,omitempty"`
	Payload     []byte `protobuf:"bytes,6,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *GovernanceDeltaswapMessageProposal) Reset()      { *m = GovernanceDeltaswapMessageProposal{} }
func (*GovernanceDeltaswapMessageProposal) ProtoMessage() {}
func (*GovernanceDeltaswapMessageProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_69d0ecfe102360ce, []int{1}
}
func (m *GovernanceDeltaswapMessageProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GovernanceDeltaswapMessageProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GovernanceDeltaswapMessageProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GovernanceDeltaswapMessageProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GovernanceDeltaswapMessageProposal.Merge(m, src)
}
func (m *GovernanceDeltaswapMessageProposal) XXX_Size() int {
	return m.Size()
}
func (m *GovernanceDeltaswapMessageProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_GovernanceDeltaswapMessageProposal.DiscardUnknown(m)
}

var xxx_messageInfo_GovernanceDeltaswapMessageProposal proto.InternalMessageInfo

func (m *GovernanceDeltaswapMessageProposal) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *GovernanceDeltaswapMessageProposal) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *GovernanceDeltaswapMessageProposal) GetAction() uint32 {
	if m != nil {
		return m.Action
	}
	return 0
}

func (m *GovernanceDeltaswapMessageProposal) GetModule() []byte {
	if m != nil {
		return m.Module
	}
	return nil
}

func (m *GovernanceDeltaswapMessageProposal) GetTargetChain() uint32 {
	if m != nil {
		return m.TargetChain
	}
	return 0
}

func (m *GovernanceDeltaswapMessageProposal) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func init() {
	proto.RegisterType((*PhylaxSetUpdateProposal)(nil), "deltaswapio.deltachain.deltaswap.PhylaxSetUpdateProposal")
	proto.RegisterType((*GovernanceDeltaswapMessageProposal)(nil), "deltaswapio.deltachain.deltaswap.GovernanceDeltaswapMessageProposal")
}

func init() { proto.RegisterFile("deltaswap/governance.proto", fileDescriptor_69d0ecfe102360ce) }

var fileDescriptor_69d0ecfe102360ce = []byte{
	// 343 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x91, 0x31, 0x4f, 0xc2, 0x40,
	0x14, 0xc7, 0x7b, 0x0a, 0xa8, 0x07, 0x2e, 0x0d, 0xc1, 0x0b, 0x43, 0x69, 0x98, 0x48, 0x4c, 0xda,
	0x44, 0x36, 0x47, 0x34, 0x71, 0x22, 0x21, 0x35, 0x2c, 0x6e, 0x47, 0xfb, 0x52, 0x9a, 0x94, 0xbe,
	0x4b, 0xef, 0x10, 0xf8, 0x16, 0x8e, 0x8e, 0x7c, 0x06, 0x3f, 0x05, 0x23, 0x71, 0x72, 0x32, 0x06,
	0x16, 0x3f, 0x86, 0xe9, 0x01, 0x2d, 0x4e, 0x0e, 0x6e, 0xef, 0xff, 0xff, 0xdf, 0xef, 0xe5, 0xbd,
	0x77, 0xb4, 0x19, 0x40, 0xac, 0xb8, 0x9c, 0x71, 0xe1, 0x86, 0xf8, 0x0c, 0x69, 0xc2, 0x13, 0x1f,
	0x1c, 0x91, 0xa2, 0x42, 0xd3, 0xce, 0xb3, 0x08, 0x1d, 0x5d, 0xfb, 0x63, 0x1e, 0x25, 0x4e, 0x6e,
	0x37, 0x1b, 0x05, 0x2d, 0xc6, 0x8b, 0x98, 0xcf, 0x77, 0x64, 0xb3, 0x1e, 0x62, 0x88, 0xba, 0x74,
	0xb3, 0x6a, 0xe7, 0xb6, 0xdf, 0x08, 0xbd, 0x1a, 0xe8, 0x67, 0x8f, 0xa0, 0x86, 0x22, 0xe0, 0x0a,
	0x06, 0x29, 0x0a, 0x94, 0x3c, 0x36, 0xeb, 0xb4, 0xac, 0x22, 0x15, 0x03, 0x23, 0x36, 0xe9, 0x5c,
	0x78, 0x3b, 0x61, 0xda, 0xb4, 0x1a, 0x80, 0xf4, 0xd3, 0x48, 0xa8, 0x08, 0x13, 0x76, 0xa2, 0xb3,
	0x63, 0xcb, 0x1c, 0xd2, 0x5a, 0x02, 0xb3, 0xbc, 0x2b, 0x3b, 0xb5, 0x49, 0xa7, 0x7a, 0x73, 0xed,
	0xfc, 0x35, 0xba, 0x93, 0x23, 0xbd, 0xd2, 0xea, 0xb3, 0x65, 0x78, 0xbf, 0xda, 0xdc, 0x9e, 0xbf,
	0x2e, 0x5b, 0xc6, 0xf7, 0xb2, 0x45, 0xda, 0xef, 0x84, 0xb6, 0x1f, 0xf2, 0xcb, 0xdc, 0x1f, 0xf8,
	0x3e, 0x48, 0xc9, 0xc3, 0xff, 0xcf, 0xdf, 0xa0, 0x15, 0xee, 0xeb, 0x30, 0x9b, 0xfc, 0xd2, 0xdb,
	0xab, 0xcc, 0x9f, 0x60, 0x30, 0x8d, 0x81, 0x95, 0x6c, 0xd2, 0xa9, 0x79, 0x7b, 0x95, 0x75, 0x54,
	0x3c, 0x0d, 0x41, 0xdd, 0x65, 0xfb, 0xb0, 0xb2, 0x86, 0x8e, 0x2d, 0x93, 0xd1, 0x33, 0xc1, 0x17,
	0x31, 0xf2, 0x80, 0x55, 0x34, 0x7a, 0x90, 0xc5, 0x52, 0xbd, 0xfe, 0x6a, 0x63, 0x91, 0xf5, 0xc6,
	0x22, 0x5f, 0x1b, 0x8b, 0xbc, 0x6c, 0x2d, 0x63, 0xbd, 0xb5, 0x8c, 0x8f, 0xad, 0x65, 0x3c, 0x75,
	0xc3, 0x48, 0x8d, 0xa7, 0x23, 0xc7, 0xc7, 0x89, 0x7b, 0x74, 0x43, 0xb7, 0xb8, 0xa1, 0x3b, 0x2f,
	0x02, 0x57, 0x2d, 0x04, 0xc8, 0x51, 0x45, 0xff, 0x6f, 0xf7, 0x27, 0x00, 0x00, 0xff, 0xff, 0xce,
	0xe2, 0x86, 0x20, 0x4d, 0x02, 0x00, 0x00,
}

func (this *PhylaxSetUpdateProposal) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*PhylaxSetUpdateProposal)
	if !ok {
		that2, ok := that.(PhylaxSetUpdateProposal)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Title != that1.Title {
		return false
	}
	if this.Description != that1.Description {
		return false
	}
	if !this.NewPhylaxSet.Equal(&that1.NewPhylaxSet) {
		return false
	}
	return true
}
func (this *GovernanceDeltaswapMessageProposal) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*GovernanceDeltaswapMessageProposal)
	if !ok {
		that2, ok := that.(GovernanceDeltaswapMessageProposal)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Title != that1.Title {
		return false
	}
	if this.Description != that1.Description {
		return false
	}
	if this.Action != that1.Action {
		return false
	}
	if !bytes.Equal(this.Module, that1.Module) {
		return false
	}
	if this.TargetChain != that1.TargetChain {
		return false
	}
	if !bytes.Equal(this.Payload, that1.Payload) {
		return false
	}
	return true
}
func (m *PhylaxSetUpdateProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PhylaxSetUpdateProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PhylaxSetUpdateProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.NewPhylaxSet.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGovernance(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintGovernance(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintGovernance(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *GovernanceDeltaswapMessageProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GovernanceDeltaswapMessageProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GovernanceDeltaswapMessageProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Payload) > 0 {
		i -= len(m.Payload)
		copy(dAtA[i:], m.Payload)
		i = encodeVarintGovernance(dAtA, i, uint64(len(m.Payload)))
		i--
		dAtA[i] = 0x32
	}
	if m.TargetChain != 0 {
		i = encodeVarintGovernance(dAtA, i, uint64(m.TargetChain))
		i--
		dAtA[i] = 0x28
	}
	if len(m.Module) > 0 {
		i -= len(m.Module)
		copy(dAtA[i:], m.Module)
		i = encodeVarintGovernance(dAtA, i, uint64(len(m.Module)))
		i--
		dAtA[i] = 0x22
	}
	if m.Action != 0 {
		i = encodeVarintGovernance(dAtA, i, uint64(m.Action))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintGovernance(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintGovernance(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGovernance(dAtA []byte, offset int, v uint64) int {
	offset -= sovGovernance(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *PhylaxSetUpdateProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovGovernance(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovGovernance(uint64(l))
	}
	l = m.NewPhylaxSet.Size()
	n += 1 + l + sovGovernance(uint64(l))
	return n
}

func (m *GovernanceDeltaswapMessageProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovGovernance(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovGovernance(uint64(l))
	}
	if m.Action != 0 {
		n += 1 + sovGovernance(uint64(m.Action))
	}
	l = len(m.Module)
	if l > 0 {
		n += 1 + l + sovGovernance(uint64(l))
	}
	if m.TargetChain != 0 {
		n += 1 + sovGovernance(uint64(m.TargetChain))
	}
	l = len(m.Payload)
	if l > 0 {
		n += 1 + l + sovGovernance(uint64(l))
	}
	return n
}

func sovGovernance(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGovernance(x uint64) (n int) {
	return sovGovernance(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PhylaxSetUpdateProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGovernance
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
			return fmt.Errorf("proto: PhylaxSetUpdateProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PhylaxSetUpdateProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernance
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGovernance
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGovernance
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernance
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGovernance
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGovernance
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NewPhylaxSet", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernance
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
				return ErrInvalidLengthGovernance
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGovernance
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.NewPhylaxSet.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGovernance(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGovernance
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
func (m *GovernanceDeltaswapMessageProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGovernance
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
			return fmt.Errorf("proto: GovernanceDeltaswapMessageProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GovernanceDeltaswapMessageProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernance
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGovernance
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGovernance
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernance
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGovernance
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGovernance
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Action", wireType)
			}
			m.Action = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernance
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Action |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Module", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernance
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthGovernance
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGovernance
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Module = append(m.Module[:0], dAtA[iNdEx:postIndex]...)
			if m.Module == nil {
				m.Module = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TargetChain", wireType)
			}
			m.TargetChain = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernance
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TargetChain |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Payload", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGovernance
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthGovernance
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGovernance
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Payload = append(m.Payload[:0], dAtA[iNdEx:postIndex]...)
			if m.Payload == nil {
				m.Payload = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGovernance(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGovernance
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
func skipGovernance(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGovernance
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
					return 0, ErrIntOverflowGovernance
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
					return 0, ErrIntOverflowGovernance
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
				return 0, ErrInvalidLengthGovernance
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGovernance
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGovernance
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGovernance        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGovernance          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGovernance = fmt.Errorf("proto: unexpected end of group")
)
