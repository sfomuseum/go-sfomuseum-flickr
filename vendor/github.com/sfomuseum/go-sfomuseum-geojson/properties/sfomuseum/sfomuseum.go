package sfomuseum

import (
	"github.com/whosonfirst/go-whosonfirst-flags"
	"github.com/whosonfirst/go-whosonfirst-flags/existential"	
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
	_ "log"
)

func Placetype(f geojson.Feature) string {

	possible := []string{
		"properties.sfomuseum:placetype",
	}
	
	return utils.StringProperty(f.Bytes(), possible, f.Placetype())
}

func IsSFO(f geojson.Feature) (flags.ExistentialFlag, error) {

	possible := []string{
		"properties.sfomuseum:is_sfo",
	}

	v := utils.Int64Property(f.Bytes(), possible, -1)

	if v == 1 || v == 0 {
		return existential.NewKnownUnknownFlag(v)
	}

	return existential.NewKnownUnknownFlag(-1)	
}
