// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.15.8
// source: waypoint/builtin/aws/alb/plugin.proto

package alb

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type TargetGroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Arn    string `protobuf:"bytes,1,opt,name=arn,proto3" json:"arn,omitempty"`
	Region string `protobuf:"bytes,2,opt,name=region,proto3" json:"region,omitempty"`
}

func (x *TargetGroup) Reset() {
	*x = TargetGroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_waypoint_builtin_aws_alb_plugin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TargetGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TargetGroup) ProtoMessage() {}

func (x *TargetGroup) ProtoReflect() protoreflect.Message {
	mi := &file_waypoint_builtin_aws_alb_plugin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TargetGroup.ProtoReflect.Descriptor instead.
func (*TargetGroup) Descriptor() ([]byte, []int) {
	return file_waypoint_builtin_aws_alb_plugin_proto_rawDescGZIP(), []int{0}
}

func (x *TargetGroup) GetArn() string {
	if x != nil {
		return x.Arn
	}
	return ""
}

func (x *TargetGroup) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

type Release struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url             string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	LoadBalancerArn string `protobuf:"bytes,2,opt,name=load_balancer_arn,json=loadBalancerArn,proto3" json:"load_balancer_arn,omitempty"`
	ManagementLevel string `protobuf:"bytes,3,opt,name=management_level,json=managementLevel,proto3" json:"management_level,omitempty"`
	ListenerArn     string `protobuf:"bytes,4,opt,name=listener_arn,json=listenerArn,proto3" json:"listener_arn,omitempty"`
	TargetGroupArn  string `protobuf:"bytes,5,opt,name=target_group_arn,json=targetGroupArn,proto3" json:"target_group_arn,omitempty"`
	SecurityGroupId string `protobuf:"bytes,7,opt,name=security_group_id,json=securityGroupId,proto3" json:"security_group_id,omitempty"`
	ZoneId          string `protobuf:"bytes,8,opt,name=zone_id,json=zoneId,proto3" json:"zone_id,omitempty"`
	Fqdn            string `protobuf:"bytes,9,opt,name=fqdn,proto3" json:"fqdn,omitempty"`
}

func (x *Release) Reset() {
	*x = Release{}
	if protoimpl.UnsafeEnabled {
		mi := &file_waypoint_builtin_aws_alb_plugin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Release) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Release) ProtoMessage() {}

func (x *Release) ProtoReflect() protoreflect.Message {
	mi := &file_waypoint_builtin_aws_alb_plugin_proto_msgTypes[1]
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
	return file_waypoint_builtin_aws_alb_plugin_proto_rawDescGZIP(), []int{1}
}

func (x *Release) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Release) GetLoadBalancerArn() string {
	if x != nil {
		return x.LoadBalancerArn
	}
	return ""
}

func (x *Release) GetManagementLevel() string {
	if x != nil {
		return x.ManagementLevel
	}
	return ""
}

func (x *Release) GetListenerArn() string {
	if x != nil {
		return x.ListenerArn
	}
	return ""
}

func (x *Release) GetTargetGroupArn() string {
	if x != nil {
		return x.TargetGroupArn
	}
	return ""
}

func (x *Release) GetSecurityGroupId() string {
	if x != nil {
		return x.SecurityGroupId
	}
	return ""
}

func (x *Release) GetZoneId() string {
	if x != nil {
		return x.ZoneId
	}
	return ""
}

func (x *Release) GetFqdn() string {
	if x != nil {
		return x.Fqdn
	}
	return ""
}

var File_waypoint_builtin_aws_alb_plugin_proto protoreflect.FileDescriptor

var file_waypoint_builtin_aws_alb_plugin_proto_rawDesc = []byte{
	0x0a, 0x25, 0x77, 0x61, 0x79, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x74,
	0x69, 0x6e, 0x2f, 0x61, 0x77, 0x73, 0x2f, 0x61, 0x6c, 0x62, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61, 0x6c, 0x62, 0x22, 0x37, 0x0a, 0x0b,
	0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x10, 0x0a, 0x03, 0x61,
	0x72, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x61, 0x72, 0x6e, 0x12, 0x16, 0x0a,
	0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72,
	0x65, 0x67, 0x69, 0x6f, 0x6e, 0x22, 0x98, 0x02, 0x0a, 0x07, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x75, 0x72, 0x6c, 0x12, 0x2a, 0x0a, 0x11, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x62, 0x61, 0x6c, 0x61,
	0x6e, 0x63, 0x65, 0x72, 0x5f, 0x61, 0x72, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f,
	0x6c, 0x6f, 0x61, 0x64, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x72, 0x41, 0x72, 0x6e, 0x12,
	0x29, 0x0a, 0x10, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x6c, 0x65,
	0x76, 0x65, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x21, 0x0a, 0x0c, 0x6c, 0x69,
	0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x5f, 0x61, 0x72, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x41, 0x72, 0x6e, 0x12, 0x28, 0x0a,
	0x10, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x61, 0x72,
	0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x41, 0x72, 0x6e, 0x12, 0x2a, 0x0a, 0x11, 0x73, 0x65, 0x63, 0x75, 0x72,
	0x69, 0x74, 0x79, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0f, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x7a, 0x6f, 0x6e, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x7a, 0x6f, 0x6e, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x66, 0x71, 0x64, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x71, 0x64, 0x6e,
	0x42, 0x1a, 0x5a, 0x18, 0x77, 0x61, 0x79, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2f, 0x62, 0x75, 0x69,
	0x6c, 0x74, 0x69, 0x6e, 0x2f, 0x61, 0x77, 0x73, 0x2f, 0x61, 0x6c, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_waypoint_builtin_aws_alb_plugin_proto_rawDescOnce sync.Once
	file_waypoint_builtin_aws_alb_plugin_proto_rawDescData = file_waypoint_builtin_aws_alb_plugin_proto_rawDesc
)

func file_waypoint_builtin_aws_alb_plugin_proto_rawDescGZIP() []byte {
	file_waypoint_builtin_aws_alb_plugin_proto_rawDescOnce.Do(func() {
		file_waypoint_builtin_aws_alb_plugin_proto_rawDescData = protoimpl.X.CompressGZIP(file_waypoint_builtin_aws_alb_plugin_proto_rawDescData)
	})
	return file_waypoint_builtin_aws_alb_plugin_proto_rawDescData
}

var file_waypoint_builtin_aws_alb_plugin_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_waypoint_builtin_aws_alb_plugin_proto_goTypes = []interface{}{
	(*TargetGroup)(nil), // 0: alb.TargetGroup
	(*Release)(nil),     // 1: alb.Release
}
var file_waypoint_builtin_aws_alb_plugin_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_waypoint_builtin_aws_alb_plugin_proto_init() }
func file_waypoint_builtin_aws_alb_plugin_proto_init() {
	if File_waypoint_builtin_aws_alb_plugin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_waypoint_builtin_aws_alb_plugin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TargetGroup); i {
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
		file_waypoint_builtin_aws_alb_plugin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_waypoint_builtin_aws_alb_plugin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_waypoint_builtin_aws_alb_plugin_proto_goTypes,
		DependencyIndexes: file_waypoint_builtin_aws_alb_plugin_proto_depIdxs,
		MessageInfos:      file_waypoint_builtin_aws_alb_plugin_proto_msgTypes,
	}.Build()
	File_waypoint_builtin_aws_alb_plugin_proto = out.File
	file_waypoint_builtin_aws_alb_plugin_proto_rawDesc = nil
	file_waypoint_builtin_aws_alb_plugin_proto_goTypes = nil
	file_waypoint_builtin_aws_alb_plugin_proto_depIdxs = nil
}
