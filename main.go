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
	"sql":   &sqlEnumFormatter{},
}

func getNested(d interface{}) []*descriptor.DescriptorProto {
	if fdp, ok := d.(*descriptor.FileDescriptorProto); ok {
		return fdp.GetMessageType()
	} else if dp, ok := d.(*descriptor.DescriptorProto); ok {
		return dp.GetNestedType()
	} else {
		return nil
	}
}

func getDescName(d interface{}) string {
	if fdp, ok := d.(*descriptor.FileDescriptorProto); ok {
		return strings.Replace(fdp.GetPackage(), ".", "_", -1) + "__"
	} else if dp, ok := d.(*descriptor.DescriptorProto); ok {
		return dp.GetName() + "_"
	} else {
		return ""
	}
}

func appendNestedEnum(file []*plugin.CodeGeneratorResponse_File, formatter enumFormatter, prefix string, desc []*descriptor.DescriptorProto) []*plugin.CodeGeneratorResponse_File {
	for _, d := range desc {
		descName := prefix + getDescName(d)
		for _, e := range d.GetEnumType() {
			var contents []string
			for _, ev := range e.GetValue() {
				contents = append(contents, formatter.printLine(ev))
			}
			file = append(file, &plugin.CodeGeneratorResponse_File{
				Name:    proto.String(descName + e.GetName() + formatter.extension()),
				Content: proto.String(formatter.printHeader(descName+e.GetName()) + strings.Join(contents, "")),
			})
		}
		file = appendNestedEnum(file, formatter, descName, getNested(d))
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
		descName := getDescName(f)
		for _, e := range f.GetEnumType() {
			var contents []string
			for _, ev := range e.GetValue() {
				contents = append(contents, formatter.printLine(ev))
			}
			resp.File = append(resp.File, &plugin.CodeGeneratorResponse_File{
				Name:    proto.String(descName + e.GetName() + formatter.extension()),
				Content: proto.String(formatter.printHeader(descName+e.GetName()) + strings.Join(contents, "")),
			})
		}
		resp.File = appendNestedEnum(resp.File, formatter, descName, getNested(f))
	}

	buf, err = proto.Marshal(&resp)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stdout.Write(buf); err != nil {
		log.Fatal(err)
	}
}
