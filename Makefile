build:
	go build -o protoc-gen-enummap

fmt:
	gofmt -w *.go
