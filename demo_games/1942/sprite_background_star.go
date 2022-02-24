package main

import (
	"math/rand"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
)

var (
	bgStarName  = "background_star"
	bgStarFrame = cgame.FrameFromString(".", cwin.Attr{Fg: cterm.ColorDarkGray})
)

func createBackgroundStar(m *myGame) {
	s := cgame.NewSpriteBase(m.g, m.winArena, bgStarName, bgStarFrame,
		rand.Int()%(m.winArena.ClientRect().W-cgame.FrameRect(bgStarFrame).W), 0)
	s.SendToBottom()
	dist := 1000 // large enough to go out of window (and auto destroy)
	a := cgame.NewAnimatorWaypoint(s, cgame.AnimatorWaypointCfg{
		Waypoints: cgame.NewWaypointProviderSimple([]cgame.Waypoint{
			{
				DX: 0,
				DY: 1 * dist,
				T:  time.Duration((float64(dist) / float64(bgStarSpeed)) * float64(time.Second)),
			}}),
	})
	s.AddAnimator(a)
	m.g.SpriteMgr.AddSprite(s)
}
