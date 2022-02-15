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

func (b *spriteBullet) IsCollidableWith(other cgame.Collidable) bool {
	if b.Name() == alphaBulletName {
		switch other.Name() {
		case betaName, gammaName, deltaName, bossName:
			return true
		}
		return false
	}
	return other.Name() == alphaName
}

func (b *spriteBullet) ResolveCollision(other cgame.Collidable) cgame.CollisionResolution {
	b.Mgr().DeleteSprite(b)
	return cgame.CollisionAllowed
}

func createBullet(m *myGame, name string, attr cwin.ChAttr,
	dx, dy int, speed cgame.CharPerSec, x, y int) {

	s := &spriteBullet{cgame.NewSpriteBase(
		m.g, m.winArena, name, cgame.FrameFromString(bulletFrameTxt, attr), x, y)}

	dist := 1000 // large enough to go out of window (and auto destroy)
	a := cgame.NewAnimatorWaypoint(s.SpriteBase, cgame.AnimatorWaypointCfg{
		Waypoints: cgame.NewSimpleWaypoints([]cgame.Waypoint{
			{
				Type: cgame.WaypointRelative,
				X:    dx * dist,
				Y:    dy * dist,
				T:    time.Duration((float64(dist) / float64(speed)) * float64(time.Second)),
			}})})
	s.AddAnimator(a)
	m.g.SpriteMgr.AddSprite(s)
}
