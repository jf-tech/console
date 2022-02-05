package main

import (
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cwin"
)

var (
	bulletFrameTxt = "•"
)

type spriteBullet struct {
	*cgame.SpriteBase
}

func (b *spriteBullet) Collided(other cgame.Sprite) {
	if b.Name() == alphaBulletName {
		if other.Name() == betaName || other.Name() == gammaName {
			b.Mgr().AddEvent(cgame.NewSpriteEventDelete(b))
		}
	} else if other.Name() == alphaName {
		b.Mgr().AddEvent(cgame.NewSpriteEventDelete(b))
	}
}

func createBullet(m *myGame, name string, attr cwin.ChAttr,
	dx, dy int, speed cgame.CharPerSec, x, y int) {
	dist := 1000 // large enough to go out of window (and auto destroy)
	a := cgame.NewAnimatorWaypoint(cgame.AnimatorWaypointCfg{Waypoints: []cgame.Waypoint{
		{
			Type: cgame.WaypointRelative,
			X:    dx * dist,
			Y:    dy * dist,
			T:    time.Duration((float64(dist) / float64(speed)) * float64(time.Second)),
		},
	}})
	s := &spriteBullet{cgame.NewSpriteBase(
		m.g, m.winArena, name, cgame.FrameFromString(bulletFrameTxt, attr), x, y)}
	m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(s, a))
}
