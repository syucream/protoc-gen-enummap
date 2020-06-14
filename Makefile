all: build

build:
	go build -o protoc-gen-enummap

check: build
	go test -v ./...
	protoc -I. --plugin=./protoc-gen-enummap --enummap_opt=jsonl --enummap_out=./test/dest test/**/*.proto
	rm test/dest/*

fmt:
	gofmt -w *.go
