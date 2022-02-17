package cutil

import (
	"path"
	"runtime"
)

func GetCurFileDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
