package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

var formatters = map[string]enumFormatter{
	"csv":   &csvEnumFormatter{},
	"jsonl": &jsonlEnumFormatter{},
	"sql":   &sqlEnumFormatter{},
}

func merge(l, r map[string]*descriptor.EnumDescriptorProto) map[string]*descriptor.EnumDescriptorProto {
	merged := make(map[string]*descriptor.EnumDescriptorProto)

	for k, v := range l {
		merged[k] = v
	}
	for k, v := range r {
		merged[k] = v // NOTE: Overwrite if value appears
	}

	return merged
}

func appendNestedEnum(parent string, desc []*descriptor.DescriptorProto) map[string]*descriptor.EnumDescriptorProto {
	entries := make(map[string]*descriptor.EnumDescriptorProto)

	for _, d := range desc {
		current := parent + "_" + d.GetName()
		for _, e := range d.GetEnumType() {
			fqdn := current + "_" + e.GetName()
			entries[fqdn] = e
		}
		entries = merge(entries, appendNestedEnum(current, d.GetNestedType()))
	}

	return entries
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

	entries := make(map[string]*descriptor.EnumDescriptorProto)
	for _, f := range req.GetProtoFile() {
		current := strings.ReplaceAll(f.GetPackage(), ".", "_")
		for _, e := range f.GetEnumType() {
			fqdn := current + "_" + e.GetName()
			entries[fqdn] = e
		}
		entries = merge(entries, appendNestedEnum(current, f.GetMessageType()))
	}

	resp := plugin.CodeGeneratorResponse{}
	for descName, contentEntries := range entries {
		resp.File = append(resp.File, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(descName + formatter.extension()),
			Content: proto.String(formatter.printContent(descName, contentEntries)),
		})
	}

	buf, err = proto.Marshal(&resp)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stdout.Write(buf); err != nil {
		log.Fatal(err)
	}
}
