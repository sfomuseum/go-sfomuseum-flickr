package queue

import (
	"github.com/aaronland/go-flickr-archive/flickr"
	"github.com/aaronland/go-flickr-archive/photo"
	"github.com/aaronland/go-storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	sfom_flickr "github.com/sfomuseum/go-sfomuseum-flickr"
	sfom_storage "github.com/sfomuseum/go-sfomuseum-flickr/storage"
	"log"
	"net/url"
	"strconv"
	"sync/atomic"
)

// yes, it's called "QueueThingy" ... for now anyway
// (20181128/thisisaaronland)

type QueueThingy struct {
	svc     *sqs.SQS
	queue   string
	store   storage.Store
	depicts []int64
	Debug   bool
	Queued  int64
}

func NewQueueThingy(svc *sqs.SQS, queue string, store storage.Store, depicts ...int64) (*QueueThingy, error) {

	qt := QueueThingy{
		svc:     svc,
		queue:   queue,
		store:   store,
		depicts: depicts,
		Debug:   false,
		Queued:  int64(0),
	}

	return &qt, nil
}

func (qt *QueueThingy) QueuePhotos(candidates ...photo.Photo) error {

	photos, err := sfom_storage.PruneExisting(qt.store, candidates...)

	if err != nil {
		return err
	}

	for _, ph := range photos {

		str_ph, err := sfom_flickr.MarshalPhotoString(ph)

		if err != nil {
			return err
		}

		// please for to add a proper logging thing

		if qt.Debug {
			// log.Println("QUEUE", ph.Id())
			go atomic.AddInt64(&qt.Queued, 1)
			continue
		}

		msg := &sqs.SendMessageInput{
			DelaySeconds: aws.Int64(0),
			MessageBody:  aws.String(str_ph),
			QueueUrl:     aws.String(qt.queue),
		}

		rsp, err := qt.svc.SendMessage(msg)

		if err != nil {
			return err
		}

		log.Println("QUEUE", ph.Id(), *rsp.MessageId)
		go atomic.AddInt64(&qt.Queued, 1)
	}

	return nil
}

func (qt *QueueThingy) QueuePhotosForCLI(str_ids ...string) error {

	candidates := make([]photo.Photo, 0)

	for _, str_id := range str_ids {

		photo_id, err := strconv.ParseInt(str_id, 10, 64)

		if err != nil {
			return err
		}

		ph, err := sfom_flickr.NewSFOMuseumFlickrPhoto(photo_id, qt.depicts...)

		if err != nil {
			return err
		}

		candidates = append(candidates, ph)
	}

	return qt.QueuePhotos(candidates...)
}

func (qt *QueueThingy) QueuePhotosForWOEID(api flickr.API, woe_id int) error {

	query := url.Values{}
	query.Set("woe_id", strconv.Itoa(woe_id))

	return qt.QueuePhotosForSearch(api, query)
}

func (qt *QueueThingy) QueuePhotosForSearch(api flickr.API, query url.Values) error {

	// https://www.flickr.com/services/api/flickr.photos.search

	method := "flickr.photos.search"

	query.Set("license", "1,2,3,4,5,6,7,8,9,10")
	query.Set("safe_search", "1")
	query.Set("media", "photos")

	return qt.QueuePhotosForSPR(api, method, query)
}

func (qt *QueueThingy) QueuePhotosForSPR(api flickr.API, method string, query url.Values) error {

	// query.Set("per_page", "100")

	cb := func(spr flickr.StandardPhotoResponse) error {

		candidates := make([]photo.Photo, 0)

		for _, spr_ph := range spr.Photos.Photos {

			photo_id, err := strconv.ParseInt(spr_ph.ID, 10, 64)

			if err != nil {
				return err
			}

			ph, err := sfom_flickr.NewSFOMuseumFlickrPhoto(photo_id, qt.depicts...)

			if err != nil {
				return err
			}

			candidates = append(candidates, ph)
		}

		return qt.QueuePhotos(candidates...)
	}

	return api.ExecuteMethodPaginated(method, query, cb)
}
