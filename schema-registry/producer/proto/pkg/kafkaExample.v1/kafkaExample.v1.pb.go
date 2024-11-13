// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.3
// source: kafkaExample.v1.proto

package kafkaExample_v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ExampleMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      int32   `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Name    string  `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Age     int32   `protobuf:"varint,3,opt,name=Age,proto3" json:"Age,omitempty"`
	Student bool    `protobuf:"varint,4,opt,name=Student,proto3" json:"Student,omitempty"`
	City    *string `protobuf:"bytes,5,opt,name=City,proto3,oneof" json:"City,omitempty"`
}

func (x *ExampleMessage) Reset() {
	*x = ExampleMessage{}
	mi := &file_kafkaExample_v1_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExampleMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExampleMessage) ProtoMessage() {}

func (x *ExampleMessage) ProtoReflect() protoreflect.Message {
	mi := &file_kafkaExample_v1_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExampleMessage.ProtoReflect.Descriptor instead.
func (*ExampleMessage) Descriptor() ([]byte, []int) {
	return file_kafkaExample_v1_proto_rawDescGZIP(), []int{0}
}

func (x *ExampleMessage) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ExampleMessage) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ExampleMessage) GetAge() int32 {
	if x != nil {
		return x.Age
	}
	return 0
}

func (x *ExampleMessage) GetStudent() bool {
	if x != nil {
		return x.Student
	}
	return false
}

func (x *ExampleMessage) GetCity() string {
	if x != nil && x.City != nil {
		return *x.City
	}
	return ""
}

var File_kafkaExample_v1_proto protoreflect.FileDescriptor

var file_kafkaExample_v1_proto_rawDesc = []byte{
	0x0a, 0x15, 0x6b, 0x61, 0x66, 0x6b, 0x61, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x6b, 0x61, 0x66, 0x6b, 0x61, 0x45, 0x78,
	0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x22, 0x82, 0x01, 0x0a, 0x0e, 0x45, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x10, 0x0a, 0x03, 0x41, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x41, 0x67,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x12, 0x17, 0x0a, 0x04, 0x43,
	0x69, 0x74, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x43, 0x69, 0x74,
	0x79, 0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x43, 0x69, 0x74, 0x79, 0x42, 0x17, 0x5a,
	0x15, 0x2e, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6b, 0x61, 0x66, 0x6b, 0x61, 0x45, 0x78, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x2e, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_kafkaExample_v1_proto_rawDescOnce sync.Once
	file_kafkaExample_v1_proto_rawDescData = file_kafkaExample_v1_proto_rawDesc
)

func file_kafkaExample_v1_proto_rawDescGZIP() []byte {
	file_kafkaExample_v1_proto_rawDescOnce.Do(func() {
		file_kafkaExample_v1_proto_rawDescData = protoimpl.X.CompressGZIP(file_kafkaExample_v1_proto_rawDescData)
	})
	return file_kafkaExample_v1_proto_rawDescData
}

var file_kafkaExample_v1_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_kafkaExample_v1_proto_goTypes = []any{
	(*ExampleMessage)(nil), // 0: kafkaExample.v1.ExampleMessage
}
var file_kafkaExample_v1_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_kafkaExample_v1_proto_init() }
func file_kafkaExample_v1_proto_init() {
	if File_kafkaExample_v1_proto != nil {
		return
	}
	file_kafkaExample_v1_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_kafkaExample_v1_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_kafkaExample_v1_proto_goTypes,
		DependencyIndexes: file_kafkaExample_v1_proto_depIdxs,
		MessageInfos:      file_kafkaExample_v1_proto_msgTypes,
	}.Build()
	File_kafkaExample_v1_proto = out.File
	file_kafkaExample_v1_proto_rawDesc = nil
	file_kafkaExample_v1_proto_goTypes = nil
	file_kafkaExample_v1_proto_depIdxs = nil
}
