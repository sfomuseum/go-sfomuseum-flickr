package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"log"
)

func main() {

	do_lambda := flag.Bool("lambda", false, "...")

	flag.Parse()

	err := flags.SetFlagsFromEnvVars("SFOMUSEUM_FLICKR")

	if err != nil {
		log.Fatal(err)
	}

	if *do_lambda {

		handler := func(ctx context.Context, s3Event events.S3Event) {

			for _, record := range s3Event.Records {
				s3 := record.S3
				fmt.Printf("[%s - %s] Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, s3.Bucket.Name, s3.Object.Key)
			}
		}

		lambda.Start(handler)

	} else {

		log.Println("PLEASE WRITE ME")
	}

}
