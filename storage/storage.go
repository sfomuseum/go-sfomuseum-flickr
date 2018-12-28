package storage

import (
	"errors"
	"github.com/aaronland/go-storage"
	"github.com/aaronland/go-storage-s3"
	"github.com/aaronland/go-string/dsn"
	_ "log"
	"path/filepath"
)

func NewStoreWithSuffix(str_dsn string, suffix string) (storage.Store, error) {

	dsn_map, err := dsn.StringToDSNWithKeys(str_dsn, "storage")

	if err != nil {
		return nil, err
	}

	switch dsn_map["storage"] {

	case "fs":
		dsn_map["root"] = filepath.Join(dsn_map["root"], suffix)
	case "s3":
		dsn_map["prefix"] = filepath.Join(dsn_map["prefix"], suffix)
	default:
		return nil, errors.New("Invalid storage layer")

	}

	return NewStore(dsn_map.String())
}

func NewStore(str_dsn string) (storage.Store, error) {

	dsn_map, err := dsn.StringToDSNWithKeys(str_dsn, "storage")

	if err != nil {
		return nil, err
	}

	var store storage.Store
	var store_err error

	switch dsn_map["storage"] {

	case "fs":
		store, store_err = storage.NewFSStore(str_dsn)
	case "s3":
		store, store_err = s3.NewS3Store(str_dsn)
	default:
		store_err = errors.New("Invalid storage layer")

	}

	return store, store_err
}
