package cgame

import (
	"time"

	"github.com/jf-tech/console/cutil"
	"github.com/jf-tech/console/cwin"
)

type AnimatorWaypointCfg struct {
	Waypoints WaypointProvider
	// instead of a "continuous" move, after the specified T of a waypoint, the move will be
	// completed in one go. In most cases, we like the smooth movement; but in some case, like
	// for example in tetris, we do want the tetris piece to drop in a discrete manner
	SingleMovePerWaypoint bool
	AnimatorCfgCommon
}

type AnimatorWaypoint struct {
	cfg AnimatorWaypointCfg
	s   *SpriteBase

	clock *cutil.Clock

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
			aw.s.Mgr().AsyncDeleteSprite(aw.s)
		}
	}

	elapsed := aw.clock.Now() - aw.wpStartedTime
	ratio := float64(1)
	if elapsed < aw.wp.T {
		if aw.cfg.SingleMovePerWaypoint {
			return
		}
		ratio = float64(elapsed) / float64(aw.wp.T)
	}
	// move proportionally to the elapsed time over current waypoint duration aw.wp.T
	dx, dy := int(float64(aw.wp.DX)*ratio), int(float64(aw.wp.DY)*ratio)
	if aw.dxDone != dx || aw.dyDone != dy {
		// If collision is detected or in-bounds check fails, and PreUpdateNotify decides to abandon
		// then the this animator is finished.
		if !aw.s.Update(UpdateArg{
			DXY: &cwin.Point{X: dx - aw.dxDone, Y: dy - aw.dyDone},
			IBC: aw.cfg.InBoundsCheckType,
			CD:  aw.cfg.CollisionDetectionType}) {
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
