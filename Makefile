
.PHONY: build

build:
	go build github.com/yuce/importtst/cmd/generate-import
	go build github.com/yuce/importtst/cmd/pilosa-import