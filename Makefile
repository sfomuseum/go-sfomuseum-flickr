CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/sfomuseum/go-sfomuseum-flickr
	cp *.go src/github.com/sfomuseum/go-sfomuseum-flickr/
	cp -r queue src/github.com/sfomuseum/go-sfomuseum-flickr/
	cp -r storage src/github.com/sfomuseum/go-sfomuseum-flickr/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-cli"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-index"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-aws"
	@GOPATH=$(GOPATH) go get -u "github.com/sfomuseum/go-sfomuseum-geojson"
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/gjson"
	@GOPATH=$(GOPATH) go get -u "github.com/aaronland/go-storage-s3"
	@GOPATH=$(GOPATH) go get -u "github.com/aaronland/go-string"
	@GOPATH=$(GOPATH) go get -u "github.com/aaronland/go-flickr-archive/flickr"
	@GOPATH=$(GOPATH) go get -u "github.com/aws/aws-lambda-go/..."
	mv src/github.com/aaronland/go-storage-s3/vendor/github.com/aaronland/go-storage src/github.com/aaronland/
	mv src/github.com/sfomuseum/go-sfomuseum-geojson/vendor/github.com/whosonfirst/go-whosonfirst-geojson-v2 src/github.com/whosonfirst/
	mv src/github.com/whosonfirst/go-whosonfirst-geojson-v2/vendor/github.com/whosonfirst/go-whosonfirst-flags src/github.com/whosonfirst/
	mv src/github.com/whosonfirst/go-whosonfirst-geojson-v2/vendor/github.com/whosonfirst/warning src/github.com/whosonfirst/
	mv src/github.com/whosonfirst/go-whosonfirst-aws/vendor/github.com/aws/aws-sdk-go src/github.com/aws/
	rm -rf src/github.com/aaronland/go-storage-s3/vendor/github.com/aws/aws-sdk-go
	rm -rf src/github.com/aaronland/go-storage-s3/vendor/github.com/whosonfirst/go-whosonfirst-aws
	rm -rf src/github.com/whosonfirst/go-whosonfirst-aws/vendor/github.com/aws/aws-sdk-go
	rm -rf src/github.com/aaronland/go-flickr-archive/vendor/github.com/aaronland/go-storage

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt *.go
	go fmt cmd/*.go
	go fmt queue/*.go
	go fmt storage/*.go

bin: 	self
	rm -rf bin/*
	@GOPATH=$(GOPATH) go build -o bin/flickr-queue-airports cmd/flickr-queue-airports.go
	@GOPATH=$(GOPATH) go build -o bin/flickr-queue-photos cmd/flickr-queue-photos.go
	@GOPATH=$(GOPATH) go build -o bin/flickr-archive-photos cmd/flickr-archive-photos.go
	@GOPATH=$(GOPATH) go build -o bin/flickr-process-photos cmd/flickr-process-photos.go
	@GOPATH=$(GOPATH) go build -o bin/flickr-list-pending cmd/flickr-list-pending.go

lambda: lambda-archive lambda-process

lambda-archive:	
	@make self
	if test -f main; then rm -f main; fi
	if test -f archive.zip; then rm -f archive.zip; fi
	@GOPATH=$(GOPATH) GOOS=linux go build -o main cmd/flickr-archive-photos.go
	zip archive.zip main
	rm -f main

lambda-process:	
	@make self
	if test -f main; then rm -f main; fi
	if test -f process.zip; then rm -f process.zip; fi
	@GOPATH=$(GOPATH) GOOS=linux go build -o main cmd/flickr-process-photos.go
	zip process.zip main
	rm -f main
