package main

import (
	"context"
	"flag"
	"github.com/aaronland/go-flickr-archive/flickr"
	"github.com/sfomuseum/go-sfomuseum-flickr/aws"
	"github.com/sfomuseum/go-sfomuseum-flickr/queue"
	"github.com/sfomuseum/go-sfomuseum-flickr/storage"
	"github.com/sfomuseum/go-sfomuseum-geojson/properties/sfomuseum"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-index"
	"github.com/whosonfirst/go-whosonfirst-index/utils"
	"github.com/whosonfirst/warning"
	"io"
	"log"
	"strconv"
	"sync/atomic"
)

func main() {

	sqs_dsn := flag.String("sqs-dsn", "", "...")
	storage_dsn := flag.String("storage-dsn", "", "...")

	mode := flag.String("mode", "repo", "...")

	api_key := flag.String("api-key", "", "...")
	api_secret := flag.String("api-secret", "", "...")

	debug := flag.Bool("debug", false, "...")

	flag.Parse()

	svc, sqs_queue, err := aws.NewSQSServiceFromString(*sqs_dsn)

	if err != nil {
		log.Fatal(err)
	}

	api, err := flickr.NewFlickrAuthAPI(*api_key, *api_secret)

	if err != nil {
		log.Fatal(err)
	}

	total := int64(0)
	airports := int64(0)

	cb := func(fh io.Reader, ctx context.Context, args ...interface{}) error {

		_, err := index.PathForContext(ctx)

		if err != nil {
			return err
		}

		ok, err := utils.IsPrincipalWOFRecord(fh, ctx)

		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		f, err := feature.LoadFeatureFromReader(fh)

		if err != nil && !warning.IsWarning(err) {
			return err
		}

		if sfomuseum.Placetype(f) != "airport" {
			return nil
		}

		concordances, err := whosonfirst.Concordances(f)

		if err != nil {
			return err
		}

		str_id, ok := concordances["gp:id"]

		if !ok {
			return nil
		}

		iata_code, ok := concordances["iata:code"]

		if !ok {
			return nil
		}

		woe_id, err := strconv.Atoi(str_id)

		if err != nil {
			return err
		}

		store, err := storage.NewStoreWithSuffix(*storage_dsn, f.Id())

		if err != nil {
			return err
		}

		depicts := whosonfirst.Id(f)

		qt, err := queue.NewQueueThingy(svc, sqs_queue, store, depicts)

		if err != nil {
			return err
		}

		qt.Debug = *debug

		err = qt.QueuePhotosForWOEID(api, woe_id)

		count := qt.Queued

		atomic.AddInt64(&airports, 1)
		atomic.AddInt64(&total, count)

		log.Println("QUEUED", iata_code, count, atomic.LoadInt64(&total), atomic.LoadInt64(&airports))

		if err != nil {
			log.Println("WHAT", iata_code, err)
		}

		return nil
	}

	i, err := index.NewIndexer(*mode, cb)

	if err != nil {
		log.Fatal(err)
	}

	paths := flag.Args()

	err = i.IndexPaths(paths)

	log.Println("QUEUED", total, airports)

	if err != nil {
		log.Fatal(err)
	}

}
