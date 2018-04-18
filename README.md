# protoc-gen-enummap

A protoc plugin generates name/number pairs from enum type.

## How to use

```
$ protoc -I. --plugin=path/to/protoc-gen-enummap --enummap_out=./ --enummap_opt=<csv|jsonl> target.proto
```
