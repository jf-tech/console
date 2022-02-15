package cgame

import (
	"time"

	"github.com/jf-tech/console/cwin"
)

type AnimatorWaypointCfg struct {
	Waypoints WaypointProvider
	AnimatorCfgCommon
}

type AnimatorWaypoint struct {
	cfg AnimatorWaypointCfg
	s   *SpriteBase

	clock *Clock

	wp             Waypoint
	dxDone, dyDone int
	wpStartedTime  time.Duration
}

func (aw *AnimatorWaypoint) Animate() {
	aw.checkToInit()

	finish := func() {
		aw.s.DeleteAnimator(aw)
		if aw.cfg.AfterFinish != nil {
			aw.cfg.AfterFinish()
		}
		if !aw.cfg.KeepAliveWhenFinished {
			aw.s.Mgr().DeleteSprite(aw.s)
		}
	}

	elapsed := aw.clock.Now() - aw.wpStartedTime
	ratio := float64(1)
	if elapsed < aw.wp.T {
		ratio = float64(elapsed) / float64(aw.wp.T)
	}
	// move proportionally to the elapsed time over current waypoint duration aw.wp.T
	dx, dy := int(float64(aw.wp.X)*ratio), int(float64(aw.wp.Y)*ratio)
	if aw.dxDone != dx || aw.dyDone != dy {
		// If collision is detected or in-bounds check fails, and PreUpdateNotify decides to abandon
		// then the this animator is finished.
		if !aw.s.Update(UpdateArg{
			DXY:    &cwin.Point{X: dx - aw.dxDone, Y: dy - aw.dyDone},
			IBC:    aw.cfg.InBoundsCheckTypeToFinish,
			CD:     aw.cfg.CollisionDetectionTypeToFinish,
			Notify: aw.cfg.PreUpdateNotify}) {
			finish()
			return
		}
		aw.dxDone, aw.dyDone = dx, dy
		if aw.cfg.AfterUpdate != nil {
			aw.cfg.AfterUpdate()
		}
	}
	if elapsed < aw.wp.T {
		return
	}
	if aw.setupNextWaypoint() {
		return
	}
	finish()
}

func (aw *AnimatorWaypoint) setupNextWaypoint() (more bool) {
	if aw.wp, more = aw.cfg.Waypoints.Next(); !more {
		return false
	}
	if aw.wp.Type == WaypointAbs {
		aw.wp.X -= aw.s.Rect().X
		aw.wp.Y -= aw.s.Rect().Y
		aw.wp.Type = WaypointRelative
	}
	aw.dxDone, aw.dyDone = 0, 0
	aw.wpStartedTime = aw.clock.Now()
	return true
}

func (aw *AnimatorWaypoint) checkToInit() {
	if aw.clock != nil {
		return
	}
	aw.clock = aw.s.Game().MasterClock
	if !aw.setupNextWaypoint() {
		panic("Waypoints cannot be empty")
	}
}

func NewAnimatorWaypoint(s *SpriteBase, c AnimatorWaypointCfg) *AnimatorWaypoint {
	return &AnimatorWaypoint{cfg: c, s: s}
}
