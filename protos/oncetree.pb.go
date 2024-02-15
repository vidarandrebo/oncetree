// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.2
// source: protos/oncetree.proto

package protos

import (
	_ "github.com/relab/gorums"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ReadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key int64 `protobuf:"varint,1,opt,name=Key,proto3" json:"Key,omitempty"`
}

func (x *ReadRequest) Reset() {
	*x = ReadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_oncetree_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadRequest) ProtoMessage() {}

func (x *ReadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_oncetree_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadRequest.ProtoReflect.Descriptor instead.
func (*ReadRequest) Descriptor() ([]byte, []int) {
	return file_protos_oncetree_proto_rawDescGZIP(), []int{0}
}

func (x *ReadRequest) GetKey() int64 {
	if x != nil {
		return x.Key
	}
	return 0
}

type ReadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value int64 `protobuf:"varint,1,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *ReadResponse) Reset() {
	*x = ReadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_oncetree_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadResponse) ProtoMessage() {}

func (x *ReadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_oncetree_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadResponse.ProtoReflect.Descriptor instead.
func (*ReadResponse) Descriptor() ([]byte, []int) {
	return file_protos_oncetree_proto_rawDescGZIP(), []int{1}
}

func (x *ReadResponse) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

// ReadAllResponse contains a map Address -> Value
type ReadAllResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value map[string]int64 `protobuf:"bytes,1,rep,name=Value,proto3" json:"Value,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *ReadAllResponse) Reset() {
	*x = ReadAllResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_oncetree_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadAllResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadAllResponse) ProtoMessage() {}

func (x *ReadAllResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_oncetree_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadAllResponse.ProtoReflect.Descriptor instead.
func (*ReadAllResponse) Descriptor() ([]byte, []int) {
	return file_protos_oncetree_proto_rawDescGZIP(), []int{2}
}

func (x *ReadAllResponse) GetValue() map[string]int64 {
	if x != nil {
		return x.Value
	}
	return nil
}

type WriteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   int64 `protobuf:"varint,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value int64 `protobuf:"varint,2,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *WriteRequest) Reset() {
	*x = WriteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_oncetree_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WriteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriteRequest) ProtoMessage() {}

func (x *WriteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_oncetree_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WriteRequest.ProtoReflect.Descriptor instead.
func (*WriteRequest) Descriptor() ([]byte, []int) {
	return file_protos_oncetree_proto_rawDescGZIP(), []int{3}
}

func (x *WriteRequest) GetKey() int64 {
	if x != nil {
		return x.Key
	}
	return 0
}

func (x *WriteRequest) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type GroupInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeID           string `protobuf:"bytes,1,opt,name=NodeID,proto3" json:"NodeID,omitempty"`
	NeighbourID      string `protobuf:"bytes,2,opt,name=NeighbourID,proto3" json:"NeighbourID,omitempty"`
	NeighbourAddress string `protobuf:"bytes,3,opt,name=NeighbourAddress,proto3" json:"NeighbourAddress,omitempty"`
	NeighbourRole    int64  `protobuf:"varint,4,opt,name=NeighbourRole,proto3" json:"NeighbourRole,omitempty"`
}

func (x *GroupInfo) Reset() {
	*x = GroupInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_oncetree_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupInfo) ProtoMessage() {}

func (x *GroupInfo) ProtoReflect() protoreflect.Message {
	mi := &file_protos_oncetree_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupInfo.ProtoReflect.Descriptor instead.
func (*GroupInfo) Descriptor() ([]byte, []int) {
	return file_protos_oncetree_proto_rawDescGZIP(), []int{4}
}

func (x *GroupInfo) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

func (x *GroupInfo) GetNeighbourID() string {
	if x != nil {
		return x.NeighbourID
	}
	return ""
}

func (x *GroupInfo) GetNeighbourAddress() string {
	if x != nil {
		return x.NeighbourAddress
	}
	return ""
}

func (x *GroupInfo) GetNeighbourRole() int64 {
	if x != nil {
		return x.NeighbourRole
	}
	return 0
}

type GossipMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeID    string `protobuf:"bytes,1,opt,name=NodeID,proto3" json:"NodeID,omitempty"`
	Key       int64  `protobuf:"varint,2,opt,name=Key,proto3" json:"Key,omitempty"`
	Value     int64  `protobuf:"varint,3,opt,name=Value,proto3" json:"Value,omitempty"`
	Timestamp int64  `protobuf:"varint,4,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
}

func (x *GossipMessage) Reset() {
	*x = GossipMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_oncetree_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GossipMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GossipMessage) ProtoMessage() {}

func (x *GossipMessage) ProtoReflect() protoreflect.Message {
	mi := &file_protos_oncetree_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GossipMessage.ProtoReflect.Descriptor instead.
func (*GossipMessage) Descriptor() ([]byte, []int) {
	return file_protos_oncetree_proto_rawDescGZIP(), []int{5}
}

func (x *GossipMessage) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

func (x *GossipMessage) GetKey() int64 {
	if x != nil {
		return x.Key
	}
	return 0
}

func (x *GossipMessage) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *GossipMessage) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type HeartbeatMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeID string `protobuf:"bytes,1,opt,name=NodeID,proto3" json:"NodeID,omitempty"`
}

func (x *HeartbeatMessage) Reset() {
	*x = HeartbeatMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_oncetree_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeartbeatMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartbeatMessage) ProtoMessage() {}

func (x *HeartbeatMessage) ProtoReflect() protoreflect.Message {
	mi := &file_protos_oncetree_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartbeatMessage.ProtoReflect.Descriptor instead.
func (*HeartbeatMessage) Descriptor() ([]byte, []int) {
	return file_protos_oncetree_proto_rawDescGZIP(), []int{6}
}

func (x *HeartbeatMessage) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

type PrepareMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeID string `protobuf:"bytes,1,opt,name=NodeID,proto3" json:"NodeID,omitempty"`
	Hash   string `protobuf:"bytes,2,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Round  int64  `protobuf:"varint,3,opt,name=Round,proto3" json:"Round,omitempty"`
}

func (x *PrepareMessage) Reset() {
	*x = PrepareMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_oncetree_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrepareMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrepareMessage) ProtoMessage() {}

func (x *PrepareMessage) ProtoReflect() protoreflect.Message {
	mi := &file_protos_oncetree_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrepareMessage.ProtoReflect.Descriptor instead.
func (*PrepareMessage) Descriptor() ([]byte, []int) {
	return file_protos_oncetree_proto_rawDescGZIP(), []int{7}
}

func (x *PrepareMessage) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

func (x *PrepareMessage) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *PrepareMessage) GetRound() int64 {
	if x != nil {
		return x.Round
	}
	return 0
}

type PromiseMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeID string `protobuf:"bytes,1,opt,name=NodeID,proto3" json:"NodeID,omitempty"`
	Hash   string `protobuf:"bytes,2,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Round  int64  `protobuf:"varint,3,opt,name=Round,proto3" json:"Round,omitempty"`
}

func (x *PromiseMessage) Reset() {
	*x = PromiseMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_oncetree_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PromiseMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PromiseMessage) ProtoMessage() {}

func (x *PromiseMessage) ProtoReflect() protoreflect.Message {
	mi := &file_protos_oncetree_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PromiseMessage.ProtoReflect.Descriptor instead.
func (*PromiseMessage) Descriptor() ([]byte, []int) {
	return file_protos_oncetree_proto_rawDescGZIP(), []int{8}
}

func (x *PromiseMessage) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

func (x *PromiseMessage) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *PromiseMessage) GetRound() int64 {
	if x != nil {
		return x.Round
	}
	return 0
}

type AcceptMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value int64  `protobuf:"varint,2,opt,name=Value,proto3" json:"Value,omitempty"`
	Round int64  `protobuf:"varint,3,opt,name=Round,proto3" json:"Round,omitempty"`
}

func (x *AcceptMessage) Reset() {
	*x = AcceptMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_oncetree_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AcceptMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptMessage) ProtoMessage() {}

func (x *AcceptMessage) ProtoReflect() protoreflect.Message {
	mi := &file_protos_oncetree_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcceptMessage.ProtoReflect.Descriptor instead.
func (*AcceptMessage) Descriptor() ([]byte, []int) {
	return file_protos_oncetree_proto_rawDescGZIP(), []int{9}
}

func (x *AcceptMessage) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *AcceptMessage) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *AcceptMessage) GetRound() int64 {
	if x != nil {
		return x.Round
	}
	return 0
}

type LearnMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value int64  `protobuf:"varint,2,opt,name=Value,proto3" json:"Value,omitempty"`
	Round int64  `protobuf:"varint,3,opt,name=Round,proto3" json:"Round,omitempty"`
}

func (x *LearnMessage) Reset() {
	*x = LearnMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_oncetree_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LearnMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LearnMessage) ProtoMessage() {}

func (x *LearnMessage) ProtoReflect() protoreflect.Message {
	mi := &file_protos_oncetree_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LearnMessage.ProtoReflect.Descriptor instead.
func (*LearnMessage) Descriptor() ([]byte, []int) {
	return file_protos_oncetree_proto_rawDescGZIP(), []int{10}
}

func (x *LearnMessage) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *LearnMessage) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *LearnMessage) GetRound() int64 {
	if x != nil {
		return x.Round
	}
	return 0
}

var File_protos_oncetree_proto protoreflect.FileDescriptor

var file_protos_oncetree_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6f, 0x6e, 0x63, 0x65, 0x74, 0x72, 0x65,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x1a,
	0x0c, 0x67, 0x6f, 0x72, 0x75, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65,
	0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1f, 0x0a, 0x0b, 0x52, 0x65,
	0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x22, 0x24, 0x0a, 0x0c, 0x52,
	0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x22, 0x85, 0x01, 0x0a, 0x0f, 0x52, 0x65, 0x61, 0x64, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x52, 0x65,
	0x61, 0x64, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a,
	0x38, 0x0a, 0x0a, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x36, 0x0a, 0x0c, 0x57, 0x72, 0x69,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x22, 0x97, 0x01, 0x0a, 0x09, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x16, 0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x12, 0x20, 0x0a, 0x0b, 0x4e, 0x65, 0x69, 0x67, 0x68,
	0x62, 0x6f, 0x75, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x4e, 0x65,
	0x69, 0x67, 0x68, 0x62, 0x6f, 0x75, 0x72, 0x49, 0x44, 0x12, 0x2a, 0x0a, 0x10, 0x4e, 0x65, 0x69,
	0x67, 0x68, 0x62, 0x6f, 0x75, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x10, 0x4e, 0x65, 0x69, 0x67, 0x68, 0x62, 0x6f, 0x75, 0x72, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x4e, 0x65, 0x69, 0x67, 0x68, 0x62, 0x6f,
	0x75, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x4e, 0x65,
	0x69, 0x67, 0x68, 0x62, 0x6f, 0x75, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x22, 0x6d, 0x0a, 0x0d, 0x47,
	0x6f, 0x73, 0x73, 0x69, 0x70, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4e, 0x6f,
	0x64, 0x65, 0x49, 0x44, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1c, 0x0a, 0x09,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x2a, 0x0a, 0x10, 0x48, 0x65,
	0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x22, 0x52, 0x0a, 0x0e, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72,
	0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44,
	0x12, 0x12, 0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x48, 0x61, 0x73, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x05, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x22, 0x52, 0x0a, 0x0e, 0x50, 0x72,
	0x6f, 0x6d, 0x69, 0x73, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4e, 0x6f,
	0x64, 0x65, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x48, 0x61, 0x73, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x52, 0x6f, 0x75, 0x6e,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x22, 0x4d,
	0x0a, 0x0d, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x4b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x52, 0x6f, 0x75, 0x6e, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x22, 0x4c, 0x0a,
	0x0c, 0x4c, 0x65, 0x61, 0x72, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x32, 0xc1, 0x03, 0x0a, 0x0f,
	0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x37, 0x0a, 0x05, 0x57, 0x72, 0x69, 0x74, 0x65, 0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x73, 0x2e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x33, 0x0a, 0x04, 0x52, 0x65, 0x61, 0x64,
	0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x52,
	0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3d, 0x0a,
	0x07, 0x52, 0x65, 0x61, 0x64, 0x41, 0x6c, 0x6c, 0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x73, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x41, 0x6c, 0x6c, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x04, 0xa0, 0xb5, 0x18, 0x01, 0x12, 0x43, 0x0a, 0x09,
	0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x73, 0x2e, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x04, 0x98, 0xb5, 0x18,
	0x01, 0x12, 0x41, 0x0a, 0x0e, 0x53, 0x65, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x12, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x04,
	0x98, 0xb5, 0x18, 0x01, 0x12, 0x39, 0x0a, 0x06, 0x47, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x12, 0x15,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x47, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12,
	0x3e, 0x0a, 0x0a, 0x50, 0x72, 0x69, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42,
	0x29, 0x5a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x76, 0x69,
	0x64, 0x61, 0x72, 0x61, 0x6e, 0x64, 0x72, 0x65, 0x62, 0x6f, 0x2f, 0x6f, 0x6e, 0x63, 0x65, 0x74,
	0x72, 0x65, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_protos_oncetree_proto_rawDescOnce sync.Once
	file_protos_oncetree_proto_rawDescData = file_protos_oncetree_proto_rawDesc
)

func file_protos_oncetree_proto_rawDescGZIP() []byte {
	file_protos_oncetree_proto_rawDescOnce.Do(func() {
		file_protos_oncetree_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_oncetree_proto_rawDescData)
	})
	return file_protos_oncetree_proto_rawDescData
}

var file_protos_oncetree_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_protos_oncetree_proto_goTypes = []interface{}{
	(*ReadRequest)(nil),      // 0: protos.ReadRequest
	(*ReadResponse)(nil),     // 1: protos.ReadResponse
	(*ReadAllResponse)(nil),  // 2: protos.ReadAllResponse
	(*WriteRequest)(nil),     // 3: protos.WriteRequest
	(*GroupInfo)(nil),        // 4: protos.GroupInfo
	(*GossipMessage)(nil),    // 5: protos.GossipMessage
	(*HeartbeatMessage)(nil), // 6: protos.HeartbeatMessage
	(*PrepareMessage)(nil),   // 7: protos.PrepareMessage
	(*PromiseMessage)(nil),   // 8: protos.PromiseMessage
	(*AcceptMessage)(nil),    // 9: protos.AcceptMessage
	(*LearnMessage)(nil),     // 10: protos.LearnMessage
	nil,                      // 11: protos.ReadAllResponse.ValueEntry
	(*emptypb.Empty)(nil),    // 12: google.protobuf.Empty
}
var file_protos_oncetree_proto_depIdxs = []int32{
	11, // 0: protos.ReadAllResponse.Value:type_name -> protos.ReadAllResponse.ValueEntry
	3,  // 1: protos.KeyValueService.Write:input_type -> protos.WriteRequest
	0,  // 2: protos.KeyValueService.Read:input_type -> protos.ReadRequest
	0,  // 3: protos.KeyValueService.ReadAll:input_type -> protos.ReadRequest
	6,  // 4: protos.KeyValueService.Heartbeat:input_type -> protos.HeartbeatMessage
	4,  // 5: protos.KeyValueService.SetGroupMember:input_type -> protos.GroupInfo
	5,  // 6: protos.KeyValueService.Gossip:input_type -> protos.GossipMessage
	12, // 7: protos.KeyValueService.PrintState:input_type -> google.protobuf.Empty
	12, // 8: protos.KeyValueService.Write:output_type -> google.protobuf.Empty
	1,  // 9: protos.KeyValueService.Read:output_type -> protos.ReadResponse
	2,  // 10: protos.KeyValueService.ReadAll:output_type -> protos.ReadAllResponse
	12, // 11: protos.KeyValueService.Heartbeat:output_type -> google.protobuf.Empty
	12, // 12: protos.KeyValueService.SetGroupMember:output_type -> google.protobuf.Empty
	12, // 13: protos.KeyValueService.Gossip:output_type -> google.protobuf.Empty
	12, // 14: protos.KeyValueService.PrintState:output_type -> google.protobuf.Empty
	8,  // [8:15] is the sub-list for method output_type
	1,  // [1:8] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_protos_oncetree_proto_init() }
func file_protos_oncetree_proto_init() {
	if File_protos_oncetree_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_oncetree_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_oncetree_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_oncetree_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadAllResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_oncetree_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WriteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_oncetree_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GroupInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_oncetree_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GossipMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_oncetree_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HeartbeatMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_oncetree_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrepareMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_oncetree_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PromiseMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_oncetree_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AcceptMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_oncetree_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LearnMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protos_oncetree_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_oncetree_proto_goTypes,
		DependencyIndexes: file_protos_oncetree_proto_depIdxs,
		MessageInfos:      file_protos_oncetree_proto_msgTypes,
	}.Build()
	File_protos_oncetree_proto = out.File
	file_protos_oncetree_proto_rawDesc = nil
	file_protos_oncetree_proto_goTypes = nil
	file_protos_oncetree_proto_depIdxs = nil
}
