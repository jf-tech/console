package cwin

import (
	"github.com/jf-tech/console/cterm"
)

func Keys(keys ...interface{}) []cterm.Event {
	var ks []cterm.Event
	for _, k := range keys {
		if ch, ok := k.(rune); ok {
			if ch == 0 {
				panic("rune cannot be zero")
			}
			ks = append(ks, cterm.Event{Type: cterm.EventKey, Ch: ch})
			continue
		}
		if key, ok := k.(cterm.Key); ok {
			ks = append(ks, cterm.Event{Type: cterm.EventKey, Key: key})
			continue
		}
	}
	return ks
}

func FindKey(keys []cterm.Event, key cterm.Event) bool {
	for _, ev := range keys {
		if ev.Type != cterm.EventKey {
			continue
		}
		if key.Ch != 0 && ev.Ch == key.Ch {
			return true
		}
		if key.Ch == 0 && ev.Ch == 0 && ev.Key == key.Key {
			return true
		}
	}
	return false
}
