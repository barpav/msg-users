// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.23.4
// source: users_service_go_grpc/users_service.proto

package users_service_go_grpc

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

type Credentials struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *Credentials) Reset() {
	*x = Credentials{}
	if protoimpl.UnsafeEnabled {
		mi := &file_users_service_go_grpc_users_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Credentials) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Credentials) ProtoMessage() {}

func (x *Credentials) ProtoReflect() protoreflect.Message {
	mi := &file_users_service_go_grpc_users_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Credentials.ProtoReflect.Descriptor instead.
func (*Credentials) Descriptor() ([]byte, []int) {
	return file_users_service_go_grpc_users_service_proto_rawDescGZIP(), []int{0}
}

func (x *Credentials) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Credentials) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type ValidationResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Valid bool `protobuf:"varint,1,opt,name=valid,proto3" json:"valid,omitempty"`
}

func (x *ValidationResult) Reset() {
	*x = ValidationResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_users_service_go_grpc_users_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidationResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidationResult) ProtoMessage() {}

func (x *ValidationResult) ProtoReflect() protoreflect.Message {
	mi := &file_users_service_go_grpc_users_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidationResult.ProtoReflect.Descriptor instead.
func (*ValidationResult) Descriptor() ([]byte, []int) {
	return file_users_service_go_grpc_users_service_proto_rawDescGZIP(), []int{1}
}

func (x *ValidationResult) GetValid() bool {
	if x != nil {
		return x.Valid
	}
	return false
}

var File_users_service_go_grpc_users_service_proto protoreflect.FileDescriptor

var file_users_service_go_grpc_users_service_proto_rawDesc = []byte{
	0x0a, 0x29, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f,
	0x67, 0x6f, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x6d, 0x73, 0x67,
	0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x22, 0x39, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x61, 0x6c, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x22, 0x28, 0x0a, 0x10, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x32, 0x48, 0x0a, 0x05, 0x55,
	0x73, 0x65, 0x72, 0x73, 0x12, 0x3f, 0x0a, 0x08, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x12, 0x16, 0x2e, 0x6d, 0x73, 0x67, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x43, 0x72, 0x65,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x1a, 0x1b, 0x2e, 0x6d, 0x73, 0x67, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x42, 0x33, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x61, 0x72, 0x70, 0x61, 0x76, 0x2f, 0x6d, 0x73, 0x67, 0x2d, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x5f, 0x67, 0x6f, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_users_service_go_grpc_users_service_proto_rawDescOnce sync.Once
	file_users_service_go_grpc_users_service_proto_rawDescData = file_users_service_go_grpc_users_service_proto_rawDesc
)

func file_users_service_go_grpc_users_service_proto_rawDescGZIP() []byte {
	file_users_service_go_grpc_users_service_proto_rawDescOnce.Do(func() {
		file_users_service_go_grpc_users_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_users_service_go_grpc_users_service_proto_rawDescData)
	})
	return file_users_service_go_grpc_users_service_proto_rawDescData
}

var file_users_service_go_grpc_users_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_users_service_go_grpc_users_service_proto_goTypes = []interface{}{
	(*Credentials)(nil),      // 0: msg.users.Credentials
	(*ValidationResult)(nil), // 1: msg.users.ValidationResult
}
var file_users_service_go_grpc_users_service_proto_depIdxs = []int32{
	0, // 0: msg.users.Users.Validate:input_type -> msg.users.Credentials
	1, // 1: msg.users.Users.Validate:output_type -> msg.users.ValidationResult
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_users_service_go_grpc_users_service_proto_init() }
func file_users_service_go_grpc_users_service_proto_init() {
	if File_users_service_go_grpc_users_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_users_service_go_grpc_users_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Credentials); i {
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
		file_users_service_go_grpc_users_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValidationResult); i {
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
			RawDescriptor: file_users_service_go_grpc_users_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_users_service_go_grpc_users_service_proto_goTypes,
		DependencyIndexes: file_users_service_go_grpc_users_service_proto_depIdxs,
		MessageInfos:      file_users_service_go_grpc_users_service_proto_msgTypes,
	}.Build()
	File_users_service_go_grpc_users_service_proto = out.File
	file_users_service_go_grpc_users_service_proto_rawDesc = nil
	file_users_service_go_grpc_users_service_proto_goTypes = nil
	file_users_service_go_grpc_users_service_proto_depIdxs = nil
}
