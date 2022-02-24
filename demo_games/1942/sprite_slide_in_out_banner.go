package main

import (
	"time"

	"github.com/jf-tech/console/cgame"
)

func createSlideInOutBanner(
	m *myGame, frame cgame.Frame, inOut, stay time.Duration, afterFinish func()) {
	frameR := cgame.FrameRect(frame)
	y := (m.winArena.ClientRect().H - frameR.H) / 2
	s := cgame.NewSpriteBase(m.g, m.winArena, "banner", frame, -frameR.W, y)
	a := cgame.NewAnimatorWaypoint(s, cgame.AnimatorWaypointCfg{
		Waypoints: cgame.NewWaypointProviderSimple([]cgame.Waypoint{
			{
				DX: frameR.W + (m.winArena.ClientRect().W-frameR.W)/2,
				DY: 0,
				T:  inOut,
			},
			{
				DX: 0,
				DY: 0,
				T:  stay,
			},
			{
				DX: m.winArena.ClientRect().W - (m.winArena.ClientRect().W-frameR.W)/2,
				DY: 0,
				T:  inOut,
			},
		}),
		AnimatorCfgCommon: cgame.AnimatorCfgCommon{
			AfterFinish: func() {
				afterFinish()
			},
		},
	})
	s.AddAnimator(a)
	m.g.SpriteMgr.AddSprite(s)
}
