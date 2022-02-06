package main

import (
	"time"

	"github.com/jf-tech/console/cgame"
)

func createSlideInOutBanner(
	m *myGame, frame cgame.SpriteFrame, inOut, stay time.Duration, afterFinish func()) {

	frameR := cgame.FrameRect(frame)
	y := (m.winArena.ClientRect().H - frameR.H) / 2
	a := cgame.NewAnimatorWaypoint(cgame.AnimatorWaypointCfg{
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
		AfterFinish: func(cgame.Sprite) {
			afterFinish()
		},
	})
	s := cgame.NewSpriteBase(m.g, m.winArena, "banner", frame, -frameR.W+1, y)
	m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(s, a))
}
