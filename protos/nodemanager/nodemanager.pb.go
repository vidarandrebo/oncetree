// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.3
// source: protos/nodemanager/nodemanager.proto

package nodemanager

import (
	reflect "reflect"
	sync "sync"

	_ "github.com/relab/gorums"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type JoinRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeID  string `protobuf:"bytes,1,opt,name=NodeID,proto3" json:"NodeID,omitempty"`
	Address string `protobuf:"bytes,2,opt,name=Address,proto3" json:"Address,omitempty"`
}

func (x *JoinRequest) Reset() {
	*x = JoinRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_nodemanager_nodemanager_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinRequest) ProtoMessage() {}

func (x *JoinRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_nodemanager_nodemanager_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinRequest.ProtoReflect.Descriptor instead.
func (*JoinRequest) Descriptor() ([]byte, []int) {
	return file_protos_nodemanager_nodemanager_proto_rawDescGZIP(), []int{0}
}

func (x *JoinRequest) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

func (x *JoinRequest) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type JoinResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OK          bool   `protobuf:"varint,1,opt,name=OK,proto3" json:"OK,omitempty"`
	NodeID      string `protobuf:"bytes,2,opt,name=NodeID,proto3" json:"NodeID,omitempty"`
	NextAddress string `protobuf:"bytes,3,opt,name=NextAddress,proto3" json:"NextAddress,omitempty"`
}

func (x *JoinResponse) Reset() {
	*x = JoinResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_nodemanager_nodemanager_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinResponse) ProtoMessage() {}

func (x *JoinResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_nodemanager_nodemanager_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinResponse.ProtoReflect.Descriptor instead.
func (*JoinResponse) Descriptor() ([]byte, []int) {
	return file_protos_nodemanager_nodemanager_proto_rawDescGZIP(), []int{1}
}

func (x *JoinResponse) GetOK() bool {
	if x != nil {
		return x.OK
	}
	return false
}

func (x *JoinResponse) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

func (x *JoinResponse) GetNextAddress() string {
	if x != nil {
		return x.NextAddress
	}
	return ""
}

type ReadyMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeID string `protobuf:"bytes,1,opt,name=NodeID,proto3" json:"NodeID,omitempty"`
}

func (x *ReadyMessage) Reset() {
	*x = ReadyMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_nodemanager_nodemanager_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadyMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadyMessage) ProtoMessage() {}

func (x *ReadyMessage) ProtoReflect() protoreflect.Message {
	mi := &file_protos_nodemanager_nodemanager_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadyMessage.ProtoReflect.Descriptor instead.
func (*ReadyMessage) Descriptor() ([]byte, []int) {
	return file_protos_nodemanager_nodemanager_proto_rawDescGZIP(), []int{2}
}

func (x *ReadyMessage) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

type GroupMemberInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Role    int64  `protobuf:"varint,1,opt,name=Role,proto3" json:"Role,omitempty"`
	Address string `protobuf:"bytes,2,opt,name=Address,proto3" json:"Address,omitempty"`
	ID      string `protobuf:"bytes,3,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (x *GroupMemberInfo) Reset() {
	*x = GroupMemberInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_nodemanager_nodemanager_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupMemberInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupMemberInfo) ProtoMessage() {}

func (x *GroupMemberInfo) ProtoReflect() protoreflect.Message {
	mi := &file_protos_nodemanager_nodemanager_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupMemberInfo.ProtoReflect.Descriptor instead.
func (*GroupMemberInfo) Descriptor() ([]byte, []int) {
	return file_protos_nodemanager_nodemanager_proto_rawDescGZIP(), []int{3}
}

func (x *GroupMemberInfo) GetRole() int64 {
	if x != nil {
		return x.Role
	}
	return 0
}

func (x *GroupMemberInfo) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *GroupMemberInfo) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

type GroupInfoMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epoch   int64              `protobuf:"varint,1,opt,name=Epoch,proto3" json:"Epoch,omitempty"`
	NodeID  string             `protobuf:"bytes,2,opt,name=NodeID,proto3" json:"NodeID,omitempty"`
	Members []*GroupMemberInfo `protobuf:"bytes,3,rep,name=Members,proto3" json:"Members,omitempty"`
}

func (x *GroupInfoMessage) Reset() {
	*x = GroupInfoMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_nodemanager_nodemanager_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupInfoMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupInfoMessage) ProtoMessage() {}

func (x *GroupInfoMessage) ProtoReflect() protoreflect.Message {
	mi := &file_protos_nodemanager_nodemanager_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupInfoMessage.ProtoReflect.Descriptor instead.
func (*GroupInfoMessage) Descriptor() ([]byte, []int) {
	return file_protos_nodemanager_nodemanager_proto_rawDescGZIP(), []int{4}
}

func (x *GroupInfoMessage) GetEpoch() int64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *GroupInfoMessage) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

func (x *GroupInfoMessage) GetMembers() []*GroupMemberInfo {
	if x != nil {
		return x.Members
	}
	return nil
}

type PrepareMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeID string `protobuf:"bytes,1,opt,name=NodeID,proto3" json:"NodeID,omitempty"`
}

func (x *PrepareMessage) Reset() {
	*x = PrepareMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_nodemanager_nodemanager_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrepareMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrepareMessage) ProtoMessage() {}

func (x *PrepareMessage) ProtoReflect() protoreflect.Message {
	mi := &file_protos_nodemanager_nodemanager_proto_msgTypes[5]
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
	return file_protos_nodemanager_nodemanager_proto_rawDescGZIP(), []int{5}
}

func (x *PrepareMessage) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

type PromiseMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeID string `protobuf:"bytes,1,opt,name=NodeID,proto3" json:"NodeID,omitempty"`
}

func (x *PromiseMessage) Reset() {
	*x = PromiseMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_nodemanager_nodemanager_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PromiseMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PromiseMessage) ProtoMessage() {}

func (x *PromiseMessage) ProtoReflect() protoreflect.Message {
	mi := &file_protos_nodemanager_nodemanager_proto_msgTypes[6]
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
	return file_protos_nodemanager_nodemanager_proto_rawDescGZIP(), []int{6}
}

func (x *PromiseMessage) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

type AcceptMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epoch int64 `protobuf:"varint,1,opt,name=Epoch,proto3" json:"Epoch,omitempty"`
	Key   int64 `protobuf:"varint,2,opt,name=Key,proto3" json:"Key,omitempty"`
	Value int64 `protobuf:"varint,3,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *AcceptMessage) Reset() {
	*x = AcceptMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_nodemanager_nodemanager_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AcceptMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptMessage) ProtoMessage() {}

func (x *AcceptMessage) ProtoReflect() protoreflect.Message {
	mi := &file_protos_nodemanager_nodemanager_proto_msgTypes[7]
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
	return file_protos_nodemanager_nodemanager_proto_rawDescGZIP(), []int{7}
}

func (x *AcceptMessage) GetEpoch() int64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *AcceptMessage) GetKey() int64 {
	if x != nil {
		return x.Key
	}
	return 0
}

func (x *AcceptMessage) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type LearnMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epoch int64 `protobuf:"varint,1,opt,name=Epoch,proto3" json:"Epoch,omitempty"`
	Key   int64 `protobuf:"varint,2,opt,name=Key,proto3" json:"Key,omitempty"`
}

func (x *LearnMessage) Reset() {
	*x = LearnMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_nodemanager_nodemanager_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LearnMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LearnMessage) ProtoMessage() {}

func (x *LearnMessage) ProtoReflect() protoreflect.Message {
	mi := &file_protos_nodemanager_nodemanager_proto_msgTypes[8]
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
	return file_protos_nodemanager_nodemanager_proto_rawDescGZIP(), []int{8}
}

func (x *LearnMessage) GetEpoch() int64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *LearnMessage) GetKey() int64 {
	if x != nil {
		return x.Key
	}
	return 0
}

type CommitMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epoch int64 `protobuf:"varint,1,opt,name=Epoch,proto3" json:"Epoch,omitempty"`
	Key   int64 `protobuf:"varint,2,opt,name=Key,proto3" json:"Key,omitempty"`
}

func (x *CommitMessage) Reset() {
	*x = CommitMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_nodemanager_nodemanager_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommitMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommitMessage) ProtoMessage() {}

func (x *CommitMessage) ProtoReflect() protoreflect.Message {
	mi := &file_protos_nodemanager_nodemanager_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommitMessage.ProtoReflect.Descriptor instead.
func (*CommitMessage) Descriptor() ([]byte, []int) {
	return file_protos_nodemanager_nodemanager_proto_rawDescGZIP(), []int{9}
}

func (x *CommitMessage) GetEpoch() int64 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *CommitMessage) GetKey() int64 {
	if x != nil {
		return x.Key
	}
	return 0
}

var File_protos_nodemanager_nodemanager_proto protoreflect.FileDescriptor

var file_protos_nodemanager_nodemanager_proto_rawDesc = []byte{
	0x0a, 0x24, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x6d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x72, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x6e, 0x6f, 0x64, 0x65, 0x6d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x72, 0x1a, 0x0c, 0x67, 0x6f, 0x72, 0x75, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3f,
	0x0a, 0x0b, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4e,
	0x6f, 0x64, 0x65, 0x49, 0x44, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22,
	0x58, 0x0a, 0x0c, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x0e, 0x0a, 0x02, 0x4f, 0x4b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x4f, 0x4b, 0x12,
	0x16, 0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x12, 0x20, 0x0a, 0x0b, 0x4e, 0x65, 0x78, 0x74, 0x41,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x4e, 0x65,
	0x78, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x26, 0x0a, 0x0c, 0x52, 0x65, 0x61,
	0x64, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x4e, 0x6f, 0x64,
	0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49,
	0x44, 0x22, 0x4f, 0x0a, 0x0f, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x49, 0x44, 0x22, 0x78, 0x0a, 0x10, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x16, 0x0a, 0x06,
	0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4e, 0x6f,
	0x64, 0x65, 0x49, 0x44, 0x12, 0x36, 0x0a, 0x07, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x6d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x72, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x07, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x22, 0x28, 0x0a, 0x0e,
	0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x22, 0x28, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x6d, 0x69, 0x73,
	0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44,
	0x22, 0x4d, 0x0a, 0x0d, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22,
	0x36, 0x0a, 0x0c, 0x4c, 0x65, 0x61, 0x72, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05,
	0x45, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x22, 0x37, 0x0a, 0x0d, 0x43, 0x6f, 0x6d, 0x6d, 0x69,
	0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x70, 0x6f, 0x63,
	0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x10,
	0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x4b, 0x65, 0x79,
	0x32, 0xb1, 0x03, 0x0a, 0x12, 0x4e, 0x6f, 0x64, 0x65, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x48, 0x0a, 0x09, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1d, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x72, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x04, 0x98, 0xb5, 0x18,
	0x01, 0x12, 0x3d, 0x0a, 0x04, 0x4a, 0x6f, 0x69, 0x6e, 0x12, 0x18, 0x2e, 0x6e, 0x6f, 0x64, 0x65,
	0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x72, 0x2e, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x3c, 0x0a, 0x05, 0x52, 0x65, 0x61, 0x64, 0x79, 0x12, 0x19, 0x2e, 0x6e, 0x6f, 0x64, 0x65,
	0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x79, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x49,
	0x0a, 0x07, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x12, 0x1b, 0x2e, 0x6e, 0x6f, 0x64, 0x65,
	0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x1b, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x6d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x72, 0x2e, 0x50, 0x72, 0x6f, 0x6d, 0x69, 0x73, 0x65, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x22, 0x04, 0xa0, 0xb5, 0x18, 0x01, 0x12, 0x45, 0x0a, 0x06, 0x41, 0x63, 0x63,
	0x65, 0x70, 0x74, 0x12, 0x1a, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x72, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a,
	0x19, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x4c, 0x65,
	0x61, 0x72, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x04, 0xa0, 0xb5, 0x18, 0x01,
	0x12, 0x42, 0x0a, 0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x1a, 0x2e, 0x6e, 0x6f, 0x64,
	0x65, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x04,
	0x98, 0xb5, 0x18, 0x01, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x76, 0x69, 0x64, 0x61, 0x72, 0x61, 0x6e, 0x64, 0x72, 0x65, 0x62, 0x6f, 0x2f,
	0x6f, 0x6e, 0x63, 0x65, 0x74, 0x72, 0x65, 0x65, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x6d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_nodemanager_nodemanager_proto_rawDescOnce sync.Once
	file_protos_nodemanager_nodemanager_proto_rawDescData = file_protos_nodemanager_nodemanager_proto_rawDesc
)

func file_protos_nodemanager_nodemanager_proto_rawDescGZIP() []byte {
	file_protos_nodemanager_nodemanager_proto_rawDescOnce.Do(func() {
		file_protos_nodemanager_nodemanager_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_nodemanager_nodemanager_proto_rawDescData)
	})
	return file_protos_nodemanager_nodemanager_proto_rawDescData
}

var (
	file_protos_nodemanager_nodemanager_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
	file_protos_nodemanager_nodemanager_proto_goTypes  = []interface{}{
		(*JoinRequest)(nil),      // 0: nodemanager.JoinRequest
		(*JoinResponse)(nil),     // 1: nodemanager.JoinResponse
		(*ReadyMessage)(nil),     // 2: nodemanager.ReadyMessage
		(*GroupMemberInfo)(nil),  // 3: nodemanager.GroupMemberInfo
		(*GroupInfoMessage)(nil), // 4: nodemanager.GroupInfoMessage
		(*PrepareMessage)(nil),   // 5: nodemanager.PrepareMessage
		(*PromiseMessage)(nil),   // 6: nodemanager.PromiseMessage
		(*AcceptMessage)(nil),    // 7: nodemanager.AcceptMessage
		(*LearnMessage)(nil),     // 8: nodemanager.LearnMessage
		(*CommitMessage)(nil),    // 9: nodemanager.CommitMessage
		(*emptypb.Empty)(nil),    // 10: google.protobuf.Empty
	}
)
var file_protos_nodemanager_nodemanager_proto_depIdxs = []int32{
	3,  // 0: nodemanager.GroupInfoMessage.Members:type_name -> nodemanager.GroupMemberInfo
	4,  // 1: nodemanager.NodeManagerService.GroupInfo:input_type -> nodemanager.GroupInfoMessage
	0,  // 2: nodemanager.NodeManagerService.Join:input_type -> nodemanager.JoinRequest
	2,  // 3: nodemanager.NodeManagerService.Ready:input_type -> nodemanager.ReadyMessage
	5,  // 4: nodemanager.NodeManagerService.Prepare:input_type -> nodemanager.PrepareMessage
	7,  // 5: nodemanager.NodeManagerService.Accept:input_type -> nodemanager.AcceptMessage
	9,  // 6: nodemanager.NodeManagerService.Commit:input_type -> nodemanager.CommitMessage
	10, // 7: nodemanager.NodeManagerService.GroupInfo:output_type -> google.protobuf.Empty
	1,  // 8: nodemanager.NodeManagerService.Join:output_type -> nodemanager.JoinResponse
	10, // 9: nodemanager.NodeManagerService.Ready:output_type -> google.protobuf.Empty
	6,  // 10: nodemanager.NodeManagerService.Prepare:output_type -> nodemanager.PromiseMessage
	8,  // 11: nodemanager.NodeManagerService.Accept:output_type -> nodemanager.LearnMessage
	10, // 12: nodemanager.NodeManagerService.Commit:output_type -> google.protobuf.Empty
	7,  // [7:13] is the sub-list for method output_type
	1,  // [1:7] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_protos_nodemanager_nodemanager_proto_init() }
func file_protos_nodemanager_nodemanager_proto_init() {
	if File_protos_nodemanager_nodemanager_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_nodemanager_nodemanager_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoinRequest); i {
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
		file_protos_nodemanager_nodemanager_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoinResponse); i {
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
		file_protos_nodemanager_nodemanager_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadyMessage); i {
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
		file_protos_nodemanager_nodemanager_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GroupMemberInfo); i {
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
		file_protos_nodemanager_nodemanager_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GroupInfoMessage); i {
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
		file_protos_nodemanager_nodemanager_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
		file_protos_nodemanager_nodemanager_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
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
		file_protos_nodemanager_nodemanager_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
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
		file_protos_nodemanager_nodemanager_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
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
		file_protos_nodemanager_nodemanager_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommitMessage); i {
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
			RawDescriptor: file_protos_nodemanager_nodemanager_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_nodemanager_nodemanager_proto_goTypes,
		DependencyIndexes: file_protos_nodemanager_nodemanager_proto_depIdxs,
		MessageInfos:      file_protos_nodemanager_nodemanager_proto_msgTypes,
	}.Build()
	File_protos_nodemanager_nodemanager_proto = out.File
	file_protos_nodemanager_nodemanager_proto_rawDesc = nil
	file_protos_nodemanager_nodemanager_proto_goTypes = nil
	file_protos_nodemanager_nodemanager_proto_depIdxs = nil
}
