# go-sfomuseum-flickr

Go package for working with Flickr photos in an SFO Museum context, principally archiving photos.

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.7 so let's just assume you need [Go 1.11](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Important

This should still be considered experimental. It is not well documented yet nor is it a general-purpose "archive Flickr" package; it is tailored to the specific needs of SFO Museum.

It is provided as-is in the spirit of sharing and maybe serving as a point of reference for others. 

For a generic (and equally experiemental) Flickr archiving tool there is the [go-flickr-archive](https://github.com/aaronland/go-flickr-archive) package which this package uses.

## Tools

### flickr-queue-photos

Put one or more photo IDs in a (SQS) queue for processing. Eventually this code will be wrapped by something that collects all the photos for a given airport (WOE ID).

```
$> ./bin/flickr-queue-photos -storage-dsn 'storage=s3 bucket={S3_BUCKET} prefix={S3_PREFIX} region={S3_REGION} credentials={S3_CREDENTIALS}' -sqs-dsn 'queue={SQS_QUEUE} credentials={SQS_CREDENTIALS} region={SQS_REGION}' -depicts 102535681 37340024961

Success c1e60e95-bc01-4fae-b11b-b976f4844da2
```

Queue-ing photos from a `flickr.photos.search` API query:

```
$> ./bin/flickr-queue-photos -storage-dsn 'storage=s3 bucket={S3_BUCKET} prefix={S3_PREFIX} region={S3_REGION} credentials={S3_CREDENTIALS}' -sqs-dsn 'region={SQS_REGION} credentials={SQS_CREDENTIALS} queue={SQS_QUEUE}' -search -api-key {FLICKR_APIKEY} -param woe_id=22411879 -param per_page=100 -depicts 102545733
```

### flickr-archive-photos

Archive one or more photo IDs where "archive" means fetch the largest available photo and the output of the `flickr.photos.getInfo` API method and put them in a place for later processing. Normally this is expected to be run as a Lambda function connected to an SQS queue.

```
$> ./bin/flickr-archive-photos -api-key {FLICKR_APIKEY} -storage-dsn 'storage=s3 bucket={S3_BUCKET} prefix={S3_PREFIX} region={S3_REGION} credentials={S3_CREDENTIALS}' 17834003703
```

Or the same thing but archiving things locally to the current directory (`.`):

```
$> ./bin/flickr-archive-photos -api-key {FLICKR_APIKEY} -depicts 1234 -storage-dsn 'storage=fs root=.' 34607346293

# time passes...

$> ls -la 34607346293/
total 3528
-rw-r--r--  1 asc  staff     3202 Dec 28 11:03 34607346293_2e50d8058d_i.json
-rw-r--r--  1 asc  staff  1796141 Dec 28 11:03 34607346293_a4da86aaff_k.jpg
-rw-r--r--  1 asc  staff       35 Dec 28 11:03 34607346293_r.json

$> cat 34607346293/34607346293_r.json
{"id":34607346293,"depicts":[1234]}
```

Currently the (`go-flickr-archive`) code will archive the response value of a call to the `flickr.photos.getInfo` API method as well as the largest photo is that is publicly available.

The `{PHOTO_ID}_r.json` file contains metadata specific to that particular archiving request. It is left up to consumers to determine its use and function.

## Lambda

Yes. If you run the handy `make lambda-archive` Makefile target then a `archive.zip` binary will be created (derived from the `flickr-archive-photos.go` tool) that you can upload and run as a Lambda function.

As of this writing the Lambda function is hard-wired to receive SQS messages containing JSON-encoded `SFOMuseumFlickrPhoto` documents:

```
type SFOMuseumFlickrPhoto struct {
	photo.Photo `json:",omitempty"`
	ID          int64   `json:"id"`
	Depicts     []int64 `json:"depicts"`
}
```

Which should look familiar since it's basically the same as the `34607346293/34607346293_r.json` file above.

When configuring your Lambda function be sure the add a `SFOMUSEUM_FLICKR_LAMBDA` environment variable whose value is `true`.

It is expected that some or all of the Lambda-related stuff (code, environment variables, etc.) _will_ change in time.

## To do still

1. When queue-ing photos pass along some amount of metadata to denote what (in SFO Museum -land) the photo depicts (DONE)
2. Check to see whether photo has already been stored (DONE)
3. When archiving photos store the SFO Museum depiction information (set in `flickr-queue-photos`) (DONE)
4. Search query/criteria to queue
5. Final image processing and association (depictions)

The 1-3 should be prioritized so that we can get ahead of any photo purging that may or may not happen to non-pro Flickr photos come January 2019. Good times...

## See also

* https://github.com/aaronland/go-flickr-archive
* https://github.com/aaronland/go-storage
* https://github.com/aaronland/go-storage-s3
* https://github.com/whosonfirst/go-whosonfirst-aws