// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.3
// source: protos/keyvaluestorage/keyvaluestorage.proto

package keyvaluestorage

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

type ReadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key int64 `protobuf:"varint,1,opt,name=Key,proto3" json:"Key,omitempty"`
}

func (x *ReadRequest) Reset() {
	*x = ReadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_keyvaluestorage_keyvaluestorage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadRequest) ProtoMessage() {}

func (x *ReadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_keyvaluestorage_keyvaluestorage_proto_msgTypes[0]
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
	return file_protos_keyvaluestorage_keyvaluestorage_proto_rawDescGZIP(), []int{0}
}

func (x *ReadRequest) GetKey() int64 {
	if x != nil {
		return x.Key
	}
	return 0
}

type ReadLocalRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key    int64  `protobuf:"varint,1,opt,name=Key,proto3" json:"Key,omitempty"`
	NodeID string `protobuf:"bytes,2,opt,name=GroupID,proto3" json:"GroupID,omitempty"`
}

func (x *ReadLocalRequest) Reset() {
	*x = ReadLocalRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_keyvaluestorage_keyvaluestorage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadLocalRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadLocalRequest) ProtoMessage() {}

func (x *ReadLocalRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_keyvaluestorage_keyvaluestorage_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadLocalRequest.ProtoReflect.Descriptor instead.
func (*ReadLocalRequest) Descriptor() ([]byte, []int) {
	return file_protos_keyvaluestorage_keyvaluestorage_proto_rawDescGZIP(), []int{1}
}

func (x *ReadLocalRequest) GetKey() int64 {
	if x != nil {
		return x.Key
	}
	return 0
}

func (x *ReadLocalRequest) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
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
		mi := &file_protos_keyvaluestorage_keyvaluestorage_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadResponse) ProtoMessage() {}

func (x *ReadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_keyvaluestorage_keyvaluestorage_proto_msgTypes[2]
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
	return file_protos_keyvaluestorage_keyvaluestorage_proto_rawDescGZIP(), []int{2}
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
		mi := &file_protos_keyvaluestorage_keyvaluestorage_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadAllResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadAllResponse) ProtoMessage() {}

func (x *ReadAllResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_keyvaluestorage_keyvaluestorage_proto_msgTypes[3]
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
	return file_protos_keyvaluestorage_keyvaluestorage_proto_rawDescGZIP(), []int{3}
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

	Key     int64 `protobuf:"varint,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value   int64 `protobuf:"varint,2,opt,name=Value,proto3" json:"Value,omitempty"`
	WriteID int64 `protobuf:"varint,3,opt,name=WriteID,proto3" json:"WriteID,omitempty"`
}

func (x *WriteRequest) Reset() {
	*x = WriteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_keyvaluestorage_keyvaluestorage_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WriteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriteRequest) ProtoMessage() {}

func (x *WriteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_keyvaluestorage_keyvaluestorage_proto_msgTypes[4]
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
	return file_protos_keyvaluestorage_keyvaluestorage_proto_rawDescGZIP(), []int{4}
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

func (x *WriteRequest) GetWriteID() int64 {
	if x != nil {
		return x.WriteID
	}
	return 0
}

type GossipMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeID         string `protobuf:"bytes,1,opt,name=GroupID,proto3" json:"GroupID,omitempty"`
	Key            int64  `protobuf:"varint,2,opt,name=Key,proto3" json:"Key,omitempty"`
	AggValue       int64  `protobuf:"varint,3,opt,name=AggValue,proto3" json:"AggValue,omitempty"`
	AggTimestamp   int64  `protobuf:"varint,4,opt,name=AggTimestamp,proto3" json:"AggTimestamp,omitempty"`
	LocalValue     int64  `protobuf:"varint,5,opt,name=LocalValue,proto3" json:"LocalValue,omitempty"`
	LocalTimestamp int64  `protobuf:"varint,6,opt,name=LocalTimestamp,proto3" json:"LocalTimestamp,omitempty"`
	WriteID        int64  `protobuf:"varint,7,opt,name=WriteID,proto3" json:"WriteID,omitempty"`
}

func (x *GossipMessage) Reset() {
	*x = GossipMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_keyvaluestorage_keyvaluestorage_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GossipMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GossipMessage) ProtoMessage() {}

func (x *GossipMessage) ProtoReflect() protoreflect.Message {
	mi := &file_protos_keyvaluestorage_keyvaluestorage_proto_msgTypes[5]
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
	return file_protos_keyvaluestorage_keyvaluestorage_proto_rawDescGZIP(), []int{5}
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

func (x *GossipMessage) GetAggValue() int64 {
	if x != nil {
		return x.AggValue
	}
	return 0
}

func (x *GossipMessage) GetAggTimestamp() int64 {
	if x != nil {
		return x.AggTimestamp
	}
	return 0
}

func (x *GossipMessage) GetLocalValue() int64 {
	if x != nil {
		return x.LocalValue
	}
	return 0
}

func (x *GossipMessage) GetLocalTimestamp() int64 {
	if x != nil {
		return x.LocalTimestamp
	}
	return 0
}

func (x *GossipMessage) GetWriteID() int64 {
	if x != nil {
		return x.WriteID
	}
	return 0
}

var File_protos_keyvaluestorage_keyvaluestorage_proto protoreflect.FileDescriptor

var file_protos_keyvaluestorage_keyvaluestorage_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6b, 0x65, 0x79, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2f, 0x6b, 0x65, 0x79, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f,
	0x6b, 0x65, 0x79, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x1a,
	0x0c, 0x67, 0x6f, 0x72, 0x75, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65,
	0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1f, 0x0a, 0x0b, 0x52, 0x65,
	0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x22, 0x3c, 0x0a, 0x10, 0x52,
	0x65, 0x61, 0x64, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x4b, 0x65,
	0x79, 0x12, 0x16, 0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x22, 0x24, 0x0a, 0x0c, 0x52, 0x65, 0x61,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22,
	0x8e, 0x01, 0x0a, 0x0f, 0x52, 0x65, 0x61, 0x64, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x6b, 0x65, 0x79, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x38, 0x0a, 0x0a, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x22, 0x50, 0x0a, 0x0c, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x4b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x57, 0x72, 0x69, 0x74,
	0x65, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x57, 0x72, 0x69, 0x74, 0x65,
	0x49, 0x44, 0x22, 0xdb, 0x01, 0x0a, 0x0d, 0x47, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x12, 0x10, 0x0a, 0x03,
	0x4b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x1a,
	0x0a, 0x08, 0x41, 0x67, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x08, 0x41, 0x67, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x41, 0x67,
	0x67, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0c, 0x41, 0x67, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1e,
	0x0a, 0x0a, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0a, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x26,
	0x0a, 0x0e, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x57, 0x72, 0x69, 0x74, 0x65, 0x49,
	0x44, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x57, 0x72, 0x69, 0x74, 0x65, 0x49, 0x44,
	0x32, 0xc0, 0x03, 0x0a, 0x0f, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x53, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x12, 0x40, 0x0a, 0x05, 0x57, 0x72, 0x69, 0x74, 0x65, 0x12, 0x1d, 0x2e,
	0x6b, 0x65, 0x79, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e,
	0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x45, 0x0a, 0x04, 0x52, 0x65, 0x61, 0x64, 0x12, 0x1c,
	0x2e, 0x6b, 0x65, 0x79, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x2e, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x6b,
	0x65, 0x79, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x52,
	0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4f, 0x0a,
	0x07, 0x52, 0x65, 0x61, 0x64, 0x41, 0x6c, 0x6c, 0x12, 0x1c, 0x2e, 0x6b, 0x65, 0x79, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x6b, 0x65, 0x79, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x41, 0x6c, 0x6c,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x04, 0xa0, 0xb5, 0x18, 0x01, 0x12, 0x4f,
	0x0a, 0x09, 0x52, 0x65, 0x61, 0x64, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x12, 0x21, 0x2e, 0x6b, 0x65,
	0x79, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x52, 0x65,
	0x61, 0x64, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d,
	0x2e, 0x6b, 0x65, 0x79, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x2e, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x42, 0x0a, 0x06, 0x47, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x12, 0x1e, 0x2e, 0x6b, 0x65, 0x79, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x47, 0x6f, 0x73, 0x73,
	0x69, 0x70, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x00, 0x12, 0x3e, 0x0a, 0x0a, 0x50, 0x72, 0x69, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x00, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x76, 0x69, 0x64, 0x61, 0x72, 0x61, 0x6e, 0x64, 0x72, 0x65, 0x62, 0x6f, 0x2f, 0x6f,
	0x6e, 0x63, 0x65, 0x74, 0x72, 0x65, 0x65, 0x2f, 0x6b, 0x65, 0x79, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_keyvaluestorage_keyvaluestorage_proto_rawDescOnce sync.Once
	file_protos_keyvaluestorage_keyvaluestorage_proto_rawDescData = file_protos_keyvaluestorage_keyvaluestorage_proto_rawDesc
)

func file_protos_keyvaluestorage_keyvaluestorage_proto_rawDescGZIP() []byte {
	file_protos_keyvaluestorage_keyvaluestorage_proto_rawDescOnce.Do(func() {
		file_protos_keyvaluestorage_keyvaluestorage_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_keyvaluestorage_keyvaluestorage_proto_rawDescData)
	})
	return file_protos_keyvaluestorage_keyvaluestorage_proto_rawDescData
}

var (
	file_protos_keyvaluestorage_keyvaluestorage_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
	file_protos_keyvaluestorage_keyvaluestorage_proto_goTypes  = []interface{}{
		(*ReadRequest)(nil),      // 0: keyvaluestorage.ReadRequest
		(*ReadLocalRequest)(nil), // 1: keyvaluestorage.ReadLocalRequest
		(*ReadResponse)(nil),     // 2: keyvaluestorage.ReadResponse
		(*ReadAllResponse)(nil),  // 3: keyvaluestorage.ReadAllResponse
		(*WriteRequest)(nil),     // 4: keyvaluestorage.WriteRequest
		(*GossipMessage)(nil),    // 5: keyvaluestorage.GossipMessage
		nil,                      // 6: keyvaluestorage.ReadAllResponse.ValueEntry
		(*emptypb.Empty)(nil),    // 7: google.protobuf.Empty
	}
)

var file_protos_keyvaluestorage_keyvaluestorage_proto_depIdxs = []int32{
	6, // 0: keyvaluestorage.ReadAllResponse.Value:type_name -> keyvaluestorage.ReadAllResponse.ValueEntry
	4, // 1: keyvaluestorage.KeyValueStorage.Write:input_type -> keyvaluestorage.WriteRequest
	0, // 2: keyvaluestorage.KeyValueStorage.Read:input_type -> keyvaluestorage.ReadRequest
	0, // 3: keyvaluestorage.KeyValueStorage.ReadAll:input_type -> keyvaluestorage.ReadRequest
	1, // 4: keyvaluestorage.KeyValueStorage.ReadLocal:input_type -> keyvaluestorage.ReadLocalRequest
	5, // 5: keyvaluestorage.KeyValueStorage.Gossip:input_type -> keyvaluestorage.GossipMessage
	7, // 6: keyvaluestorage.KeyValueStorage.PrintState:input_type -> google.protobuf.Empty
	7, // 7: keyvaluestorage.KeyValueStorage.Write:output_type -> google.protobuf.Empty
	2, // 8: keyvaluestorage.KeyValueStorage.Read:output_type -> keyvaluestorage.ReadResponse
	3, // 9: keyvaluestorage.KeyValueStorage.ReadAll:output_type -> keyvaluestorage.ReadAllResponse
	2, // 10: keyvaluestorage.KeyValueStorage.ReadLocal:output_type -> keyvaluestorage.ReadResponse
	7, // 11: keyvaluestorage.KeyValueStorage.Gossip:output_type -> google.protobuf.Empty
	7, // 12: keyvaluestorage.KeyValueStorage.PrintState:output_type -> google.protobuf.Empty
	7, // [7:13] is the sub-list for method output_type
	1, // [1:7] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_protos_keyvaluestorage_keyvaluestorage_proto_init() }
func file_protos_keyvaluestorage_keyvaluestorage_proto_init() {
	if File_protos_keyvaluestorage_keyvaluestorage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_keyvaluestorage_keyvaluestorage_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_protos_keyvaluestorage_keyvaluestorage_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadLocalRequest); i {
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
		file_protos_keyvaluestorage_keyvaluestorage_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_protos_keyvaluestorage_keyvaluestorage_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
		file_protos_keyvaluestorage_keyvaluestorage_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
		file_protos_keyvaluestorage_keyvaluestorage_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protos_keyvaluestorage_keyvaluestorage_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_keyvaluestorage_keyvaluestorage_proto_goTypes,
		DependencyIndexes: file_protos_keyvaluestorage_keyvaluestorage_proto_depIdxs,
		MessageInfos:      file_protos_keyvaluestorage_keyvaluestorage_proto_msgTypes,
	}.Build()
	File_protos_keyvaluestorage_keyvaluestorage_proto = out.File
	file_protos_keyvaluestorage_keyvaluestorage_proto_rawDesc = nil
	file_protos_keyvaluestorage_keyvaluestorage_proto_goTypes = nil
	file_protos_keyvaluestorage_keyvaluestorage_proto_depIdxs = nil
}
