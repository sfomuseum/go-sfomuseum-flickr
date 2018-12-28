package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-aws/s3"
	"log"
	"path/filepath"
	"strings"
	"sync/atomic"
)

func main() {

	s3_dsn := flag.String("dsn", "", "...")

	flag.Parse()

	cfg, err := s3.NewS3ConfigFromString(*s3_dsn)

	if err != nil {
		log.Fatal(err)
	}

	conn, err := s3.NewS3Connection(cfg)

	if err != nil {
		log.Fatal(err)
	}

	opts := s3.DefaultS3ListOptions()

	count := int64(0)

	cb := func(obj *s3.S3Object) error {

		ext := filepath.Ext(obj.Key)

		if ext == ".json" {
			return nil
		}

		atomic.AddInt64(&count, 1)
		return nil

		key := strings.Replace(obj.Key, ext, "", -1)
		parts := strings.Split(key, "_")

		key_info := fmt.Sprintf("%s_%s_i.json", parts[0], parts[1])
		key_depicts := fmt.Sprintf("%s_d.json", parts[0])

		log.Println(key_info, key_depicts)

		row := make(map[string]string)
		row["uri"] = obj.KeyRaw

		log.Println(row)
		return nil
	}

	err = conn.List(cb, opts)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("COUNT", count)
}
