package cutil

import (
	"io/ioutil"

	"github.com/jf-tech/go-corelib/caches"
)

var fileCache = caches.NewLoadingCache()

func LoadCachedFile(filepath string) ([]byte, error) {
	b, err := fileCache.Get(filepath, func(interface{}) (interface{}, error) {
		return ioutil.ReadFile(filepath)
	})
	if err != nil {
		return nil, err
	}
	return b.([]byte), nil
}
