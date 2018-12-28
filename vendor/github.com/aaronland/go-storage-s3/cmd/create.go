package main

import (
	"github.com/aaronland/go-storage-s3"
	"flag"
	"log"
)

func main() {

	storage_dsn := flag.String("dsn", "", "...")
	path := flag.String("path", "", "...")
	
	flag.Parse()

	store, err := s3.NewS3Store(*storage_dsn)

	if err != nil {
		log.Fatal(err)
	}

	fh, err := store.Create(*path)

	if err != nil {
		log.Fatal(err)
	}

	for _, word := range flag.Args() {

		_, err := fh.Write([]byte(word))

		if err != nil {
			log.Fatal(err)
		}
	}

	fh.Close()
}
