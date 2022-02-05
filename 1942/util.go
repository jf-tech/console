package main

import (
	"io/ioutil"
	"os"
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func createTempFile(dir, pattern, content string) *os.File {
	f, _ := ioutil.TempFile(dir, pattern)
	f.Write([]byte(content))
	defer f.Close()
	return f
}
