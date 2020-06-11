package main

import (
	"fmt"
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// enumFormatter formats enum proto's to specific format.
type enumFormatter interface {
	printContent(string, []*descriptor.EnumDescriptorProto) string
	extension() string
}

// csvEnumFormatter formats enum proto's to CSV.
type csvEnumFormatter struct{}

func (f *csvEnumFormatter) printContent(criteria string, entries []*descriptor.EnumDescriptorProto) string {
	var contents []string

	for _, e := range entries {
		for _, ev := range e.Value {
			content := fmt.Sprintf("%d,%s,%s\n", ev.GetNumber(), ev.GetName(), *e.Name)
			contents = append(contents, content)
		}
	}

	return strings.Join(contents, "")
}

func (f *csvEnumFormatter) extension() string { return ".csv" }

// jsonEnumFormatter formats enum proto's to JSON(Newline delimited).
type jsonlEnumFormatter struct{}

func (f *jsonlEnumFormatter) printContent(criteria string, entries []*descriptor.EnumDescriptorProto) string {
	var contents []string

	for _, e := range entries {
		for _, ev := range e.Value {
			content := fmt.Sprintf("{\"number\": %d, \"name\": \"%s\", \"message_name\": \"%s\"}\n", ev.GetNumber(), ev.GetName(), *e.Name)
			contents = append(contents, content)
		}
	}

	return strings.Join(contents, "")
}

func (f *jsonlEnumFormatter) extension() string { return ".json" }

// sqlEnumFormatter formats enum proto's to SQL DDL and DMLs.
type sqlEnumFormatter struct{}

func (f *sqlEnumFormatter) printContent(criteria string, entries []*descriptor.EnumDescriptorProto) string {
	var contents []string

	contents = append(contents, fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	number BIGINT UNSIGNED NOT NULL,
	name VARCHAR(255) NOT NULL,
	message_name VARCHAR(255) NOT NULL
);
`, criteria))

	for _, e := range entries {
		for _, ev := range e.Value {
			content := fmt.Sprintf("INSERT INTO %s (number, name, message_name) VALUES (%d, \"%s\", \"%s\"); \n", criteria, ev.GetNumber(), ev.GetName(), *e.Name)
			contents = append(contents, content)
		}
	}

	return strings.Join(contents, "")
}

func (f *sqlEnumFormatter) extension() string { return ".sql" }
