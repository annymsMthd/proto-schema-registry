package protobuf_test

import (
	"testing"

	"github.com/syncromatics/proto-schema-registry/internal/testing/testProto/gen"
	v1 "github.com/syncromatics/proto-schema-registry/internal/testing/testProto/v1"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/syncromatics/proto-schema-registry/pkg/protobuf"
	"gotest.tools/assert"
)

func Test_SchemaExtractor(t *testing.T) {
	o := &v1.TestObjects{}
	s, err := protobuf.ExtractSchema(o)
	if err != nil {
		t.Fatal(err)
	}

	expected := `syntax = "proto3";
package gen;
message google_protobuf_StringValue {
	string value = 1;
}
message v1_TestObject {
	int64 id = 1;
	google_protobuf_StringValue stuff = 2;
	bool bool_message = 3;
	bytes bytes_message = 4;
	double double_message = 5;
	fixed32 fixed32_message = 6;
	fixed64 fixed64_message = 7;
	float float_message = 8;
	sfixed32 sfixed32_message = 9;
	sfixed64 sfixed64_message = 10;
	sint32 sint32_message = 11;
	sint64 sint64_message = 12;
	uint32 uint32_message = 13;
	uint64 uint64_message = 14;
}
message v1_TestObject2_MapObjectsEntry {
	string key = 1;
	v1_TestObject value = 2;
}
message v1_TestObject2_StringMapsEntry {
	int32 key = 1;
	string value = 2;
}
message v1_TestObject2 {
	repeated google_protobuf_StringValue strings = 1;
	repeated v1_TestObject2_MapObjectsEntry map_objects = 2;
	repeated v1_TestObject2_StringMapsEntry string_maps = 3;
}
message v1_TestObject3 {
	oneof oneof_0 {
		v1_TestObject request_test_object = 1;
		v1_TestObject2 request_test_object_2 = 2;
		string request_string = 4;
	}
	string bla = 3;
	oneof oneof_1 {
		string request2_string = 5;
		int32 request2_int32 = 6;
	}
}
enum v1_Enum1 {
	ZERO = 0;
	ONE = 1;
	TWO = 2;
	THREE = 3;
	FOUR = 4;
}
message v1_EnumMessage {
	v1_Enum1 enum1_message = 1;
	Enum2 enum2_message = 2;
	enum Enum2 {
		STUFF = 0;
		PIE = 1;
	}
}
message record {
	v1_TestObject object = 1;
	v1_TestObject2 object_2 = 2;
	v1_TestObject3 object_3 = 3;
	v1_EnumMessage enum_message = 4;
	v1_EnumMessage.Enum2 enum_inside_message = 5;
}
`

	assert.Equal(t, expected, s)
}

func Test_UnmarshalWithGen(t *testing.T) {
	message := &v1.TestObjects{
		Object: &v1.TestObject{
			Id: 51345,
			Stuff: &wrappers.StringValue{
				Value: "hey yo",
			},
			BoolMessage:     true,
			BytesMessage:    []byte{0x1, 0x2},
			DoubleMessage:   7.8,
			Fixed32Message:  6,
			Fixed64Message:  89,
			FloatMessage:    9.0776,
			Sfixed32Message: 34,
			Sfixed64Message: 234,
			Sint32Message:   12,
			Sint64Message:   32,
		},
		Object_2: &v1.TestObject2{
			Strings: []*wrappers.StringValue{
				&wrappers.StringValue{
					Value: "item1",
				},
				&wrappers.StringValue{
					Value: "item2",
				},
			},
			StringMaps: map[int32]string{
				4: "stuff",
			},
			MapObjects: map[string]*v1.TestObject{
				"object1": &v1.TestObject{},
			},
		},
		Object_3: &v1.TestObject3{
			Bla: "bla",
			Request: &v1.TestObject3_RequestTestObject{
				RequestTestObject: &v1.TestObject{},
			},
			Request2: &v1.TestObject3_Request2Int32{
				Request2Int32: 54,
			},
		},
		EnumMessage: &v1.EnumMessage{
			Enum1Message: v1.Enum1_FOUR,
			Enum2Message: v1.EnumMessage_PIE,
		},
		EnumInsideMessage: v1.EnumMessage_STUFF,
	}

	bytes, err := proto.Marshal(message)
	if err != nil {
		t.Fatal(err)
	}

	new := &gen.Record{}
	err = proto.Unmarshal(bytes, new)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, new.String(), message.String())
}

func Test_SchemaExtractor_ShouldHandleRecursiveMessage(t *testing.T) {
	o := &v1.AnyValue{}
	s, err := protobuf.ExtractSchema(o)
	if err != nil {
		t.Fatal(err)
	}

	expected := `syntax = "proto3";
package gen;
message record {
	oneof oneof_0 {
		v1_ArrayValue array_value = 1;
		string string_value = 2;
	}
}
message v1_ArrayValue {
	repeated v1_AnyValue values = 1;
}
`

	assert.Equal(t, expected, s)
}
