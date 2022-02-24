package cgame

import (
	"time"

	"github.com/jf-tech/console/cutil"
)

type AnimatorFrameCfg struct {
	Frames FrameProvider
	AnimatorCfgCommon
}

type AnimatorFrame struct {
	cfg AnimatorFrameCfg
	s   Sprite

	clock *cutil.Clock

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
	af.s.Base().DeleteAnimator(af)
	if !af.s.Base().IsDestroyed() && af.cfg.AfterFinish != nil {
		af.cfg.AfterFinish()
	}
	if !af.cfg.KeepAliveWhenFinished {
		af.s.Mgr().DeleteSprite(af.s)
	}
}

func (af *AnimatorFrame) setNextFrame() (more bool) {
	var f Frame
	if f, af.curFrameDuration, more = af.cfg.Frames.Next(); !more {
		return false
	}
	if !af.s.Base().Update(UpdateArg{
		F:   f,
		IBC: af.cfg.InBoundsCheckType,
		CD:  af.cfg.CollisionDetectionType}) {
		return false
	}
	af.curFrameStartedTime = af.clock.Now()
	if !af.s.Base().IsDestroyed() && af.cfg.AfterUpdate != nil {
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

func NewAnimatorFrame(s Sprite, c AnimatorFrameCfg) *AnimatorFrame {
	return &AnimatorFrame{cfg: c, s: s}
}
