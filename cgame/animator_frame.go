package cgame

import (
	"time"
)

type AnimatorFrameCfg struct {
	Frames FrameProvider
	AnimatorCfgCommon
}

type AnimatorFrame struct {
	cfg AnimatorFrameCfg
	s   *SpriteBase

	clock *Clock

	curFrameDuration    time.Duration
	curFrameStartedTime time.Duration
}

func (af *AnimatorFrame) Animate() {
	af.checkToInit()
	elapsed := af.clock.Now() - af.curFrameStartedTime
	if elapsed < af.curFrameDuration {
		return
	}
	if af.setNextFrame() {
		return
	}
	af.s.DeleteAnimator(af)
	if af.cfg.AfterFinish != nil {
		af.cfg.AfterFinish()
	}
	if !af.cfg.KeepAliveWhenFinished {
		af.s.Mgr().AsyncDeleteSprite(af.s)
	}
}

func (af *AnimatorFrame) setNextFrame() (more bool) {
	var f Frame
	if f, af.curFrameDuration, more = af.cfg.Frames.Next(); !more {
		return false
	}
	if !af.s.Update(UpdateArg{
		F:   f,
		IBC: af.cfg.InBoundsCheckType,
		CD:  af.cfg.CollisionDetectionType}) {
		return false
	}
	af.curFrameStartedTime = af.clock.Now()
	if af.cfg.AfterUpdate != nil {
		af.cfg.AfterUpdate()
	}
	return true
}

func (af *AnimatorFrame) checkToInit() {
	if af.clock != nil {
		return
	}
	af.clock = af.s.Game().MasterClock
	if !af.setNextFrame() {
		panic("Frames cannot be empty")
	}
}

func NewAnimatorFrame(s *SpriteBase, c AnimatorFrameCfg) *AnimatorFrame {
	return &AnimatorFrame{cfg: c, s: s}
}
