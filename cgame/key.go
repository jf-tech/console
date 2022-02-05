package cgame

import "github.com/nsf/termbox-go"

func Keys(keys ...interface{}) []termbox.Event {
	var ks []termbox.Event
	for _, k := range keys {
		if ch, ok := k.(rune); ok {
			if ch == 0 {
				panic("rune cannot be zero")
			}
			ks = append(ks, termbox.Event{Type: termbox.EventKey, Ch: ch})
			continue
		}
		if key, ok := k.(termbox.Key); ok {
			ks = append(ks, termbox.Event{Type: termbox.EventKey, Key: key})
			continue
		}
	}
	return ks
}

func FindKey(keys []termbox.Event, key termbox.Event) bool {
	for _, ev := range keys {
		if ev.Type != termbox.EventKey {
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
