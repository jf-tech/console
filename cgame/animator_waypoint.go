package cgame

import (
	"time"
)

type WaypointType int

const (
	WaypointAbs WaypointType = iota
	WaypointRelative
)

// WaypointAbs: from current position to (X, Y) using time T
// WaypointRel: from current position to (curX + X, curY + Y) using time T. Note X=Y=0, means
// the animator will keep the sprite at the current position for time T.
type Waypoint struct {
	Type WaypointType
	X, Y int
	T    time.Duration
}

type AnimatorWaypointCfg struct {
	Waypoints               []Waypoint
	Loop                    bool
	KeepAliveWhenOutOfBound bool
	KeepAliveWhenFinished   bool
	AfterMove, AfterFinish  func(Sprite)
}

type AnimatorWaypoint struct {
	cfg   AnimatorWaypointCfg
	clock *Clock

	wpIdx                  int
	curWPOrigX, curWPOrigY int
	curWPDestX, curWPDestY int
	curWPStartedTime       time.Duration
}

func (aw *AnimatorWaypoint) Animate(s Sprite) AnimatorState {
	aw.checkToInit(s)
	wp := aw.cfg.Waypoints[aw.wpIdx]
	elapsed := aw.clock.Now() - aw.curWPStartedTime
	ratio := float64(1)
	if elapsed < wp.T {
		ratio = float64(elapsed) / float64(wp.T)
	}
	// move proportionally to the elapsed time over wp.T
	newX := aw.curWPOrigX + int(float64(aw.curWPDestX-aw.curWPOrigX)*ratio)
	newY := aw.curWPOrigY + int(float64(aw.curWPDestY-aw.curWPOrigY)*ratio)
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
	if elapsed < wp.T {
		return AnimatorRunning
	}
	aw.wpIdx++
	if aw.wpIdx >= len(aw.cfg.Waypoints) && aw.cfg.Loop {
		aw.wpIdx = 0
	}
	if aw.wpIdx < len(aw.cfg.Waypoints) {
		aw.initCurWaypoint(s)
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

func (aw *AnimatorWaypoint) initCurWaypoint(s Sprite) {
	curR := s.Win().Rect()
	aw.curWPOrigX, aw.curWPOrigY = curR.X, curR.Y
	wp := aw.cfg.Waypoints[aw.wpIdx]
	aw.curWPDestX, aw.curWPDestY = wp.X, wp.Y
	if wp.Type == WaypointRelative {
		aw.curWPDestX += curR.X
		aw.curWPDestY += curR.Y
	}
	aw.curWPStartedTime = aw.clock.Now()
}

func (aw *AnimatorWaypoint) checkToInit(s Sprite) {
	if aw.wpIdx < 0 {
		aw.clock = s.Game().MasterClock
		aw.wpIdx = 0
		aw.initCurWaypoint(s)
	}
}

func NewAnimatorWaypoint(c AnimatorWaypointCfg) *AnimatorWaypoint {
	return &AnimatorWaypoint{cfg: c, wpIdx: -1}
}
