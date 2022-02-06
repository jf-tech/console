package main

import (
	"math/rand"
	"strings"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cwin"
	"github.com/nsf/termbox-go"
)

var (
	deltaName  = "delta"
	deltaFrame = cgame.FrameFromString(strings.Trim(`
   /^\
 << X >>
   \v/
`, "\n"), cwin.ChAttr{Fg: termbox.ColorLightRed})
)

type spriteDelta struct {
	*cgame.SpriteBase
}

func (d *spriteDelta) Collided(other cgame.Sprite) {
	if other.Name() == alphaBulletName || other.Name() == alphaName {
		d.Mgr().AddEvent(cgame.NewSpriteEventDelete(d))
		d.Mgr().FindByName(alphaName).(*spriteAlpha).deltaKills++
	}
}

func createDelta(m *myGame) {
	dist := 1000 // large enough to go out of window (and auto destroy)
	y := 0
	x := rand.Int() % (m.winArena.ClientRect().W - cgame.FrameRect(deltaFrame).W)
	a := cgame.NewAnimatorWaypoint(cgame.AnimatorWaypointCfg{
		Waypoints: []cgame.Waypoint{
			{
				Type: cgame.WaypointRelative,
				X:    0,
				Y:    1 * dist,
				T:    time.Duration((float64(dist) / float64(deltaSpeed)) * float64(time.Second)),
			}},
	})
	s := &spriteDelta{cgame.NewSpriteBase(m.g, m.winArena, deltaName, deltaFrame, x, y)}
	m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(s, a))
}
