CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/sfomuseum/go-sfomuseum-geojson
	cp *.go src/github.com/sfomuseum/go-sfomuseum-geojson/
	cp -r properties src/github.com/sfomuseum/go-sfomuseum-geojson/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	mv vendor/github.com/whosonfirst/go-whosonfirst-geojson-v2/vendor/github.com/whosonfirst/go-whosonfirst-flags vendor/github.com/whosonfirst/

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt *.go
	go fmt properties/*.go
	go fmt properties/*/*.go

bin: 	self
	rm -rf bin/*
	# @GOPATH=$(GOPATH) go build -o bin/sfom-render-ids cmd/sfom-render-ids.go
