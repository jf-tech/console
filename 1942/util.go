package main

import (
	"io/ioutil"
	"math/rand"
	"os"
)

func testProb(prob int) bool {
	return rand.Int()%prob == 0
}

func createTempFile(dir, pattern, content string) *os.File {
	f, _ := ioutil.TempFile(dir, pattern)
	f.Write([]byte(content))
	defer f.Close()
	return f
}
