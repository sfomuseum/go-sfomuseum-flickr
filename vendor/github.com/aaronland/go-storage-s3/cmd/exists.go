package main

import (
	"github.com/aaronland/go-storage-s3"
	"flag"
	"log"
)

func main() {

	storage_dsn := flag.String("dsn", "", "...")	
	flag.Parse()

	store, err := s3.NewS3Store(*storage_dsn)

	if err != nil {
		log.Fatal(err)
	}

	for _, path := range flag.Args() {

		ok, err := store.Exists(path)
		log.Println(path, ok, err)
	}
}
