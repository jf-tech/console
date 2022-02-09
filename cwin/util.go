package cwin

import (
	"fmt"
	"strings"
	"sync/atomic"
	"unicode/utf8"

	"github.com/jf-tech/go-corelib/maths"
)

func TextDimension(s string) Rect {
	if len(s) <= 0 {
		return Rect{0, 0, 0, 0}
	}
	lines := strings.Split(s, "\n")
	maxLineLen := 0
	for _, line := range lines {
		maxLineLen = maths.MaxInt(maxLineLen, utf8.RuneCountInString(line))
	}
	return Rect{0, 0, maxLineLen, len(lines)}
}

var (
	globalUIDCounter int64 = 0
)

func GenUID() int64 {
	return atomic.AddInt64(&globalUIDCounter, 1)
}

var (
	byteSizeStrs = []string{"B", "KB", "MB", "GB", "TB", "EB", "ZB"}
)

func ByteSizeStr(s int64) string {
	ss, p := s, 0
	for ; ss >= 1024; ss /= 1024 {
		p++
	}
	// we're safe as the max int64 value is about 9 ZB.
	return fmt.Sprintf("%d %s", ss, byteSizeStrs[p])
}
