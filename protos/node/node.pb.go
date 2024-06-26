// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.3
// source: protos/node/node.proto

package node

import (
	reflect "reflect"
	sync "sync"

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

type NodesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeMap map[string]string `protobuf:"bytes,1,rep,name=NodeMap,proto3" json:"NodeMap,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *NodesResponse) Reset() {
	*x = NodesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_node_node_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodesResponse) ProtoMessage() {}

func (x *NodesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_node_node_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodesResponse.ProtoReflect.Descriptor instead.
func (*NodesResponse) Descriptor() ([]byte, []int) {
	return file_protos_node_node_proto_rawDescGZIP(), []int{0}
}

func (x *NodesResponse) GetNodeMap() map[string]string {
	if x != nil {
		return x.NodeMap
	}
	return nil
}

type NodesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Origin string `protobuf:"bytes,1,opt,name=Origin,proto3" json:"Origin,omitempty"`
}

func (x *NodesRequest) Reset() {
	*x = NodesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_node_node_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodesRequest) ProtoMessage() {}

func (x *NodesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_node_node_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodesRequest.ProtoReflect.Descriptor instead.
func (*NodesRequest) Descriptor() ([]byte, []int) {
	return file_protos_node_node_proto_rawDescGZIP(), []int{1}
}

func (x *NodesRequest) GetOrigin() string {
	if x != nil {
		return x.Origin
	}
	return ""
}

type IDResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (x *IDResponse) Reset() {
	*x = IDResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_node_node_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IDResponse) ProtoMessage() {}

func (x *IDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_node_node_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IDResponse.ProtoReflect.Descriptor instead.
func (*IDResponse) Descriptor() ([]byte, []int) {
	return file_protos_node_node_proto_rawDescGZIP(), []int{2}
}

func (x *IDResponse) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

var File_protos_node_node_proto protoreflect.FileDescriptor

var file_protos_node_node_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x6e, 0x6f,
	0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x87, 0x01, 0x0a, 0x0d,
	0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a,
	0x07, 0x4e, 0x6f, 0x64, 0x65, 0x4d, 0x61, 0x70, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20,
	0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x07, 0x4e, 0x6f, 0x64, 0x65, 0x4d, 0x61, 0x70, 0x1a, 0x3a, 0x0a, 0x0c, 0x4e, 0x6f, 0x64,
	0x65, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x26, 0x0a, 0x0c, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x22, 0x1c, 0x0a,
	0x0a, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49,
	0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x32, 0xb2, 0x01, 0x0a, 0x0b,
	0x4e, 0x6f, 0x64, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x39, 0x0a, 0x05, 0x43,
	0x72, 0x61, 0x73, 0x68, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x32, 0x0a, 0x05, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x12,
	0x12, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x34, 0x0a, 0x06, 0x4e, 0x6f,
	0x64, 0x65, 0x49, 0x44, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x10, 0x2e, 0x6e,
	0x6f, 0x64, 0x65, 0x2e, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x27, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x76,
	0x69, 0x64, 0x61, 0x72, 0x61, 0x6e, 0x64, 0x72, 0x65, 0x62, 0x6f, 0x2f, 0x6f, 0x6e, 0x63, 0x65,
	0x74, 0x72, 0x65, 0x65, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_protos_node_node_proto_rawDescOnce sync.Once
	file_protos_node_node_proto_rawDescData = file_protos_node_node_proto_rawDesc
)

func file_protos_node_node_proto_rawDescGZIP() []byte {
	file_protos_node_node_proto_rawDescOnce.Do(func() {
		file_protos_node_node_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_node_node_proto_rawDescData)
	})
	return file_protos_node_node_proto_rawDescData
}

var (
	file_protos_node_node_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
	file_protos_node_node_proto_goTypes  = []interface{}{
		(*NodesResponse)(nil), // 0: node.NodesResponse
		(*NodesRequest)(nil),  // 1: node.NodesRequest
		(*IDResponse)(nil),    // 2: node.IDResponse
		nil,                   // 3: node.NodesResponse.NodeMapEntry
		(*emptypb.Empty)(nil), // 4: google.protobuf.Empty
	}
)

var file_protos_node_node_proto_depIdxs = []int32{
	3, // 0: node.NodesResponse.NodeMap:type_name -> node.NodesResponse.NodeMapEntry
	4, // 1: node.NodeService.Crash:input_type -> google.protobuf.Empty
	1, // 2: node.NodeService.Nodes:input_type -> node.NodesRequest
	4, // 3: node.NodeService.NodeID:input_type -> google.protobuf.Empty
	4, // 4: node.NodeService.Crash:output_type -> google.protobuf.Empty
	0, // 5: node.NodeService.Nodes:output_type -> node.NodesResponse
	2, // 6: node.NodeService.NodeID:output_type -> node.IDResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_protos_node_node_proto_init() }
func file_protos_node_node_proto_init() {
	if File_protos_node_node_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_node_node_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodesResponse); i {
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
		file_protos_node_node_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodesRequest); i {
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
		file_protos_node_node_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IDResponse); i {
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
			RawDescriptor: file_protos_node_node_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_node_node_proto_goTypes,
		DependencyIndexes: file_protos_node_node_proto_depIdxs,
		MessageInfos:      file_protos_node_node_proto_msgTypes,
	}.Build()
	File_protos_node_node_proto = out.File
	file_protos_node_node_proto_rawDesc = nil
	file_protos_node_node_proto_goTypes = nil
	file_protos_node_node_proto_depIdxs = nil
}
