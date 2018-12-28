CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test ! -d src/github.com/aaronland/go-storage-s3; then mkdir -p src/github.com/aaronland/go-storage-s3; fi
	cp -r *.go src/github.com/aaronland/go-storage-s3/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/aaronland/go-storage"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-aws"
	mv src/github.com/whosonfirst/go-whosonfirst-aws/vendor/github.com/aws src/github.com/

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt *.go

bin: self
	@GOPATH=$(GOPATH) go build -o bin/exists cmd/exists.go
	@GOPATH=$(GOPATH) go build -o bin/create cmd/create.go
	@GOPATH=$(GOPATH) go build -o bin/read cmd/read.go
