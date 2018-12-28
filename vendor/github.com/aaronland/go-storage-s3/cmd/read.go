package main

import (
	"github.com/aaronland/go-storage-s3"
	"flag"
	"io"
	"log"
	"os"
)

func main() {

	storage_dsn := flag.String("dsn", "", "...")
	path := flag.String("path", "", "...")
	
	flag.Parse()

	store, err := s3.NewS3Store(*storage_dsn)

	if err != nil {
		log.Fatal(err)
	}

	fh, err := store.Get(*path)

	if err != nil {
		log.Fatal(err)
	}

	defer fh.Close()
	
	_, err = io.Copy(os.Stdout, fh)

	if err != nil {
		log.Fatal(err)
	}
	
}
