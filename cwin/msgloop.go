package cwin

import "github.com/jf-tech/console/cterm"

type MsgLoopResponseType int

const (
	MsgLoopContinue = MsgLoopResponseType(iota)
	MsgLoopStop
	MsgLoopContinueWithFullRefresh // only use when absolutely necessary
)

type MsgLoopFunc func(ev cterm.Event) MsgLoopResponseType

func TrueForMsgLoopStop(b bool) MsgLoopResponseType {
	if b {
		return MsgLoopStop
	}
	return MsgLoopContinue
}

func FalseForMsgLoopStop(b bool) MsgLoopResponseType {
	return TrueForMsgLoopStop(!b)
}
