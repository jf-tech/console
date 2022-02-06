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

	curWP                    Waypoint
	curWPStartX, curWPStartY int
	curWPDestX, curWPDestY   int
	curWPStartedTime         time.Duration
}

func (aw *AnimatorWaypoint) Animate(s Sprite) AnimatorState {
	aw.checkToInit(s)
	elapsed := aw.clock.Now() - aw.curWPStartedTime
	ratio := float64(1)
	if elapsed < aw.curWP.T {
		ratio = float64(elapsed) / float64(aw.curWP.T)
	}
	// move proportionally to the elapsed time over aw.curWP.T
	newX := aw.curWPStartX + int(float64(aw.curWPDestX-aw.curWPStartX)*ratio)
	newY := aw.curWPStartY + int(float64(aw.curWPDestY-aw.curWPStartY)*ratio)
	if s.Win().Rect().X != newX || s.Win().Rect().Y != newY {
		// only make this actual move if the newX/Y is different than current position.
		s.Win().SetPosAbs(newX, newY)
		if aw.cfg.AfterMove != nil {
			aw.cfg.AfterMove(s)
		}
	}
	// always checking in case the sprite was created out of bound before the animator even
	// started, but before the actual move occurs
	if !s.Win().VisibleInParentClientRect() {
		if !aw.cfg.KeepAliveWhenOutOfBound {
			s.Mgr().AddEvent(NewSpriteEventDelete(s))
			if aw.cfg.AfterFinish != nil {
				aw.cfg.AfterFinish(s)
			}
			return AnimatorCompleted
		}
	}
	if elapsed < aw.curWP.T {
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
	if aw.curWP, more = aw.cfg.Waypoints.Next(); !more {
		return false
	}
	curR := s.Win().Rect()
	aw.curWPStartX, aw.curWPStartY = curR.X, curR.Y
	aw.curWPDestX, aw.curWPDestY = aw.curWP.X, aw.curWP.Y
	if aw.curWP.Type == WaypointRelative {
		aw.curWPDestX += curR.X
		aw.curWPDestY += curR.Y
	}
	aw.curWPStartedTime = aw.clock.Now()
	return true
}

func (aw *AnimatorWaypoint) checkToInit(s Sprite) {
	if aw.clock == nil {
		aw.clock = s.Game().MasterClock
		if !aw.setupNextWaypoint(s) {
			panic("Waypoints cannot be empty")
		}
	}
}

func NewAnimatorWaypoint(c AnimatorWaypointCfg) *AnimatorWaypoint {
	return &AnimatorWaypoint{cfg: c}
}
