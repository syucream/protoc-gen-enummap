package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
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

var formatters = map[string]enumFormatter{
	"csv":   &csvEnumFormatter{},
	"jsonl": &jsonlEnumFormatter{},
}

func main() {
	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	var req plugin.CodeGeneratorRequest
	if err := proto.Unmarshal(buf, &req); err != nil {
		log.Fatal(err)
	}

	formatter, fmtOk := formatters[req.GetParameter()]
	if !fmtOk {
		log.Fatal("Specify supported format by --enummap_opt=")
	}

	resp := plugin.CodeGeneratorResponse{}
	for _, f := range req.GetProtoFile() {
		for _, e := range f.GetEnumType() {
			var contents []string
			for _, ev := range e.GetValue() {
				contents = append(contents, formatter.printLine(ev))
			}
			resp.File = append(resp.File, &plugin.CodeGeneratorResponse_File{
				Name:    proto.String(e.GetName() + formatter.extension()),
				Content: proto.String(strings.Join(contents, "")),
			})
		}
	}

	buf, err = proto.Marshal(&resp)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stdout.Write(buf); err != nil {
		log.Fatal(err)
	}
}
