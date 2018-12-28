package main

import (
	"flag"
	"github.com/aaronland/go-flickr-archive/flickr"
	"github.com/sfomuseum/go-sfomuseum-flickr/aws"
	"github.com/sfomuseum/go-sfomuseum-flickr/queue"
	"github.com/sfomuseum/go-sfomuseum-flickr/storage"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"log"
	"net/url"
)

func main() {

	sqs_dsn := flag.String("sqs-dsn", "", "...")
	storage_dsn := flag.String("storage-dsn", "", "...")

	api_key := flag.String("api-key", "", "...")
	api_secret := flag.String("api-secret", "", "...")

	var depicts flags.MultiInt64
	flag.Var(&depicts, "depicts", "")

	search := flag.Bool("search", false, "...")
	debug := flag.Bool("debug", false, "...")

	var params flags.KeyValueArgs
	flag.Var(&params, "param", "")

	flag.Parse()

	svc, sqs_queue, err := aws.NewSQSServiceFromString(*sqs_dsn)

	if err != nil {
		log.Fatal(err)
	}

	store, err := storage.NewStore(*storage_dsn)

	if err != nil {
		log.Fatal(err)
	}

	// yes, it's called "QueueThingy" ... for now anyway
	// (20181128/thisisaaronland)

	qt, err := queue.NewQueueThingy(svc, sqs_queue, store, depicts...)

	if err != nil {
		log.Fatal(err)
	}

	qt.Debug = *debug

	// it's tempting to add a *lambda interface (where you'd pass it
	// a WOE ID and have it invoke the search stuff) but that feels
	// like yak-shaving right now (20181128/thisisaaronland)

	if *search {

		api, err := flickr.NewFlickrAuthAPI(*api_key, *api_secret)

		if err != nil {
			log.Fatal(err)
		}

		query := url.Values{}

		for _, p := range params {
			query.Set(p.Key, p.Value)
		}

		err = qt.QueuePhotosForSearch(api, query)

		if err != nil {
			log.Fatal(err)
		}

	} else {

		ids := flag.Args()

		err := qt.QueuePhotosForCLI(ids...)

		if err != nil {
			log.Fatal(err)
		}
	}
}
