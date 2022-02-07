appname := hclmergetool
version := v1.1.0
archs := -osarch="linux/amd64" -osarch="darwin/amd64" -osarch="windows/amd64" -osarch="darwin/arm64"
local_arch := $(shell go env GOARCH)
local_os := $(shell go env GOOS)

default: fmt clean build
all: fmt clean build_all_osarch

fmt:
	go mod tidy
	gofmt -w *.go

clean:
	rm -rf build
	rm -rf bin
	rm -rf archives/$(version)
	
build:
	go build -v -o build/$(appname)_$(local_os)_$(local_arch)_$(version)/$(appname)
	mkdir -p bin
	cp build/$(appname)_$(local_os)_$(local_arch)_$(version)/$(appname) bin/

build_all_osarch: 
	gox -verbose $(archs) -output="build/$(appname)_{{.OS}}_{{.Arch}}_$(version)/$(appname)"
	mkdir -p archives/$(version)
	mkdir -p bin
	cd build && for f in $$(ls); do echo $$f; tar czf ../archives/$(version)/$$f.tgz $$f;done
	cp build/$(appname)_$(local_os)_$(local_arch)_$(version)/$(appname) bin/

