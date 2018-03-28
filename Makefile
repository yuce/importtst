
.PHONY: build

DMD ?= dmd

build:
	go build github.com/yuce/importtst/cmd/pilosa-import
	$(DMD) -O -release generate_import.d