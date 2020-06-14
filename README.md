# protoc-gen-enummap

A protoc plugin generates name/number pairs in CSV/JSONL/SQL/etc... from enum type.

## How to use

```
$ protoc -I. --plugin=path/to/protoc-gen-enummap --enummap_out=./ --enummap_opt=<csv|jsonl|sql> target.proto
```

## Examples

- Sample proto files are here:

```
$ cat test/proto/root.proto
syntax = "proto3";

package root;

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

package root.sub;

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
$ cat test/dest/root_Status.json
{"number": 0, "name": "UNKNOWN"}
{"number": 1, "name": "STARTED"}
{"number": 2, "name": "RUNNING"}
$ cat test/dest/root_Foo_Status.json
{"number": 0, "name": "UNKNOWN"}
{"number": 1, "name": "STARTED"}
{"number": 2, "name": "RUNNING"}
$ cat test/dest/root_sub_Bar_Status.json
{"number": 0, "name": "UNKNOWN"}
{"number": 1, "name": "STARTED"}
{"number": 2, "name": "RUNNING"}
```

## Integration

### BigQuery

- Convert proto to csv

```
$ protoc -I. --plugin=path/to/protoc-gen-enummap --enummap_out=./ --enummap_opt=csv target.proto
```

- Load the csv

```
$ bq load --replace --source_format=CSV <project_id>:<dataset_id>.<table_name> <path_to_csv> "number:integer,name:string"
```
