package main

import (
	"fmt"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type enumFormatter interface {
	printLine(*descriptor.EnumValueDescriptorProto) string
	extension() string
}

type csvEnumFormatter struct{}

func (f *csvEnumFormatter) printLine(ev *descriptor.EnumValueDescriptorProto) string {
	return fmt.Sprintf("%d,%s\n", ev.GetNumber(), ev.GetName())
}

func (f *csvEnumFormatter) extension() string { return ".csv" }

type jsonlEnumFormatter struct{}

func (f *jsonlEnumFormatter) printLine(ev *descriptor.EnumValueDescriptorProto) string {
	return fmt.Sprintf("{\"number\": %d, \"name\": \"%s\"}\n", ev.GetNumber(), ev.GetName())
}

func (f *jsonlEnumFormatter) extension() string { return ".jsonl" }
