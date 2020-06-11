package main

import (
	"testing"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func ptrStr(v string) *string { return &v }
func ptrInt(v int32) *int32   { return &v }

func getTestEnums() []*descriptor.EnumValueDescriptorProto {
	return []*descriptor.EnumValueDescriptorProto{
		{
			Name:   ptrStr("test01"),
			Number: ptrInt(1),
		},
		{
			Name:   ptrStr("test02"),
			Number: ptrInt(2),
		},
	}
}

func TestCsvEnumFormatter_printContent(t *testing.T) {
	formatter := csvEnumFormatter{}
	entries := []*descriptor.EnumDescriptorProto{
		{
			Name:  ptrStr("test"),
			Value: getTestEnums(),
		},
	}

	actual := formatter.printContent("", entries)
	expected := `1,test01,test
2,test02,test
`

	if actual != expected {
		t.Errorf("actual: %v, expected: %v", actual, expected)
	}
}

func TestJsonlEnumFormatter_printContent(t *testing.T) {
	formatter := jsonlEnumFormatter{}
	entries := []*descriptor.EnumDescriptorProto{
		{
			Name:  ptrStr("test"),
			Value: getTestEnums(),
		},
	}

	actual := formatter.printContent("", entries)
	expected := `{"number": 1, "name": "test01", "message_name": "test"}
{"number": 2, "name": "test02", "message_name": "test"}
`

	if actual != expected {
		t.Errorf("actual: %v, expected: %v", actual, expected)
	}
}

func TestSqlEnumFormatter_printContent(t *testing.T) {
	formatter := sqlEnumFormatter{}
	entries := []*descriptor.EnumDescriptorProto{
		{
			Name:  ptrStr("test"),
			Value: getTestEnums(),
		},
	}

	actual := formatter.printContent("test_table", entries)
	expected := `CREATE TABLE IF NOT EXISTS test_table (
	number BIGINT UNSIGNED NOT NULL,
	name VARCHAR(255) NOT NULL,
	message_name VARCHAR(255) NOT NULL
);
INSERT INTO test_table (number, name, message_name) VALUES (1, "test01", "test"); 
INSERT INTO test_table (number, name, message_name) VALUES (2, "test02", "test"); 
`

	if actual != expected {
		t.Errorf("actual: %v, expected: %v", actual, expected)
	}
}
