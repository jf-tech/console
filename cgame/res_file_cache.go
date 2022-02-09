package cgame

import (
	"io/ioutil"

	"github.com/jf-tech/go-corelib/caches"
)

var resFileCache = caches.NewLoadingCache()

func getResFile(filepath string) ([]byte, error) {
	b, err := resFileCache.Get(filepath, func(interface{}) (interface{}, error) {
		return ioutil.ReadFile(filepath)
	})
	if err != nil {
		return nil, err
	}
	return b.([]byte), nil
}
