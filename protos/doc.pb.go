// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v5.29.3
// source: doc.proto

package protos

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

type Doc struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Doc           string                 `protobuf:"bytes,1,opt,name=doc,proto3" json:"doc,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Doc) Reset() {
	*x = Doc{}
	mi := &file_doc_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Doc) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Doc) ProtoMessage() {}

func (x *Doc) ProtoReflect() protoreflect.Message {
	mi := &file_doc_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Doc.ProtoReflect.Descriptor instead.
func (*Doc) Descriptor() ([]byte, []int) {
	return file_doc_proto_rawDescGZIP(), []int{0}
}

func (x *Doc) GetDoc() string {
	if x != nil {
		return x.Doc
	}
	return ""
}

var File_doc_proto protoreflect.FileDescriptor

var file_doc_proto_rawDesc = []byte{
	0x0a, 0x09, 0x64, 0x6f, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x73, 0x22, 0x17, 0x0a, 0x03, 0x44, 0x6f, 0x63, 0x12, 0x10, 0x0a, 0x03, 0x64, 0x6f,
	0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x64, 0x6f, 0x63, 0x42, 0x0a, 0x5a, 0x08,
	0x2e, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_doc_proto_rawDescOnce sync.Once
	file_doc_proto_rawDescData = file_doc_proto_rawDesc
)

func file_doc_proto_rawDescGZIP() []byte {
	file_doc_proto_rawDescOnce.Do(func() {
		file_doc_proto_rawDescData = protoimpl.X.CompressGZIP(file_doc_proto_rawDescData)
	})
	return file_doc_proto_rawDescData
}

var file_doc_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_doc_proto_goTypes = []any{
	(*Doc)(nil), // 0: protos.Doc
}
var file_doc_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_doc_proto_init() }
func file_doc_proto_init() {
	if File_doc_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_doc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_doc_proto_goTypes,
		DependencyIndexes: file_doc_proto_depIdxs,
		MessageInfos:      file_doc_proto_msgTypes,
	}.Build()
	File_doc_proto = out.File
	file_doc_proto_rawDesc = nil
	file_doc_proto_goTypes = nil
	file_doc_proto_depIdxs = nil
}
