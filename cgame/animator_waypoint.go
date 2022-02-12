package cgame

import (
	"time"
)

type AnimatorWaypointCfg struct {
	Waypoints               WaypointProvider
	KeepAliveWhenOutOfBound bool
	KeepAliveWhenFinished   bool
	AfterMove, AfterFinish  func(Sprite)
}

type AnimatorWaypoint struct {
	cfg AnimatorWaypointCfg

	clock *Clock

	wp             Waypoint
	dxDone, dyDone int
	wpStartedTime  time.Duration
}

func (aw *AnimatorWaypoint) Animate(s Sprite) AnimatorState {
	aw.checkToInit(s)
	elapsed := aw.clock.Now() - aw.wpStartedTime
	ratio := float64(1)
	if elapsed < aw.wp.T {
		ratio = float64(elapsed) / float64(aw.wp.T)
	}
	// move proportionally to the elapsed time over aw.curWP.T
	dx, dy := int(float64(aw.wp.X)*ratio), int(float64(aw.wp.Y)*ratio)
	if aw.dxDone != dx || aw.dyDone != dy {
		s.Win().SetPosRelative(dx-aw.dxDone, dy-aw.dyDone)
		aw.dxDone, aw.dyDone = dx, dy
		if aw.cfg.AfterMove != nil {
			aw.cfg.AfterMove(s)
		}
	}
	if !s.Win().VisibleInParentClientRect() {
		if !aw.cfg.KeepAliveWhenOutOfBound {
			s.Mgr().AddEvent(NewSpriteEventDelete(s))
			if aw.cfg.AfterFinish != nil {
				aw.cfg.AfterFinish(s)
			}
			return AnimatorCompleted
		}
	}
	if elapsed < aw.wp.T {
		return AnimatorRunning
	}
	if aw.setupNextWaypoint(s) {
		return AnimatorRunning
	}
	if !aw.cfg.KeepAliveWhenFinished {
		s.Mgr().AddEvent(NewSpriteEventDelete(s))
	}
	if aw.cfg.AfterFinish != nil {
		aw.cfg.AfterFinish(s)
	}
	return AnimatorCompleted
}

func (aw *AnimatorWaypoint) setupNextWaypoint(s Sprite) (more bool) {
	if aw.wp, more = aw.cfg.Waypoints.Next(); !more {
		return false
	}
	if aw.wp.Type == WaypointAbs {
		aw.wp.X -= s.Win().Rect().X
		aw.wp.Y -= s.Win().Rect().Y
		aw.wp.Type = WaypointRelative
	}
	aw.dxDone, aw.dyDone = 0, 0
	aw.wpStartedTime = aw.clock.Now()
	return true
}

func (aw *AnimatorWaypoint) checkToInit(s Sprite) {
	if aw.clock != nil {
		return
	}
	aw.clock = s.Game().MasterClock
	if !aw.setupNextWaypoint(s) {
		panic("Waypoints cannot be empty")
	}
}

func NewAnimatorWaypoint(c AnimatorWaypointCfg) *AnimatorWaypoint {
	return &AnimatorWaypoint{cfg: c}
}
