package main

import (
	"fmt"
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type enumFormatter interface {
	printContent(string, []*descriptor.EnumValueDescriptorProto) string
	extension() string
}

type csvEnumFormatter struct{}

func (f *csvEnumFormatter) printContent(criteria string, evs []*descriptor.EnumValueDescriptorProto) string {
	var contents []string

	for _, ev := range evs {
		content := fmt.Sprintf("%d,%s\n", ev.GetNumber(), ev.GetName())
		contents = append(contents, content)
	}

	return strings.Join(contents, "")
}

func (f *csvEnumFormatter) extension() string { return ".csv" }

type jsonlEnumFormatter struct{}

func (f *jsonlEnumFormatter) printContent(criteria string, evs []*descriptor.EnumValueDescriptorProto) string {
	var contents []string

	for _, ev := range evs {
		content := fmt.Sprintf("{\"number\": %d, \"name\": \"%s\"}\n", ev.GetNumber(), ev.GetName())
		contents = append(contents, content)
	}

	return strings.Join(contents, "")
}

func (f *jsonlEnumFormatter) extension() string { return ".jsonl" }

type sqlEnumFormatter struct{}

func (f *sqlEnumFormatter) printContent(criteria string, evs []*descriptor.EnumValueDescriptorProto) string {
	var contents []string

	contents = append(contents, fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	number BIGINT UNSIGNED NOT NULL,
	name VARCHAR(64) NOT NULL
);
`, criteria))

	for _, ev := range evs {
		content := fmt.Sprintf("INSERT INTO %s (number, name) VALUES (%d, \"%s\"); \n", criteria, ev.GetNumber(), ev.GetName())
		contents = append(contents, content)
	}

	return strings.Join(contents, "")
}

func (f *sqlEnumFormatter) extension() string { return ".sql" }
