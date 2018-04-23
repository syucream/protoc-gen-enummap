package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

var formatters = map[string]enumFormatter{
	"csv":   &csvEnumFormatter{},
	"jsonl": &jsonlEnumFormatter{},
}

func appendNestedEnum(file []*plugin.CodeGeneratorResponse_File, formatter enumFormatter, prefix string, desc []*descriptor.DescriptorProto) []*plugin.CodeGeneratorResponse_File {
	for _, d := range desc {
		for _, e := range d.GetEnumType() {
			var contents []string
			for _, ev := range e.GetValue() {
				contents = append(contents, formatter.printLine(ev))
			}
			file = append(file, &plugin.CodeGeneratorResponse_File{
				Name:    proto.String(prefix + d.GetName() + "_" + e.GetName() + formatter.extension()),
				Content: proto.String(strings.Join(contents, "")),
			})
		}

		file = appendNestedEnum(file, formatter, prefix+d.GetName()+"_", d.GetNestedType())
	}

	return file
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
		filename := strings.Replace(f.GetPackage(), ".", "_", -1) + "__"

		for _, e := range f.GetEnumType() {
			var contents []string
			for _, ev := range e.GetValue() {
				contents = append(contents, formatter.printLine(ev))
			}
			resp.File = append(resp.File, &plugin.CodeGeneratorResponse_File{
				Name:    proto.String(filename + e.GetName() + formatter.extension()),
				Content: proto.String(strings.Join(contents, "")),
			})
		}
		resp.File = appendNestedEnum(resp.File, formatter, filename, f.GetMessageType())
	}

	buf, err = proto.Marshal(&resp)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stdout.Write(buf); err != nil {
		log.Fatal(err)
	}
}
