package cwin

import "github.com/jf-tech/console/cterm"

type EventLoopResponseType int

const (
	EventLoopContinue = EventLoopResponseType(iota)
	EventLoopStop
)

type EventLoopFunc func(ev cterm.Event) EventLoopResponseType

func TrueForEventLoopStop(b bool) EventLoopResponseType {
	if b {
		return EventLoopStop
	}
	return EventLoopContinue
}

func FalseForEventLoopStop(b bool) EventLoopResponseType {
	return TrueForEventLoopStop(!b)
}
