package main

import (
	"fmt"
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type enumFormatter interface {
	printContent(string, []ContentEntry) string
	extension() string
}

type csvEnumFormatter struct{}

type ContentEntry struct {
	EnumValues  []*descriptor.EnumValueDescriptorProto
	MessageName string
}

func (f *csvEnumFormatter) printContent(criteria string, entries []ContentEntry) string {
	var contents []string

	for _, c := range entries {
		for _, ev := range c.EnumValues {
			content := fmt.Sprintf("%d,%s,%s\n", ev.GetNumber(), ev.GetName(), c.MessageName)
			contents = append(contents, content)
		}
	}

	return strings.Join(contents, "")
}

func (f *csvEnumFormatter) extension() string { return ".csv" }

type jsonlEnumFormatter struct{}

func (f *jsonlEnumFormatter) printContent(criteria string, entries []ContentEntry) string {
	var contents []string

	for _, c := range entries {
		for _, ev := range c.EnumValues {
			content := fmt.Sprintf("{\"number\": %d, \"name\": \"%s\", \"message_name\": \"%s\"}\n", ev.GetNumber(), ev.GetName(), c.MessageName)
			contents = append(contents, content)
		}
	}

	return strings.Join(contents, "")
}

func (f *jsonlEnumFormatter) extension() string { return ".json" }

type sqlEnumFormatter struct{}

func (f *sqlEnumFormatter) printContent(criteria string, entries []ContentEntry) string {
	var contents []string

	contents = append(contents, fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	number BIGINT UNSIGNED NOT NULL,
	name VARCHAR(255) NOT NULL,
	message_name VARCHAR(255) NOT NULL
);
`, criteria))

	for _, c := range entries {
		for _, ev := range c.EnumValues {
			content := fmt.Sprintf("INSERT INTO %s (number, name, message_name) VALUES (%d, \"%s\", \"%s\"); \n", criteria, ev.GetNumber(), ev.GetName(), c.MessageName)
			contents = append(contents, content)
		}
	}

	return strings.Join(contents, "")
}

func (f *sqlEnumFormatter) extension() string { return ".sql" }
