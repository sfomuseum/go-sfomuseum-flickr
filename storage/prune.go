package storage

import (
	"fmt"
	"github.com/aaronland/go-flickr-archive/photo"
	"github.com/aaronland/go-storage"
	"log"
)

// eventually this will need to be updated to check more than just
// pending photos but it will do for now (20181127/thisisaaronland)

func PruneExisting(store storage.Store, candidates ...photo.Photo) ([]photo.Photo, error) {

	photos := make([]photo.Photo, 0)

	for _, ph := range candidates {

		id := ph.Id()
		k := fmt.Sprintf("%d/%d_r.json", id, id)

		ok, err := store.Exists(k)

		if err != nil {
			return nil, err
		}

		if ok {
			log.Printf("%d is already pending so skipping\n", id)
			continue
		}

		photos = append(photos, ph)
	}

	return photos, nil
}
