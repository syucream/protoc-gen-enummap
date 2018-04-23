package main

import (
	"fmt"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type enumFormatter interface {
	printHeader(string) string
	printLine(*descriptor.EnumValueDescriptorProto) string
	extension() string
}

type csvEnumFormatter struct{}

func (f *csvEnumFormatter) printHeader(criteria string) string { return "" }

func (f *csvEnumFormatter) printLine(ev *descriptor.EnumValueDescriptorProto) string {
	return fmt.Sprintf("%d,%s\n", ev.GetNumber(), ev.GetName())
}

func (f *csvEnumFormatter) extension() string { return ".csv" }

type jsonlEnumFormatter struct{}

func (f *jsonlEnumFormatter) printHeader(criteria string) string { return "" }

func (f *jsonlEnumFormatter) printLine(ev *descriptor.EnumValueDescriptorProto) string {
	return fmt.Sprintf("{\"number\": %d, \"name\": \"%s\"}\n", ev.GetNumber(), ev.GetName())
}

func (f *jsonlEnumFormatter) extension() string { return ".jsonl" }

type sqlEnumFormatter struct{}

func (f *sqlEnumFormatter) printHeader(criteria string) string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	number BIGINT UNSIGNED NOT NULL,
	name VARCHAR(64) NOT NULL
);
`, criteria)
}

func (f *sqlEnumFormatter) printLine(ev *descriptor.EnumValueDescriptorProto) string {
	return fmt.Sprintf("INSERT INTO exported (number, name) VALUES (%d, \"%s\"); \n", ev.GetNumber(), ev.GetName())
}

func (f *sqlEnumFormatter) extension() string { return ".sql" }
