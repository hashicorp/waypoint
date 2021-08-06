// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: waypoint/builtin/aws/ec2/plugin.proto

package ec2

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

type Deployment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServiceName    string `protobuf:"bytes,1,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	Region         string `protobuf:"bytes,2,opt,name=region,proto3" json:"region,omitempty"`
	PublicIp       string `protobuf:"bytes,3,opt,name=public_ip,json=publicIp,proto3" json:"public_ip,omitempty"`
	PublicDns      string `protobuf:"bytes,4,opt,name=public_dns,json=publicDns,proto3" json:"public_dns,omitempty"`
	TargetGroupArn string `protobuf:"bytes,5,opt,name=target_group_arn,json=targetGroupArn,proto3" json:"target_group_arn,omitempty"`
}

func (x *Deployment) Reset() {
	*x = Deployment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_waypoint_builtin_aws_ec2_plugin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Deployment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Deployment) ProtoMessage() {}

func (x *Deployment) ProtoReflect() protoreflect.Message {
	mi := &file_waypoint_builtin_aws_ec2_plugin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Deployment.ProtoReflect.Descriptor instead.
func (*Deployment) Descriptor() ([]byte, []int) {
	return file_waypoint_builtin_aws_ec2_plugin_proto_rawDescGZIP(), []int{0}
}

func (x *Deployment) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

func (x *Deployment) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *Deployment) GetPublicIp() string {
	if x != nil {
		return x.PublicIp
	}
	return ""
}

func (x *Deployment) GetPublicDns() string {
	if x != nil {
		return x.PublicDns
	}
	return ""
}

func (x *Deployment) GetTargetGroupArn() string {
	if x != nil {
		return x.TargetGroupArn
	}
	return ""
}

var File_waypoint_builtin_aws_ec2_plugin_proto protoreflect.FileDescriptor

var file_waypoint_builtin_aws_ec2_plugin_proto_rawDesc = []byte{
	0x0a, 0x25, 0x77, 0x61, 0x79, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x74,
	0x69, 0x6e, 0x2f, 0x61, 0x77, 0x73, 0x2f, 0x65, 0x63, 0x32, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x65, 0x63, 0x32, 0x22, 0xad, 0x01, 0x0a,
	0x0a, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x5f, 0x69, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x75, 0x62, 0x6c, 0x69,
	0x63, 0x49, 0x70, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x64, 0x6e,
	0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x44,
	0x6e, 0x73, 0x12, 0x28, 0x0a, 0x10, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x5f, 0x61, 0x72, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x74, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x41, 0x72, 0x6e, 0x42, 0x1a, 0x5a, 0x18,
	0x77, 0x61, 0x79, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x74, 0x69, 0x6e,
	0x2f, 0x61, 0x77, 0x73, 0x2f, 0x65, 0x63, 0x32, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_waypoint_builtin_aws_ec2_plugin_proto_rawDescOnce sync.Once
	file_waypoint_builtin_aws_ec2_plugin_proto_rawDescData = file_waypoint_builtin_aws_ec2_plugin_proto_rawDesc
)

func file_waypoint_builtin_aws_ec2_plugin_proto_rawDescGZIP() []byte {
	file_waypoint_builtin_aws_ec2_plugin_proto_rawDescOnce.Do(func() {
		file_waypoint_builtin_aws_ec2_plugin_proto_rawDescData = protoimpl.X.CompressGZIP(file_waypoint_builtin_aws_ec2_plugin_proto_rawDescData)
	})
	return file_waypoint_builtin_aws_ec2_plugin_proto_rawDescData
}

var file_waypoint_builtin_aws_ec2_plugin_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_waypoint_builtin_aws_ec2_plugin_proto_goTypes = []interface{}{
	(*Deployment)(nil), // 0: ec2.Deployment
}
var file_waypoint_builtin_aws_ec2_plugin_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_waypoint_builtin_aws_ec2_plugin_proto_init() }
func file_waypoint_builtin_aws_ec2_plugin_proto_init() {
	if File_waypoint_builtin_aws_ec2_plugin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_waypoint_builtin_aws_ec2_plugin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Deployment); i {
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
			RawDescriptor: file_waypoint_builtin_aws_ec2_plugin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_waypoint_builtin_aws_ec2_plugin_proto_goTypes,
		DependencyIndexes: file_waypoint_builtin_aws_ec2_plugin_proto_depIdxs,
		MessageInfos:      file_waypoint_builtin_aws_ec2_plugin_proto_msgTypes,
	}.Build()
	File_waypoint_builtin_aws_ec2_plugin_proto = out.File
	file_waypoint_builtin_aws_ec2_plugin_proto_rawDesc = nil
	file_waypoint_builtin_aws_ec2_plugin_proto_goTypes = nil
	file_waypoint_builtin_aws_ec2_plugin_proto_depIdxs = nil
}
