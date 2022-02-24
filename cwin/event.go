package cwin

import (
	"time"

	"github.com/jf-tech/console/cterm"
)

type EventResponse int

const (
	EventHandled = EventResponse(iota)
	EventNotHandled
	EventLoopStop
)

type EventHandler func(ev cterm.Event) EventResponse

type EventLoopSleepDurationFunc func() time.Duration

func RunEventLoop(s *Sys, evHandler EventHandler, f ...EventLoopSleepDurationFunc) {
	if evHandler == nil {
		panic("evHandler cannot be nil")
	}
	durationF := defaultEventLoopSleepDurationFunc
	if len(f) > 0 {
		durationF = f[0]
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
		time.Sleep(durationF())
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

const (
	defaultEventLoopSleepDuration = time.Millisecond
)

var (
	defaultEventLoopSleepTimestamp = time.Time{}
)

func defaultEventLoopSleepDurationFunc() time.Duration {
	now := time.Now()
	duration := defaultEventLoopSleepDuration - now.Sub(defaultEventLoopSleepTimestamp)
	defaultEventLoopSleepTimestamp = now
	return duration
}
