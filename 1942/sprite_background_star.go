package main

import (
	"math/rand"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cwin"
	"github.com/nsf/termbox-go"
)

var (
	bgStarName  = "background_star"
	bgStarFrame = cgame.FrameFromString(".", cwin.ChAttr{Fg: termbox.ColorDarkGray})
)

func createBackgroundStar(m *myGame) {
	dist := 1000 // large enough to go out of window (and auto destroy)
	a := cgame.NewAnimatorWaypoint(cgame.AnimatorWaypointCfg{
		Waypoints: []cgame.Waypoint{
			{
				Type: cgame.WaypointRelative,
				X:    0,
				Y:    1 * dist,
				T:    time.Duration((float64(dist) / float64(bgStarSpeed)) * float64(time.Second)),
			}},
		AfterMove: func(s cgame.Sprite) { s.Win().ToBottom() },
	})
	s := cgame.NewSpriteBase(m.g, m.winArena, bgStarName, bgStarFrame,
		rand.Int()%(m.winArena.ClientRect().W-cgame.FrameRect(bgStarFrame).W), 0)
	m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(s, a))
}
