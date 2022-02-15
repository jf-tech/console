package main

import (
	"time"

	"github.com/jf-tech/console/cgame"
)

func createSlideInOutBanner(
	m *myGame, frame cgame.Frame, inOut, stay time.Duration, afterFinish func()) {
	frameR := cgame.FrameRect(frame)
	y := (m.winArena.ClientRect().H - frameR.H) / 2
	s := cgame.NewSpriteBase(m.g, m.winArena, "banner", frame, -frameR.W+1, y)
	a := cgame.NewAnimatorWaypoint(s, cgame.AnimatorWaypointCfg{
		Waypoints: cgame.NewSimpleWaypoints([]cgame.Waypoint{
			{
				Type: cgame.WaypointAbs,
				X:    (m.winArena.ClientRect().W - frameR.W) / 2,
				Y:    y,
				T:    inOut,
			},
			{
				Type: cgame.WaypointRelative,
				X:    0,
				Y:    0,
				T:    stay,
			},
			{
				Type: cgame.WaypointAbs,
				X:    m.winArena.ClientRect().W,
				Y:    y,
				T:    inOut,
			},
		}),
		AfterFinish: func() {
			afterFinish()
		},
	})
	s.AddAnimator(a)
	m.g.SpriteMgr.AddSprite(s)
}
