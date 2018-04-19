# protoc-gen-enummap

A protoc plugin generates name/number pairs from enum type.

## How to use

```
$ protoc -I. --plugin=path/to/protoc-gen-enummap --enummap_out=./ --enummap_opt=<csv|jsonl> target.proto
```

## Examples

- Sample proto files are here:

```
$ cat test/proto/root.proto
syntax = "proto3";

enum Status {
  UNKNOWN = 0;
  STARTED = 1;
  RUNNING = 2;
}

message Foo {
  enum Status {
    UNKNOWN = 0;
    STARTED = 1;
    RUNNING = 2;
  }
}

$ cat test/proto/sub/child.proto
syntax = "proto3";

message Bar {
  enum Status {
    UNKNOWN = 0;
    STARTED = 1;
    RUNNING = 2;
  }
}
```

- When executing protoc with this plugin

```
$ protoc -I. --plugin=./protoc-gen-enummap --enummap_opt=jsonl --enummap_out=./test/dest test/**/*.proto
```

- Then you can get below files.

```
$ ls test/dest/
test_proto_root_Foo_Status.jsonl  test_proto_root_Status.jsonl  test_proto_sub_child_Bar_Status.jsonl

$ cat test/dest/test_proto_root_Foo_Status.jsonl
{"number": 0, "name": "UNKNOWN"}
{"number": 1, "name": "STARTED"}
{"number": 2, "name": "RUNNING"}
$ cat test/dest/test_proto_root_Status.jsonl
{"number": 0, "name": "UNKNOWN"}
{"number": 1, "name": "STARTED"}
{"number": 2, "name": "RUNNING"}
$ cat test/dest/test_proto_sub_child_Bar_Status.jsonl
{"number": 0, "name": "UNKNOWN"}
{"number": 1, "name": "STARTED"}
{"number": 2, "name": "RUNNING"}
```

