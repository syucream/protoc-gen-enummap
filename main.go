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

func getNested(d interface{}) []*descriptor.DescriptorProto {
	switch dp := d.(type) {
	case *descriptor.FileDescriptorProto:
		return dp.GetMessageType()
	case *descriptor.DescriptorProto:
		return dp.GetNestedType()
	default:
		return nil
	}
}

func merge(l, r map[string][]*descriptor.EnumDescriptorProto) map[string][]*descriptor.EnumDescriptorProto {
	merged := make(map[string][]*descriptor.EnumDescriptorProto)

	for k, v := range l {
		merged[k] = v
	}
	for k, v := range r {
		if cur, ok := merged[k]; ok {
			merged[k] = append(cur, v...)
		} else {
			merged[k] = v
		}
	}

	return merged
}

func appendNestedEnum(formatter enumFormatter, baseDescName string, desc []*descriptor.DescriptorProto) map[string][]*descriptor.EnumDescriptorProto {
	entries := make(map[string][]*descriptor.EnumDescriptorProto)

	for _, d := range desc {
		for _, e := range d.GetEnumType() {
			if v, ok := entries[baseDescName]; ok {
				entries[baseDescName] = append(v, e)
			} else {
				entries[baseDescName] = []*descriptor.EnumDescriptorProto{e}
			}
		}
		entries = merge(entries, appendNestedEnum(formatter, baseDescName, getNested(d)))
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

	entries := make(map[string][]*descriptor.EnumDescriptorProto)
	for _, f := range req.GetProtoFile() {
		descName := strings.ReplaceAll(f.GetPackage(), ".", "_")
		for _, e := range f.GetEnumType() {
			if v, ok := entries[descName]; ok {
				entries[descName] = append(v, e)
			} else {
				entries[descName] = []*descriptor.EnumDescriptorProto{e}
			}
		}
		entries = merge(entries, appendNestedEnum(formatter, descName, getNested(f)))
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
