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

func (b *spriteBullet) CollisionNotify(_ bool, _ []cgame.Sprite) cgame.CollisionResponseType {
	b.Mgr().AsyncDeleteSprite(b)
	return cgame.CollisionResponseJustDoIt
}

func createBullet(m *myGame, name string, attr cwin.Attr,
	dx, dy int, speed cgame.CharPerSec, x, y int) {

	s := &spriteBullet{cgame.NewSpriteBase(
		m.g, m.winArena, name, cgame.FrameFromString(bulletFrameTxt, attr), x, y)}

	dist := 1000 // large enough to go out of window (and auto destroy)
	a := cgame.NewAnimatorWaypoint(s.SpriteBase, cgame.AnimatorWaypointCfg{
		Waypoints: cgame.NewSimpleWaypoints([]cgame.Waypoint{
			{
				DX: dx * dist,
				DY: dy * dist,
				T:  time.Duration((float64(dist) / float64(speed)) * float64(time.Second)),
			}}),
	})
	s.AddAnimator(a)
	m.g.SpriteMgr.AsyncCreateSprite(s)
}
