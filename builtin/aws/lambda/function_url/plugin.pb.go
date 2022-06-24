// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: waypoint/builtin/aws/lambda/function_url/plugin.proto

package function_url

import (
	opaqueany "github.com/hashicorp/opaqueany"
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

type Release struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The function's public url
	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	// The AWS region the function is deployed in
	Region string `protobuf:"bytes,2,opt,name=region,proto3" json:"region,omitempty"`
	// The ARN for the Lambda function itself.
	FuncArn string `protobuf:"bytes,3,opt,name=func_arn,json=funcArn,proto3" json:"func_arn,omitempty"`
	// The ARN for the version of the Lambda function this deployment uses.
	VerArn        string         `protobuf:"bytes,4,opt,name=ver_arn,json=verArn,proto3" json:"ver_arn,omitempty"`
	ResourceState *opaqueany.Any `protobuf:"bytes,5,opt,name=resource_state,json=resourceState,proto3" json:"resource_state,omitempty"`
}

func (x *Release) Reset() {
	*x = Release{}
	if protoimpl.UnsafeEnabled {
		mi := &file_waypoint_builtin_aws_lambda_function_url_plugin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Release) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Release) ProtoMessage() {}

func (x *Release) ProtoReflect() protoreflect.Message {
	mi := &file_waypoint_builtin_aws_lambda_function_url_plugin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Release.ProtoReflect.Descriptor instead.
func (*Release) Descriptor() ([]byte, []int) {
	return file_waypoint_builtin_aws_lambda_function_url_plugin_proto_rawDescGZIP(), []int{0}
}

func (x *Release) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Release) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *Release) GetFuncArn() string {
	if x != nil {
		return x.FuncArn
	}
	return ""
}

func (x *Release) GetVerArn() string {
	if x != nil {
		return x.VerArn
	}
	return ""
}

func (x *Release) GetResourceState() *opaqueany.Any {
	if x != nil {
		return x.ResourceState
	}
	return nil
}

// Resource contains the internal resource states.
type Resource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Resource) Reset() {
	*x = Resource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_waypoint_builtin_aws_lambda_function_url_plugin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Resource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Resource) ProtoMessage() {}

func (x *Resource) ProtoReflect() protoreflect.Message {
	mi := &file_waypoint_builtin_aws_lambda_function_url_plugin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Resource.ProtoReflect.Descriptor instead.
func (*Resource) Descriptor() ([]byte, []int) {
	return file_waypoint_builtin_aws_lambda_function_url_plugin_proto_rawDescGZIP(), []int{1}
}

type Resource_FunctionUrl struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *Resource_FunctionUrl) Reset() {
	*x = Resource_FunctionUrl{}
	if protoimpl.UnsafeEnabled {
		mi := &file_waypoint_builtin_aws_lambda_function_url_plugin_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Resource_FunctionUrl) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Resource_FunctionUrl) ProtoMessage() {}

func (x *Resource_FunctionUrl) ProtoReflect() protoreflect.Message {
	mi := &file_waypoint_builtin_aws_lambda_function_url_plugin_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Resource_FunctionUrl.ProtoReflect.Descriptor instead.
func (*Resource_FunctionUrl) Descriptor() ([]byte, []int) {
	return file_waypoint_builtin_aws_lambda_function_url_plugin_proto_rawDescGZIP(), []int{1, 0}
}

func (x *Resource_FunctionUrl) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

var File_waypoint_builtin_aws_lambda_function_url_plugin_proto protoreflect.FileDescriptor

var file_waypoint_builtin_aws_lambda_function_url_plugin_proto_rawDesc = []byte{
	0x0a, 0x35, 0x77, 0x61, 0x79, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x74,
	0x69, 0x6e, 0x2f, 0x61, 0x77, 0x73, 0x2f, 0x6c, 0x61, 0x6d, 0x62, 0x64, 0x61, 0x2f, 0x66, 0x75,
	0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x75, 0x72, 0x6c, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x75, 0x72, 0x6c, 0x1a, 0x13, 0x6f, 0x70, 0x61, 0x71, 0x75, 0x65, 0x61, 0x6e, 0x79,
	0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9e, 0x01, 0x0a, 0x07, 0x52,
	0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e,
	0x12, 0x19, 0x0a, 0x08, 0x66, 0x75, 0x6e, 0x63, 0x5f, 0x61, 0x72, 0x6e, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x66, 0x75, 0x6e, 0x63, 0x41, 0x72, 0x6e, 0x12, 0x17, 0x0a, 0x07, 0x76,
	0x65, 0x72, 0x5f, 0x61, 0x72, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x76, 0x65,
	0x72, 0x41, 0x72, 0x6e, 0x12, 0x35, 0x0a, 0x0e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6f,
	0x70, 0x61, 0x71, 0x75, 0x65, 0x61, 0x6e, 0x79, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x0d, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x22, 0x2b, 0x0a, 0x08, 0x52,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x1a, 0x1f, 0x0a, 0x0b, 0x46, 0x75, 0x6e, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x55, 0x72, 0x6c, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x42, 0x2a, 0x5a, 0x28, 0x77, 0x61, 0x79, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x74, 0x69, 0x6e, 0x2f, 0x61, 0x77, 0x73,
	0x2f, 0x6c, 0x61, 0x6d, 0x62, 0x64, 0x61, 0x2f, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x75, 0x72, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_waypoint_builtin_aws_lambda_function_url_plugin_proto_rawDescOnce sync.Once
	file_waypoint_builtin_aws_lambda_function_url_plugin_proto_rawDescData = file_waypoint_builtin_aws_lambda_function_url_plugin_proto_rawDesc
)

func file_waypoint_builtin_aws_lambda_function_url_plugin_proto_rawDescGZIP() []byte {
	file_waypoint_builtin_aws_lambda_function_url_plugin_proto_rawDescOnce.Do(func() {
		file_waypoint_builtin_aws_lambda_function_url_plugin_proto_rawDescData = protoimpl.X.CompressGZIP(file_waypoint_builtin_aws_lambda_function_url_plugin_proto_rawDescData)
	})
	return file_waypoint_builtin_aws_lambda_function_url_plugin_proto_rawDescData
}

var file_waypoint_builtin_aws_lambda_function_url_plugin_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_waypoint_builtin_aws_lambda_function_url_plugin_proto_goTypes = []interface{}{
	(*Release)(nil),              // 0: function_url.Release
	(*Resource)(nil),             // 1: function_url.Resource
	(*Resource_FunctionUrl)(nil), // 2: function_url.Resource.FunctionUrl
	(*opaqueany.Any)(nil),        // 3: opaqueany.Any
}
var file_waypoint_builtin_aws_lambda_function_url_plugin_proto_depIdxs = []int32{
	3, // 0: function_url.Release.resource_state:type_name -> opaqueany.Any
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_waypoint_builtin_aws_lambda_function_url_plugin_proto_init() }
func file_waypoint_builtin_aws_lambda_function_url_plugin_proto_init() {
	if File_waypoint_builtin_aws_lambda_function_url_plugin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_waypoint_builtin_aws_lambda_function_url_plugin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Release); i {
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
		file_waypoint_builtin_aws_lambda_function_url_plugin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Resource); i {
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
		file_waypoint_builtin_aws_lambda_function_url_plugin_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Resource_FunctionUrl); i {
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
			RawDescriptor: file_waypoint_builtin_aws_lambda_function_url_plugin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_waypoint_builtin_aws_lambda_function_url_plugin_proto_goTypes,
		DependencyIndexes: file_waypoint_builtin_aws_lambda_function_url_plugin_proto_depIdxs,
		MessageInfos:      file_waypoint_builtin_aws_lambda_function_url_plugin_proto_msgTypes,
	}.Build()
	File_waypoint_builtin_aws_lambda_function_url_plugin_proto = out.File
	file_waypoint_builtin_aws_lambda_function_url_plugin_proto_rawDesc = nil
	file_waypoint_builtin_aws_lambda_function_url_plugin_proto_goTypes = nil
	file_waypoint_builtin_aws_lambda_function_url_plugin_proto_depIdxs = nil
}
