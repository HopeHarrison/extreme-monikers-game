// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: monikers/proto/monikers.proto

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

type NewGameRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerName string `protobuf:"bytes,1,opt,name=player_name,json=playerName,proto3" json:"player_name,omitempty"`
}

func (x *NewGameRequest) Reset() {
	*x = NewGameRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_monikers_proto_monikers_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewGameRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewGameRequest) ProtoMessage() {}

func (x *NewGameRequest) ProtoReflect() protoreflect.Message {
	mi := &file_monikers_proto_monikers_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewGameRequest.ProtoReflect.Descriptor instead.
func (*NewGameRequest) Descriptor() ([]byte, []int) {
	return file_monikers_proto_monikers_proto_rawDescGZIP(), []int{0}
}

func (x *NewGameRequest) GetPlayerName() string {
	if x != nil {
		return x.PlayerName
	}
	return ""
}

type NewGameResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GameId string `protobuf:"bytes,1,opt,name=game_id,json=gameId,proto3" json:"game_id,omitempty"`
}

func (x *NewGameResponse) Reset() {
	*x = NewGameResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_monikers_proto_monikers_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewGameResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewGameResponse) ProtoMessage() {}

func (x *NewGameResponse) ProtoReflect() protoreflect.Message {
	mi := &file_monikers_proto_monikers_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewGameResponse.ProtoReflect.Descriptor instead.
func (*NewGameResponse) Descriptor() ([]byte, []int) {
	return file_monikers_proto_monikers_proto_rawDescGZIP(), []int{1}
}

func (x *NewGameResponse) GetGameId() string {
	if x != nil {
		return x.GameId
	}
	return ""
}

var File_monikers_proto_monikers_proto protoreflect.FileDescriptor

var file_monikers_proto_monikers_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x6d, 0x6f, 0x6e, 0x69, 0x6b, 0x65, 0x72, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x6d, 0x6f, 0x6e, 0x69, 0x6b, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x08, 0x6d, 0x6f, 0x6e, 0x69, 0x6b, 0x65, 0x72, 0x73, 0x22, 0x31, 0x0a, 0x0e, 0x4e, 0x65, 0x77,
	0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x70,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x2a, 0x0a, 0x0f,
	0x4e, 0x65, 0x77, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x17, 0x0a, 0x07, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x67, 0x61, 0x6d, 0x65, 0x49, 0x64, 0x32, 0x4a, 0x0a, 0x08, 0x4d, 0x6f, 0x6e, 0x69,
	0x6b, 0x65, 0x72, 0x73, 0x12, 0x3e, 0x0a, 0x07, 0x4e, 0x65, 0x77, 0x47, 0x61, 0x6d, 0x65, 0x12,
	0x18, 0x2e, 0x6d, 0x6f, 0x6e, 0x69, 0x6b, 0x65, 0x72, 0x73, 0x2e, 0x4e, 0x65, 0x77, 0x47, 0x61,
	0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x6d, 0x6f, 0x6e, 0x69,
	0x6b, 0x65, 0x72, 0x73, 0x2e, 0x4e, 0x65, 0x77, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x10, 0x5a, 0x0e, 0x6d, 0x6f, 0x6e, 0x69, 0x6b, 0x65, 0x72, 0x73,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_monikers_proto_monikers_proto_rawDescOnce sync.Once
	file_monikers_proto_monikers_proto_rawDescData = file_monikers_proto_monikers_proto_rawDesc
)

func file_monikers_proto_monikers_proto_rawDescGZIP() []byte {
	file_monikers_proto_monikers_proto_rawDescOnce.Do(func() {
		file_monikers_proto_monikers_proto_rawDescData = protoimpl.X.CompressGZIP(file_monikers_proto_monikers_proto_rawDescData)
	})
	return file_monikers_proto_monikers_proto_rawDescData
}

var file_monikers_proto_monikers_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_monikers_proto_monikers_proto_goTypes = []any{
	(*NewGameRequest)(nil),  // 0: monikers.NewGameRequest
	(*NewGameResponse)(nil), // 1: monikers.NewGameResponse
}
var file_monikers_proto_monikers_proto_depIdxs = []int32{
	0, // 0: monikers.Monikers.NewGame:input_type -> monikers.NewGameRequest
	1, // 1: monikers.Monikers.NewGame:output_type -> monikers.NewGameResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_monikers_proto_monikers_proto_init() }
func file_monikers_proto_monikers_proto_init() {
	if File_monikers_proto_monikers_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_monikers_proto_monikers_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*NewGameRequest); i {
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
		file_monikers_proto_monikers_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*NewGameResponse); i {
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
			RawDescriptor: file_monikers_proto_monikers_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_monikers_proto_monikers_proto_goTypes,
		DependencyIndexes: file_monikers_proto_monikers_proto_depIdxs,
		MessageInfos:      file_monikers_proto_monikers_proto_msgTypes,
	}.Build()
	File_monikers_proto_monikers_proto = out.File
	file_monikers_proto_monikers_proto_rawDesc = nil
	file_monikers_proto_monikers_proto_goTypes = nil
	file_monikers_proto_monikers_proto_depIdxs = nil
}
