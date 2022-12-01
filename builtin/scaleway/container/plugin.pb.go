// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.7
// source: builtin/scaleway/container/plugin.proto

package container

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

type Container struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            string         `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string         `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Image         string         `protobuf:"bytes,3,opt,name=image,proto3" json:"image,omitempty"`
	Url           string         `protobuf:"bytes,4,opt,name=url,proto3" json:"url,omitempty"`
	Region        string         `protobuf:"bytes,5,opt,name=region,proto3" json:"region,omitempty"`
	DeploymentId  string         `protobuf:"bytes,6,opt,name=deployment_id,json=deploymentId,proto3" json:"deployment_id,omitempty"`
	ResourceState *opaqueany.Any `protobuf:"bytes,7,opt,name=resource_state,json=resourceState,proto3" json:"resource_state,omitempty"`
}

func (x *Container) Reset() {
	*x = Container{}
	if protoimpl.UnsafeEnabled {
		mi := &file_builtin_scaleway_container_plugin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Container) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Container) ProtoMessage() {}

func (x *Container) ProtoReflect() protoreflect.Message {
	mi := &file_builtin_scaleway_container_plugin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Container.ProtoReflect.Descriptor instead.
func (*Container) Descriptor() ([]byte, []int) {
	return file_builtin_scaleway_container_plugin_proto_rawDescGZIP(), []int{0}
}

func (x *Container) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Container) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Container) GetImage() string {
	if x != nil {
		return x.Image
	}
	return ""
}

func (x *Container) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Container) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *Container) GetDeploymentId() string {
	if x != nil {
		return x.DeploymentId
	}
	return ""
}

func (x *Container) GetResourceState() *opaqueany.Any {
	if x != nil {
		return x.ResourceState
	}
	return nil
}

type Resource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Resource) Reset() {
	*x = Resource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_builtin_scaleway_container_plugin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Resource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Resource) ProtoMessage() {}

func (x *Resource) ProtoReflect() protoreflect.Message {
	mi := &file_builtin_scaleway_container_plugin_proto_msgTypes[1]
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
	return file_builtin_scaleway_container_plugin_proto_rawDescGZIP(), []int{1}
}

type Resource_Container struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Region string `protobuf:"bytes,2,opt,name=region,proto3" json:"region,omitempty"`
}

func (x *Resource_Container) Reset() {
	*x = Resource_Container{}
	if protoimpl.UnsafeEnabled {
		mi := &file_builtin_scaleway_container_plugin_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Resource_Container) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Resource_Container) ProtoMessage() {}

func (x *Resource_Container) ProtoReflect() protoreflect.Message {
	mi := &file_builtin_scaleway_container_plugin_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Resource_Container.ProtoReflect.Descriptor instead.
func (*Resource_Container) Descriptor() ([]byte, []int) {
	return file_builtin_scaleway_container_plugin_proto_rawDescGZIP(), []int{1, 0}
}

func (x *Resource_Container) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Resource_Container) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

var File_builtin_scaleway_container_plugin_proto protoreflect.FileDescriptor

var file_builtin_scaleway_container_plugin_proto_rawDesc = []byte{
	0x0a, 0x27, 0x62, 0x75, 0x69, 0x6c, 0x74, 0x69, 0x6e, 0x2f, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x77,
	0x61, 0x79, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x2f, 0x70, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x73, 0x63, 0x61, 0x6c, 0x65,
	0x77, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x1a, 0x09, 0x61,
	0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcb, 0x01, 0x0a, 0x09, 0x43, 0x6f, 0x6e,
	0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6d,
	0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75,
	0x72, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x64, 0x65,
	0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12,
	0x35, 0x0a, 0x0e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6f, 0x70, 0x61, 0x71, 0x75, 0x65,
	0x61, 0x6e, 0x79, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x22, 0x3f, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x1a, 0x33, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x42, 0x25, 0x5a, 0x23, 0x77, 0x61, 0x79, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x74, 0x69, 0x6e, 0x2f, 0x73, 0x63, 0x61, 0x6c,
	0x65, 0x77, 0x61, 0x79, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_builtin_scaleway_container_plugin_proto_rawDescOnce sync.Once
	file_builtin_scaleway_container_plugin_proto_rawDescData = file_builtin_scaleway_container_plugin_proto_rawDesc
)

func file_builtin_scaleway_container_plugin_proto_rawDescGZIP() []byte {
	file_builtin_scaleway_container_plugin_proto_rawDescOnce.Do(func() {
		file_builtin_scaleway_container_plugin_proto_rawDescData = protoimpl.X.CompressGZIP(file_builtin_scaleway_container_plugin_proto_rawDescData)
	})
	return file_builtin_scaleway_container_plugin_proto_rawDescData
}

var file_builtin_scaleway_container_plugin_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_builtin_scaleway_container_plugin_proto_goTypes = []interface{}{
	(*Container)(nil),          // 0: scaleway.container.Container
	(*Resource)(nil),           // 1: scaleway.container.Resource
	(*Resource_Container)(nil), // 2: scaleway.container.Resource.Container
	(*opaqueany.Any)(nil),      // 3: opaqueany.Any
}
var file_builtin_scaleway_container_plugin_proto_depIdxs = []int32{
	3, // 0: scaleway.container.Container.resource_state:type_name -> opaqueany.Any
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_builtin_scaleway_container_plugin_proto_init() }
func file_builtin_scaleway_container_plugin_proto_init() {
	if File_builtin_scaleway_container_plugin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_builtin_scaleway_container_plugin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Container); i {
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
		file_builtin_scaleway_container_plugin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_builtin_scaleway_container_plugin_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Resource_Container); i {
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
			RawDescriptor: file_builtin_scaleway_container_plugin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_builtin_scaleway_container_plugin_proto_goTypes,
		DependencyIndexes: file_builtin_scaleway_container_plugin_proto_depIdxs,
		MessageInfos:      file_builtin_scaleway_container_plugin_proto_msgTypes,
	}.Build()
	File_builtin_scaleway_container_plugin_proto = out.File
	file_builtin_scaleway_container_plugin_proto_rawDesc = nil
	file_builtin_scaleway_container_plugin_proto_goTypes = nil
	file_builtin_scaleway_container_plugin_proto_depIdxs = nil
}
