// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.6
// source: ops_enums.proto

package ops

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

type OpsOrderState int32

const (
	OpsOrderState_OPS_ORDER_STATE_NEW         OpsOrderState = 0
	OpsOrderState_OPS_ORDER_STATE_APPROVED    OpsOrderState = 1
	OpsOrderState_OPS_ORDER_STATE_IN_PROCESS  OpsOrderState = 2
	OpsOrderState_OPS_ORDER_STATE_PART_FILLED OpsOrderState = 3
	OpsOrderState_OPS_ORDER_STATE_FILLED      OpsOrderState = 4
	OpsOrderState_OPS_ORDER_STATE_DONE        OpsOrderState = 5
	OpsOrderState_OPS_ORDER_STATE_REJECTED    OpsOrderState = 6
)

// Enum value maps for OpsOrderState.
var (
	OpsOrderState_name = map[int32]string{
		0: "OPS_ORDER_STATE_NEW",
		1: "OPS_ORDER_STATE_APPROVED",
		2: "OPS_ORDER_STATE_IN_PROCESS",
		3: "OPS_ORDER_STATE_PART_FILLED",
		4: "OPS_ORDER_STATE_FILLED",
		5: "OPS_ORDER_STATE_DONE",
		6: "OPS_ORDER_STATE_REJECTED",
	}
	OpsOrderState_value = map[string]int32{
		"OPS_ORDER_STATE_NEW":         0,
		"OPS_ORDER_STATE_APPROVED":    1,
		"OPS_ORDER_STATE_IN_PROCESS":  2,
		"OPS_ORDER_STATE_PART_FILLED": 3,
		"OPS_ORDER_STATE_FILLED":      4,
		"OPS_ORDER_STATE_DONE":        5,
		"OPS_ORDER_STATE_REJECTED":    6,
	}
)

func (x OpsOrderState) Enum() *OpsOrderState {
	p := new(OpsOrderState)
	*p = x
	return p
}

func (x OpsOrderState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OpsOrderState) Descriptor() protoreflect.EnumDescriptor {
	return file_ops_enums_proto_enumTypes[0].Descriptor()
}

func (OpsOrderState) Type() protoreflect.EnumType {
	return &file_ops_enums_proto_enumTypes[0]
}

func (x OpsOrderState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OpsOrderState.Descriptor instead.
func (OpsOrderState) EnumDescriptor() ([]byte, []int) {
	return file_ops_enums_proto_rawDescGZIP(), []int{0}
}

type OpsOrderType int32

const (
	OpsOrderType_OPS_ORDER_TYPE_MARKET OpsOrderType = 0
	OpsOrderType_OPS_ORDER_TYPE_LIMIT  OpsOrderType = 1
)

// Enum value maps for OpsOrderType.
var (
	OpsOrderType_name = map[int32]string{
		0: "OPS_ORDER_TYPE_MARKET",
		1: "OPS_ORDER_TYPE_LIMIT",
	}
	OpsOrderType_value = map[string]int32{
		"OPS_ORDER_TYPE_MARKET": 0,
		"OPS_ORDER_TYPE_LIMIT":  1,
	}
)

func (x OpsOrderType) Enum() *OpsOrderType {
	p := new(OpsOrderType)
	*p = x
	return p
}

func (x OpsOrderType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OpsOrderType) Descriptor() protoreflect.EnumDescriptor {
	return file_ops_enums_proto_enumTypes[1].Descriptor()
}

func (OpsOrderType) Type() protoreflect.EnumType {
	return &file_ops_enums_proto_enumTypes[1]
}

func (x OpsOrderType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OpsOrderType.Descriptor instead.
func (OpsOrderType) EnumDescriptor() ([]byte, []int) {
	return file_ops_enums_proto_rawDescGZIP(), []int{1}
}

type OpsOrderDirection int32

const (
	OpsOrderDirection_OPS_ORDER_DIRECTION_SELL OpsOrderDirection = 0
	OpsOrderDirection_OPS_ORDER_DIRECTION_BUY  OpsOrderDirection = 1
)

// Enum value maps for OpsOrderDirection.
var (
	OpsOrderDirection_name = map[int32]string{
		0: "OPS_ORDER_DIRECTION_SELL",
		1: "OPS_ORDER_DIRECTION_BUY",
	}
	OpsOrderDirection_value = map[string]int32{
		"OPS_ORDER_DIRECTION_SELL": 0,
		"OPS_ORDER_DIRECTION_BUY":  1,
	}
)

func (x OpsOrderDirection) Enum() *OpsOrderDirection {
	p := new(OpsOrderDirection)
	*p = x
	return p
}

func (x OpsOrderDirection) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OpsOrderDirection) Descriptor() protoreflect.EnumDescriptor {
	return file_ops_enums_proto_enumTypes[2].Descriptor()
}

func (OpsOrderDirection) Type() protoreflect.EnumType {
	return &file_ops_enums_proto_enumTypes[2]
}

func (x OpsOrderDirection) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OpsOrderDirection.Descriptor instead.
func (OpsOrderDirection) EnumDescriptor() ([]byte, []int) {
	return file_ops_enums_proto_rawDescGZIP(), []int{2}
}

var File_ops_enums_proto protoreflect.FileDescriptor

var file_ops_enums_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x6f, 0x70, 0x73, 0x5f, 0x65, 0x6e, 0x75, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x03, 0x4f, 0x50, 0x53, 0x2a, 0xdb, 0x01, 0x0a, 0x0d, 0x4f, 0x70, 0x73, 0x4f, 0x72,
	0x64, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x17, 0x0a, 0x13, 0x4f, 0x50, 0x53, 0x5f,
	0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x4e, 0x45, 0x57, 0x10,
	0x00, 0x12, 0x1c, 0x0a, 0x18, 0x4f, 0x50, 0x53, 0x5f, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x53,
	0x54, 0x41, 0x54, 0x45, 0x5f, 0x41, 0x50, 0x50, 0x52, 0x4f, 0x56, 0x45, 0x44, 0x10, 0x01, 0x12,
	0x1e, 0x0a, 0x1a, 0x4f, 0x50, 0x53, 0x5f, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41,
	0x54, 0x45, 0x5f, 0x49, 0x4e, 0x5f, 0x50, 0x52, 0x4f, 0x43, 0x45, 0x53, 0x53, 0x10, 0x02, 0x12,
	0x1f, 0x0a, 0x1b, 0x4f, 0x50, 0x53, 0x5f, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41,
	0x54, 0x45, 0x5f, 0x50, 0x41, 0x52, 0x54, 0x5f, 0x46, 0x49, 0x4c, 0x4c, 0x45, 0x44, 0x10, 0x03,
	0x12, 0x1a, 0x0a, 0x16, 0x4f, 0x50, 0x53, 0x5f, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x53, 0x54,
	0x41, 0x54, 0x45, 0x5f, 0x46, 0x49, 0x4c, 0x4c, 0x45, 0x44, 0x10, 0x04, 0x12, 0x18, 0x0a, 0x14,
	0x4f, 0x50, 0x53, 0x5f, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f,
	0x44, 0x4f, 0x4e, 0x45, 0x10, 0x05, 0x12, 0x1c, 0x0a, 0x18, 0x4f, 0x50, 0x53, 0x5f, 0x4f, 0x52,
	0x44, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x52, 0x45, 0x4a, 0x45, 0x43, 0x54,
	0x45, 0x44, 0x10, 0x06, 0x2a, 0x43, 0x0a, 0x0c, 0x4f, 0x70, 0x73, 0x4f, 0x72, 0x64, 0x65, 0x72,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x19, 0x0a, 0x15, 0x4f, 0x50, 0x53, 0x5f, 0x4f, 0x52, 0x44, 0x45,
	0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4d, 0x41, 0x52, 0x4b, 0x45, 0x54, 0x10, 0x00, 0x12,
	0x18, 0x0a, 0x14, 0x4f, 0x50, 0x53, 0x5f, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x4c, 0x49, 0x4d, 0x49, 0x54, 0x10, 0x01, 0x2a, 0x4e, 0x0a, 0x11, 0x4f, 0x70, 0x73,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1c,
	0x0a, 0x18, 0x4f, 0x50, 0x53, 0x5f, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x44, 0x49, 0x52, 0x45,
	0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x53, 0x45, 0x4c, 0x4c, 0x10, 0x00, 0x12, 0x1b, 0x0a, 0x17,
	0x4f, 0x50, 0x53, 0x5f, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x44, 0x49, 0x52, 0x45, 0x43, 0x54,
	0x49, 0x4f, 0x4e, 0x5f, 0x42, 0x55, 0x59, 0x10, 0x01, 0x42, 0x06, 0x5a, 0x04, 0x2f, 0x6f, 0x70,
	0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ops_enums_proto_rawDescOnce sync.Once
	file_ops_enums_proto_rawDescData = file_ops_enums_proto_rawDesc
)

func file_ops_enums_proto_rawDescGZIP() []byte {
	file_ops_enums_proto_rawDescOnce.Do(func() {
		file_ops_enums_proto_rawDescData = protoimpl.X.CompressGZIP(file_ops_enums_proto_rawDescData)
	})
	return file_ops_enums_proto_rawDescData
}

var file_ops_enums_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_ops_enums_proto_goTypes = []interface{}{
	(OpsOrderState)(0),     // 0: OPS.OpsOrderState
	(OpsOrderType)(0),      // 1: OPS.OpsOrderType
	(OpsOrderDirection)(0), // 2: OPS.OpsOrderDirection
}
var file_ops_enums_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_ops_enums_proto_init() }
func file_ops_enums_proto_init() {
	if File_ops_enums_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ops_enums_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ops_enums_proto_goTypes,
		DependencyIndexes: file_ops_enums_proto_depIdxs,
		EnumInfos:         file_ops_enums_proto_enumTypes,
	}.Build()
	File_ops_enums_proto = out.File
	file_ops_enums_proto_rawDesc = nil
	file_ops_enums_proto_goTypes = nil
	file_ops_enums_proto_depIdxs = nil
}
