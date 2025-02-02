// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: interchain_security/ccv/consumer/v1/genesis.proto

package types

import (
	fmt "fmt"
	types2 "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	types "github.com/cosmos/ibc-go/v3/modules/light-clients/07-tendermint/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	types1 "github.com/tendermint/tendermint/abci/types"
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

// GenesisState defines the CCV consumer chain genesis state
type GenesisState struct {
	Params            Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	ProviderClientId  string `protobuf:"bytes,2,opt,name=provider_client_id,json=providerClientId,proto3" json:"provider_client_id,omitempty"`
	ProviderChannelId string `protobuf:"bytes,3,opt,name=provider_channel_id,json=providerChannelId,proto3" json:"provider_channel_id,omitempty"`
	NewChain          bool   `protobuf:"varint,4,opt,name=new_chain,json=newChain,proto3" json:"new_chain,omitempty"`
	// ProviderClientState filled in on new chain, nil on restart.
	ProviderClientState *types.ClientState `protobuf:"bytes,5,opt,name=provider_client_state,json=providerClientState,proto3" json:"provider_client_state,omitempty"`
	// ProviderConsensusState filled in on new chain, nil on restart.
	ProviderConsensusState *types.ConsensusState    `protobuf:"bytes,6,opt,name=provider_consensus_state,json=providerConsensusState,proto3" json:"provider_consensus_state,omitempty"`
	UnbondingSequences     []UnbondingSequence      `protobuf:"bytes,7,rep,name=unbonding_sequences,json=unbondingSequences,proto3" json:"unbonding_sequences"`
	InitialValSet          []types1.ValidatorUpdate `protobuf:"bytes,8,rep,name=initial_val_set,json=initialValSet,proto3" json:"initial_val_set"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_2db73a6057a27482, []int{0}
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

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetProviderClientId() string {
	if m != nil {
		return m.ProviderClientId
	}
	return ""
}

func (m *GenesisState) GetProviderChannelId() string {
	if m != nil {
		return m.ProviderChannelId
	}
	return ""
}

func (m *GenesisState) GetNewChain() bool {
	if m != nil {
		return m.NewChain
	}
	return false
}

func (m *GenesisState) GetProviderClientState() *types.ClientState {
	if m != nil {
		return m.ProviderClientState
	}
	return nil
}

func (m *GenesisState) GetProviderConsensusState() *types.ConsensusState {
	if m != nil {
		return m.ProviderConsensusState
	}
	return nil
}

func (m *GenesisState) GetUnbondingSequences() []UnbondingSequence {
	if m != nil {
		return m.UnbondingSequences
	}
	return nil
}

func (m *GenesisState) GetInitialValSet() []types1.ValidatorUpdate {
	if m != nil {
		return m.InitialValSet
	}
	return nil
}

// UnbondingSequence defines the genesis information for each unbonding packet sequence.
type UnbondingSequence struct {
	Sequence        uint64        `protobuf:"varint,1,opt,name=sequence,proto3" json:"sequence,omitempty"`
	UnbondingTime   uint64        `protobuf:"varint,2,opt,name=unbonding_time,json=unbondingTime,proto3" json:"unbonding_time,omitempty"`
	UnbondingPacket types2.Packet `protobuf:"bytes,3,opt,name=unbonding_packet,json=unbondingPacket,proto3" json:"unbonding_packet"`
}

func (m *UnbondingSequence) Reset()         { *m = UnbondingSequence{} }
func (m *UnbondingSequence) String() string { return proto.CompactTextString(m) }
func (*UnbondingSequence) ProtoMessage()    {}
func (*UnbondingSequence) Descriptor() ([]byte, []int) {
	return fileDescriptor_2db73a6057a27482, []int{1}
}
func (m *UnbondingSequence) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UnbondingSequence) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UnbondingSequence.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UnbondingSequence) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnbondingSequence.Merge(m, src)
}
func (m *UnbondingSequence) XXX_Size() int {
	return m.Size()
}
func (m *UnbondingSequence) XXX_DiscardUnknown() {
	xxx_messageInfo_UnbondingSequence.DiscardUnknown(m)
}

var xxx_messageInfo_UnbondingSequence proto.InternalMessageInfo

func (m *UnbondingSequence) GetSequence() uint64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *UnbondingSequence) GetUnbondingTime() uint64 {
	if m != nil {
		return m.UnbondingTime
	}
	return 0
}

func (m *UnbondingSequence) GetUnbondingPacket() types2.Packet {
	if m != nil {
		return m.UnbondingPacket
	}
	return types2.Packet{}
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "interchain_security.ccv.consumer.v1.GenesisState")
	proto.RegisterType((*UnbondingSequence)(nil), "interchain_security.ccv.consumer.v1.UnbondingSequence")
}

func init() {
	proto.RegisterFile("interchain_security/ccv/consumer/v1/genesis.proto", fileDescriptor_2db73a6057a27482)
}

var fileDescriptor_2db73a6057a27482 = []byte{
	// 570 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0xcb, 0x6e, 0xd4, 0x3c,
	0x14, 0x9e, 0xfc, 0x9d, 0xbf, 0x4c, 0x5d, 0x4a, 0x5b, 0x17, 0x50, 0xd4, 0x4a, 0x61, 0x28, 0x42,
	0x1a, 0x09, 0xb0, 0x35, 0x83, 0xc4, 0x82, 0x25, 0x5d, 0xa0, 0x4a, 0x08, 0xa1, 0xf4, 0xb2, 0x60,
	0x13, 0x39, 0xce, 0x51, 0xc6, 0x22, 0xb1, 0x43, 0xec, 0xa4, 0xf4, 0x2d, 0x78, 0x08, 0x1e, 0xa6,
	0xcb, 0x2e, 0x59, 0x55, 0xa8, 0x7d, 0x11, 0x14, 0xe7, 0x36, 0x2d, 0x48, 0xcc, 0xce, 0x3e, 0xe7,
	0xfb, 0xce, 0x77, 0xae, 0x68, 0x2a, 0xa4, 0x81, 0x9c, 0xcf, 0x99, 0x90, 0x81, 0x06, 0x5e, 0xe4,
	0xc2, 0x9c, 0x53, 0xce, 0x4b, 0xca, 0x95, 0xd4, 0x45, 0x0a, 0x39, 0x2d, 0xa7, 0x34, 0x06, 0x09,
	0x5a, 0x68, 0x92, 0xe5, 0xca, 0x28, 0xfc, 0xec, 0x2f, 0x14, 0xc2, 0x79, 0x49, 0x5a, 0x0a, 0x29,
	0xa7, 0xbb, 0x54, 0x84, 0x9c, 0x26, 0x22, 0x9e, 0x1b, 0x9e, 0x08, 0x90, 0x46, 0x53, 0x03, 0x32,
	0x82, 0x3c, 0x15, 0xd2, 0x54, 0x21, 0xfb, 0x5f, 0x1d, 0x75, 0xf7, 0x69, 0x45, 0xe0, 0x2a, 0x07,
	0xca, 0xe7, 0x4c, 0x4a, 0x48, 0x2a, 0x54, 0xf3, 0x6c, 0x20, 0x0f, 0x63, 0x15, 0x2b, 0xfb, 0xa4,
	0xd5, 0xab, 0xb1, 0xce, 0x96, 0xa9, 0xa0, 0x4b, 0xad, 0xe6, 0xec, 0x2d, 0x24, 0xc3, 0x42, 0x2e,
	0xa8, 0x39, 0xcf, 0xa0, 0xa9, 0x6f, 0xff, 0x6a, 0x88, 0xee, 0xbf, 0xaf, 0x2b, 0x3e, 0x32, 0xcc,
	0x00, 0x3e, 0x44, 0xab, 0x19, 0xcb, 0x59, 0xaa, 0x5d, 0x67, 0xec, 0x4c, 0xd6, 0x67, 0x2f, 0xc8,
	0x12, 0x1d, 0x20, 0x9f, 0x2c, 0xe5, 0xdd, 0xf0, 0xe2, 0xea, 0xc9, 0xc0, 0x6f, 0x02, 0xe0, 0x97,
	0x08, 0x67, 0xb9, 0x2a, 0x45, 0x04, 0x79, 0x50, 0x37, 0x26, 0x10, 0x91, 0xfb, 0xdf, 0xd8, 0x99,
	0xac, 0xf9, 0x5b, 0xad, 0xe7, 0xc0, 0x3a, 0x0e, 0x23, 0x4c, 0xd0, 0x4e, 0x8f, 0xae, 0x5b, 0x51,
	0xc1, 0x57, 0x2c, 0x7c, 0xbb, 0x83, 0xd7, 0x9e, 0xc3, 0x08, 0xef, 0xa1, 0x35, 0x09, 0x67, 0x81,
	0x4d, 0xcc, 0x1d, 0x8e, 0x9d, 0xc9, 0xc8, 0x1f, 0x49, 0x38, 0x3b, 0xa8, 0xfe, 0x38, 0x40, 0x8f,
	0xee, 0x4a, 0xeb, 0xaa, 0x3c, 0xf7, 0xff, 0xb6, 0xa8, 0x90, 0x93, 0xc5, 0x89, 0x91, 0x85, 0x19,
	0x95, 0x53, 0x52, 0x67, 0x65, 0x3b, 0xe2, 0xef, 0xdc, 0x4e, 0xb5, 0x6e, 0xd3, 0x1c, 0xb9, 0xbd,
	0x80, 0x92, 0x1a, 0xa4, 0x2e, 0x74, 0xa3, 0xb1, 0x6a, 0x35, 0xc8, 0x3f, 0x35, 0x5a, 0x5a, 0x2d,
	0xf3, 0xb8, 0x93, 0xb9, 0x65, 0xc7, 0x29, 0xda, 0x29, 0x64, 0xa8, 0x64, 0x24, 0x64, 0x1c, 0x68,
	0xf8, 0x5a, 0x80, 0xe4, 0xa0, 0xdd, 0x7b, 0xe3, 0x95, 0xc9, 0xfa, 0xec, 0xcd, 0x52, 0xd3, 0x39,
	0x69, 0xf9, 0x47, 0x0d, 0xbd, 0x19, 0x14, 0x2e, 0xee, 0x3a, 0x34, 0xfe, 0x88, 0x36, 0x85, 0x14,
	0x46, 0xb0, 0x24, 0x28, 0x59, 0x12, 0x68, 0x30, 0xee, 0xc8, 0x4a, 0x8d, 0x17, 0xd3, 0xaf, 0xf6,
	0x88, 0x9c, 0xb2, 0x44, 0x44, 0xcc, 0xa8, 0xfc, 0x24, 0x8b, 0x98, 0x69, 0x83, 0x6e, 0x34, 0xf4,
	0x53, 0x96, 0x1c, 0x81, 0xd9, 0xff, 0xe1, 0xa0, 0xed, 0x3f, 0xf4, 0xf1, 0x2e, 0x1a, 0xb5, 0xa5,
	0xd8, 0x3d, 0x1b, 0xfa, 0xdd, 0x1f, 0x3f, 0x47, 0x0f, 0xfa, 0x82, 0x8d, 0x48, 0xc1, 0xae, 0xcc,
	0xd0, 0xdf, 0xe8, 0xac, 0xc7, 0x22, 0x05, 0xfc, 0x01, 0x6d, 0xf5, 0xb0, 0x8c, 0xf1, 0x2f, 0x60,
	0xec, 0xb2, 0xac, 0xcf, 0xf6, 0x6c, 0xe7, 0xab, 0xf3, 0x22, 0xed, 0x4d, 0xd9, 0x15, 0xad, 0x20,
	0x4d, 0x92, 0x9b, 0x1d, 0xb5, 0x31, 0x1f, 0x5f, 0x5c, 0x7b, 0xce, 0xe5, 0xb5, 0xe7, 0xfc, 0xba,
	0xf6, 0x9c, 0xef, 0x37, 0xde, 0xe0, 0xf2, 0xc6, 0x1b, 0xfc, 0xbc, 0xf1, 0x06, 0x9f, 0xdf, 0xc6,
	0xc2, 0xcc, 0x8b, 0x90, 0x70, 0x95, 0x52, 0xae, 0x74, 0xaa, 0x34, 0xed, 0x7b, 0xfe, 0xaa, 0x3b,
	0xc2, 0x6f, 0xb7, 0xcf, 0xd0, 0xde, 0x58, 0xb8, 0x6a, 0x8f, 0xec, 0xf5, 0xef, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x3c, 0x93, 0x99, 0x5b, 0x79, 0x04, 0x00, 0x00,
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
	if len(m.InitialValSet) > 0 {
		for iNdEx := len(m.InitialValSet) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.InitialValSet[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x42
		}
	}
	if len(m.UnbondingSequences) > 0 {
		for iNdEx := len(m.UnbondingSequences) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.UnbondingSequences[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if m.ProviderConsensusState != nil {
		{
			size, err := m.ProviderConsensusState.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x32
	}
	if m.ProviderClientState != nil {
		{
			size, err := m.ProviderClientState.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.NewChain {
		i--
		if m.NewChain {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x20
	}
	if len(m.ProviderChannelId) > 0 {
		i -= len(m.ProviderChannelId)
		copy(dAtA[i:], m.ProviderChannelId)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.ProviderChannelId)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.ProviderClientId) > 0 {
		i -= len(m.ProviderClientId)
		copy(dAtA[i:], m.ProviderClientId)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.ProviderClientId)))
		i--
		dAtA[i] = 0x12
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *UnbondingSequence) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UnbondingSequence) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UnbondingSequence) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.UnbondingPacket.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if m.UnbondingTime != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.UnbondingTime))
		i--
		dAtA[i] = 0x10
	}
	if m.Sequence != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.Sequence))
		i--
		dAtA[i] = 0x8
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
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = len(m.ProviderClientId)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = len(m.ProviderChannelId)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.NewChain {
		n += 2
	}
	if m.ProviderClientState != nil {
		l = m.ProviderClientState.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.ProviderConsensusState != nil {
		l = m.ProviderConsensusState.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	if len(m.UnbondingSequences) > 0 {
		for _, e := range m.UnbondingSequences {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.InitialValSet) > 0 {
		for _, e := range m.InitialValSet {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func (m *UnbondingSequence) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Sequence != 0 {
		n += 1 + sovGenesis(uint64(m.Sequence))
	}
	if m.UnbondingTime != 0 {
		n += 1 + sovGenesis(uint64(m.UnbondingTime))
	}
	l = m.UnbondingPacket.Size()
	n += 1 + l + sovGenesis(uint64(l))
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
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
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
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProviderClientId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProviderClientId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProviderChannelId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProviderChannelId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NewChain", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.NewChain = bool(v != 0)
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProviderClientState", wireType)
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
			if m.ProviderClientState == nil {
				m.ProviderClientState = &types.ClientState{}
			}
			if err := m.ProviderClientState.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProviderConsensusState", wireType)
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
			if m.ProviderConsensusState == nil {
				m.ProviderConsensusState = &types.ConsensusState{}
			}
			if err := m.ProviderConsensusState.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnbondingSequences", wireType)
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
			m.UnbondingSequences = append(m.UnbondingSequences, UnbondingSequence{})
			if err := m.UnbondingSequences[len(m.UnbondingSequences)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InitialValSet", wireType)
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
			m.InitialValSet = append(m.InitialValSet, types1.ValidatorUpdate{})
			if err := m.InitialValSet[len(m.InitialValSet)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
func (m *UnbondingSequence) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: UnbondingSequence: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UnbondingSequence: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sequence", wireType)
			}
			m.Sequence = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Sequence |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnbondingTime", wireType)
			}
			m.UnbondingTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UnbondingTime |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnbondingPacket", wireType)
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
			if err := m.UnbondingPacket.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
