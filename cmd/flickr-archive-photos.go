package main

import (
	"context"
	"errors"
	"flag"
	"github.com/aaronland/go-flickr-archive"
	"github.com/aaronland/go-flickr-archive/archivist"
	"github.com/aaronland/go-flickr-archive/flickr"
	"github.com/aaronland/go-flickr-archive/photo"
	"github.com/aaronland/go-storage"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	aws_sqs "github.com/aws/aws-sdk-go/service/sqs"
	sfom_flickr "github.com/sfomuseum/go-sfomuseum-flickr"
	sfom_storage "github.com/sfomuseum/go-sfomuseum-flickr/storage"
	"github.com/whosonfirst/go-whosonfirst-aws/sqs"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"log"
	"os"
	"strconv"
)

func archive_photos(arch archive.Archivist, api flickr.API, store storage.Store, candidates ...photo.Photo) error {

	photos, err := sfom_storage.PruneExisting(store, candidates...)

	if err != nil {
		return err
	}

	if len(photos) == 0 {
		return nil
	}

	err = arch.ArchivePhotos(api, photos...)

	if err != nil {
		return err
	}

	return nil
}

func main() {

	do_lambda := flag.Bool("lambda", false, "...")
	do_sqs := flag.Bool("sqs", false, "...")

	api_key := flag.String("api-key", "", "...")
	api_secret := flag.String("api-secret", "", "...")

	storage_dsn := flag.String("storage-dsn", "", "...")
	sqs_dsn := flag.String("sqs-dsn", "", "...")

	var depicts flags.MultiInt64
	flag.Var(&depicts, "depicts", "")

	flag.Parse()

	err := flags.SetFlagsFromEnvVars("SFOMUSEUM_FLICKR")

	if err != nil {
		log.Fatal(err)
	}

	api, err := flickr.NewFlickrAuthAPI(*api_key, *api_secret)

	if err != nil {
		log.Fatal(err)
	}

	store, err := sfom_storage.NewStore(*storage_dsn)

	if err != nil {
		log.Fatal(err)
	}

	opts, err := archivist.DefaultStaticArchivistOptions()

	if err != nil {
		log.Fatal(err)
	}

	opts.ArchiveRequest = true

	arch, err := archivist.NewStaticArchivist(store, opts)

	if err != nil {
		log.Fatal(err)
	}

	if *do_lambda {

		// some day we might have a different kind of lambda trigger
		// today we do not (20181127/thisisaaronland)

		handler := func(ctx context.Context, sqsEvent events.SQSEvent) error {

			count := len(sqsEvent.Records)

			if count == 0 {
				return errors.New("zero-length records")
			}

			photos := make([]photo.Photo, count)

			for i, msg := range sqsEvent.Records {

				ph, err := sfom_flickr.NewSFOMuseumFlickrPhotoFromString(msg.Body)

				if err != nil {
					return err
				}

				photos[i] = ph
			}

			err := archive_photos(arch, api, store, photos...)

			if err != nil {
				return err
			}

			return nil
		}

		lambda.Start(handler)

	} else if *do_sqs {

		// READ from SQS - not SEND to SQS (20181228/thisisaaronland)

		// I don't really understand how the long-polling stuff is supposed
		// to work and this is mostly here as an exercise and for debugging
		// (20181127/thisisaaronland)

		svc, sqs_queue, err := sqs.NewSQSServiceWithDSN(*sqs_dsn)

		if err != nil {
			log.Fatal(err)
		}

		params := &aws_sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(sqs_queue),
			MaxNumberOfMessages: aws.Int64(1),
			VisibilityTimeout:   aws.Int64(30),
			WaitTimeSeconds:     aws.Int64(20),
		}

		rsp, err := svc.ReceiveMessage(params)

		if err != nil {
			log.Fatal(err)
		}

		for _, msg := range rsp.Messages {

			ph, err := sfom_flickr.NewSFOMuseumFlickrPhotoFromString(*msg.Body)

			if err != nil {
				log.Println(err)
				continue
			}

			err = archive_photos(arch, api, store, ph)

			if err != nil {
				continue
			}

			delete_params := &aws_sqs.DeleteMessageInput{
				QueueUrl:      aws.String(sqs_queue),
				ReceiptHandle: msg.ReceiptHandle,
			}

			_, err = svc.DeleteMessage(delete_params)

			if err != nil {
				log.Println(err)
			}

		}

	} else {

		count := len(flag.Args())
		photos := make([]photo.Photo, count)

		for i, str_id := range flag.Args() {

			photo_id, err := strconv.ParseInt(str_id, 10, 64)

			if err != nil {
				log.Fatal(err)
			}

			ph, err := sfom_flickr.NewSFOMuseumFlickrPhoto(photo_id, depicts...)

			if err != nil {
				log.Fatal(err)
			}

			photos[i] = ph
		}

		err := archive_photos(arch, api, store, photos...)

		if err != nil {
			log.Fatal(err)
		}

	}

	os.Exit(0)
}
