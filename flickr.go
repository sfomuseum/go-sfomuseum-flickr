package flickr

import (
	"encoding/json"
	"github.com/aaronland/go-flickr-archive/photo"
)

type SFOMuseumFlickrPhoto struct {
	photo.Photo `json:",omitempty"`
	ID          int64   `json:"id"`
	Depicts     []int64 `json:"depicts"`
}

func NewSFOMuseumFlickrPhoto(photo_id int64, depicts ...int64) (photo.Photo, error) {

	ph := SFOMuseumFlickrPhoto{
		ID:      photo_id,
		Depicts: depicts,
	}

	return &ph, nil
}

func NewSFOMuseumFlickrPhotoFromString(raw string) (photo.Photo, error) {

	var ph SFOMuseumFlickrPhoto
	err := json.Unmarshal([]byte(raw), &ph)

	if err != nil {
		return nil, err
	}

	return &ph, nil
}

func (ph *SFOMuseumFlickrPhoto) Id() int64 {
	return ph.ID
}

func MarshalPhoto(ph photo.Photo) ([]byte, error) {
	return json.Marshal(ph)
}

func MarshalPhotoString(ph photo.Photo) (string, error) {

	enc_ph, err := MarshalPhoto(ph)

	if err != nil {
		return "", err
	}

	return string(enc_ph), nil
}
