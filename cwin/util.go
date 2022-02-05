package cwin

import (
	"fmt"
	"strings"
	"sync/atomic"
	"unicode/utf8"

	"github.com/jf-tech/go-corelib/maths"
	"github.com/nsf/termbox-go"
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

// if f == nil, SyncExpectKey waits for any single key and then returns
// if f != nil, SyncExpectKey repeatedly waits for a key & has it processed by f, if f returns false
func SyncExpectKey(f func(termbox.Key, rune) bool) {
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey && (f == nil || f(ev.Key, ev.Ch)) {
			break
		}
	}
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
