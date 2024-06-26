// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.3
// source: metric.proto

package proto

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

type UpdatesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*Metric `protobuf:"bytes,1,rep,name=Metrics,proto3" json:"Metrics,omitempty"`
}

func (x *UpdatesRequest) Reset() {
	*x = UpdatesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metric_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdatesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatesRequest) ProtoMessage() {}

func (x *UpdatesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_metric_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatesRequest.ProtoReflect.Descriptor instead.
func (*UpdatesRequest) Descriptor() ([]byte, []int) {
	return file_metric_proto_rawDescGZIP(), []int{0}
}

func (x *UpdatesRequest) GetMetrics() []*Metric {
	if x != nil {
		return x.Metrics
	}
	return nil
}

type UpdatesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *UpdatesResponse) Reset() {
	*x = UpdatesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metric_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdatesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatesResponse) ProtoMessage() {}

func (x *UpdatesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_metric_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatesResponse.ProtoReflect.Descriptor instead.
func (*UpdatesResponse) Descriptor() ([]byte, []int) {
	return file_metric_proto_rawDescGZIP(), []int{1}
}

func (x *UpdatesResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MType string  `protobuf:"bytes,1,opt,name=MType,proto3" json:"MType,omitempty"`
	Delta int64   `protobuf:"varint,2,opt,name=Delta,proto3" json:"Delta,omitempty"`
	Value float64 `protobuf:"fixed64,3,opt,name=Value,proto3" json:"Value,omitempty"`
	ID    string  `protobuf:"bytes,4,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metric_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_metric_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_metric_proto_rawDescGZIP(), []int{2}
}

func (x *Metric) GetMType() string {
	if x != nil {
		return x.MType
	}
	return ""
}

func (x *Metric) GetDelta() int64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

func (x *Metric) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *Metric) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

var File_metric_proto protoreflect.FileDescriptor

var file_metric_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x22, 0x3f, 0x0a, 0x0e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2d, 0x0a,
	0x07, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13,
	0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x52, 0x07, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x22, 0x27, 0x0a, 0x0f,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x5a, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12,
	0x14, 0x0a, 0x05, 0x4d, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x4d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x44, 0x65, 0x6c, 0x74, 0x61, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x44, 0x65, 0x6c, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49,
	0x44, 0x32, 0x52, 0x0a, 0x0a, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x47, 0x72, 0x70, 0x63, 0x12,
	0x44, 0x0a, 0x07, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x12, 0x1b, 0x2e, 0x6d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0b, 0x5a, 0x09, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_metric_proto_rawDescOnce sync.Once
	file_metric_proto_rawDescData = file_metric_proto_rawDesc
)

func file_metric_proto_rawDescGZIP() []byte {
	file_metric_proto_rawDescOnce.Do(func() {
		file_metric_proto_rawDescData = protoimpl.X.CompressGZIP(file_metric_proto_rawDescData)
	})
	return file_metric_proto_rawDescData
}

var file_metric_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_metric_proto_goTypes = []interface{}{
	(*UpdatesRequest)(nil),  // 0: metric_grpc.UpdatesRequest
	(*UpdatesResponse)(nil), // 1: metric_grpc.UpdatesResponse
	(*Metric)(nil),          // 2: metric_grpc.Metric
}
var file_metric_proto_depIdxs = []int32{
	2, // 0: metric_grpc.UpdatesRequest.Metrics:type_name -> metric_grpc.Metric
	0, // 1: metric_grpc.metricGrpc.Updates:input_type -> metric_grpc.UpdatesRequest
	1, // 2: metric_grpc.metricGrpc.Updates:output_type -> metric_grpc.UpdatesResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_metric_proto_init() }
func file_metric_proto_init() {
	if File_metric_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_metric_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdatesRequest); i {
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
		file_metric_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdatesResponse); i {
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
		file_metric_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metric); i {
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
			RawDescriptor: file_metric_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_metric_proto_goTypes,
		DependencyIndexes: file_metric_proto_depIdxs,
		MessageInfos:      file_metric_proto_msgTypes,
	}.Build()
	File_metric_proto = out.File
	file_metric_proto_rawDesc = nil
	file_metric_proto_goTypes = nil
	file_metric_proto_depIdxs = nil
}
