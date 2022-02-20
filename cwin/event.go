package cwin

import (
	"github.com/jf-tech/console/cterm"
)

type EventResponse int

const (
	EventHandled = EventResponse(iota)
	EventNotHandled
	EventLoopStop
)

type EventHandler func(ev cterm.Event) EventResponse

func RunEventLoop(s *Sys, evHandler EventHandler) {
	if evHandler == nil {
		panic("evHandler cannot be nil")
	}
loop:
	for {
		resp := EventNotHandled
		ev := s.TryGetEvent()
		resp = evHandler(ev)
		s.Update()
		if resp == EventLoopStop {
			break loop
		}
	}
}

func NopHandledEventHandler(ev cterm.Event) EventResponse {
	return EventHandled
}

func NopNotHandledEventHandler(ev cterm.Event) EventResponse {
	return EventNotHandled
}

func TrueForEventSystemStop(b bool) EventResponse {
	if b {
		return EventLoopStop
	}
	return EventHandled
}

func FalseForEventSystemStop(b bool) EventResponse {
	return TrueForEventSystemStop(!b)
}
