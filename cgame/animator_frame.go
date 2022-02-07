package cgame

import (
	"time"
)

type AnimatorFrameCfg struct {
	Frames                  FrameProvider
	KeepAliveWhenFinished   bool
	AfterFrame, AfterFinish func(Sprite)
}

type AnimatorFrame struct {
	cfg AnimatorFrameCfg

	clock *Clock

	curFrameDuration    time.Duration
	curFrameStartedTime time.Duration
}

func (af *AnimatorFrame) Animate(s Sprite) AnimatorState {
	af.checkToInit(s)
	elapsed := af.clock.Now() - af.curFrameStartedTime
	if elapsed < af.curFrameDuration {
		return AnimatorRunning
	}
	if af.setNextFrame(s) {
		return AnimatorRunning
	}
	if !af.cfg.KeepAliveWhenFinished {
		s.Mgr().AddEvent(NewSpriteEventDelete(s))
	}
	if af.cfg.AfterFinish != nil {
		af.cfg.AfterFinish(s)
	}
	return AnimatorCompleted
}

func (af *AnimatorFrame) setNextFrame(s Sprite) (more bool) {
	var f Frame
	if f, af.curFrameDuration, more = af.cfg.Frames.Next(); !more {
		return false
	}
	FrameToWin(f, s.Win())
	af.curFrameStartedTime = af.clock.Now()
	if af.cfg.AfterFrame != nil {
		af.cfg.AfterFrame(s)
	}
	return true
}

func (af *AnimatorFrame) checkToInit(s Sprite) {
	if af.clock != nil {
		return
	}
	af.clock = s.Game().MasterClock
	if !af.setNextFrame(s) {
		panic("Frames cannot be empty")
	}
}

func NewAnimatorFrame(c AnimatorFrameCfg) *AnimatorFrame {
	return &AnimatorFrame{cfg: c}
}
