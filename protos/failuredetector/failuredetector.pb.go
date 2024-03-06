// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.3
// source: protos/failuredetector/failuredetector.proto

package failuredetector

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

type HeartbeatMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeID string `protobuf:"bytes,1,opt,name=NodeID,proto3" json:"NodeID,omitempty"`
}

func (x *HeartbeatMessage) Reset() {
	*x = HeartbeatMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_failuredetector_failuredetector_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeartbeatMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartbeatMessage) ProtoMessage() {}

func (x *HeartbeatMessage) ProtoReflect() protoreflect.Message {
	mi := &file_protos_failuredetector_failuredetector_proto_msgTypes[0]
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
	return file_protos_failuredetector_failuredetector_proto_rawDescGZIP(), []int{0}
}

func (x *HeartbeatMessage) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

var File_protos_failuredetector_failuredetector_proto protoreflect.FileDescriptor

var file_protos_failuredetector_failuredetector_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65,
	0x64, 0x65, 0x74, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2f, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65,
	0x64, 0x65, 0x74, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f,
	0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x64, 0x65, 0x74, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x1a,
	0x0c, 0x67, 0x6f, 0x72, 0x75, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65,
	0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2a, 0x0a, 0x10, 0x48, 0x65,
	0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x32, 0x66, 0x0a, 0x16, 0x46, 0x61, 0x69, 0x6c, 0x75, 0x72,
	0x65, 0x44, 0x65, 0x74, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x4c, 0x0a, 0x09, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x12, 0x21, 0x2e,
	0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x64, 0x65, 0x74, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e,
	0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x04, 0x98, 0xb5, 0x18, 0x01, 0x42, 0x32,
	0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x76, 0x69, 0x64,
	0x61, 0x72, 0x61, 0x6e, 0x64, 0x72, 0x65, 0x62, 0x6f, 0x2f, 0x6f, 0x6e, 0x63, 0x65, 0x74, 0x72,
	0x65, 0x65, 0x2f, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x64, 0x65, 0x74, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_failuredetector_failuredetector_proto_rawDescOnce sync.Once
	file_protos_failuredetector_failuredetector_proto_rawDescData = file_protos_failuredetector_failuredetector_proto_rawDesc
)

func file_protos_failuredetector_failuredetector_proto_rawDescGZIP() []byte {
	file_protos_failuredetector_failuredetector_proto_rawDescOnce.Do(func() {
		file_protos_failuredetector_failuredetector_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_failuredetector_failuredetector_proto_rawDescData)
	})
	return file_protos_failuredetector_failuredetector_proto_rawDescData
}

var (
	file_protos_failuredetector_failuredetector_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
	file_protos_failuredetector_failuredetector_proto_goTypes  = []interface{}{
		(*HeartbeatMessage)(nil), // 0: failuredetector.HeartbeatMessage
		(*emptypb.Empty)(nil),    // 1: google.protobuf.Empty
	}
)
var file_protos_failuredetector_failuredetector_proto_depIdxs = []int32{
	0, // 0: failuredetector.FailureDetectorService.Heartbeat:input_type -> failuredetector.HeartbeatMessage
	1, // 1: failuredetector.FailureDetectorService.Heartbeat:output_type -> google.protobuf.Empty
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_protos_failuredetector_failuredetector_proto_init() }
func file_protos_failuredetector_failuredetector_proto_init() {
	if File_protos_failuredetector_failuredetector_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_failuredetector_failuredetector_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protos_failuredetector_failuredetector_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_failuredetector_failuredetector_proto_goTypes,
		DependencyIndexes: file_protos_failuredetector_failuredetector_proto_depIdxs,
		MessageInfos:      file_protos_failuredetector_failuredetector_proto_msgTypes,
	}.Build()
	File_protos_failuredetector_failuredetector_proto = out.File
	file_protos_failuredetector_failuredetector_proto_rawDesc = nil
	file_protos_failuredetector_failuredetector_proto_goTypes = nil
	file_protos_failuredetector_failuredetector_proto_depIdxs = nil
}
