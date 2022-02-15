package main

import (
	"math/rand"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
)

var (
	deltaName  = "delta"
	deltaFrame = cgame.FrameFromString(`
  /^\
<< X >>
  \v/
`, cwin.ChAttr{Fg: cterm.ColorLightGreen})
)

type spriteDelta struct {
	*cgame.SpriteBase
}

// cgame.CollisionResponse
func (d *spriteDelta) CollisionNotify(_ bool, _ []cgame.Sprite) cgame.CollisionResponseType {
	cgame.CreateExplosion(d.SpriteBase, cgame.ExplosionCfg{MaxDuration: deltaExplosionDuration})
	d.Mgr().FindByName(alphaName).(*spriteAlpha).deltaKills++
	return cgame.CollisionResponseJustDoIt
}

func createDelta(m *myGame) {
	if cgame.CheckProbability(deltaVerticalProb) {
		createVerticalDelta(m)
		return
	}
	createHorizontalDelta(m)
}

func createVerticalDelta(m *myGame) {
	y := 0
	x := rand.Int() % (m.winArena.ClientRect().W - cgame.FrameRect(deltaFrame).W)
	s := &spriteDelta{cgame.NewSpriteBase(m.g, m.winArena, deltaName, deltaFrame, x, y)}
	vspeed := deltaVerticalSpeed
	if m.easyMode {
		vspeed *= cgame.CharPerSec(deltaSpeedDiscountEasy)
	}
	dist := 1000 // large enough to go out of window (and auto destroy)
	a := cgame.NewAnimatorWaypoint(s.SpriteBase, cgame.AnimatorWaypointCfg{
		Waypoints: cgame.NewSimpleWaypoints([]cgame.Waypoint{
			{
				Type: cgame.WaypointRelative,
				X:    0,
				Y:    1 * dist,
				T:    time.Duration((float64(dist) / float64(vspeed)) * float64(time.Second)),
			}})})
	s.AddAnimator(a)
	m.g.SpriteMgr.AddSprite(s)
}

func createHorizontalDelta(m *myGame) {
	x := -cgame.FrameRect(deltaFrame).W + 1
	y := rand.Int() % (m.winArena.ClientRect().H - cgame.FrameRect(deltaFrame).H)
	s := &spriteDelta{cgame.NewSpriteBase(m.g, m.winArena, deltaName, deltaFrame, x, y)}
	dist := 1000 // large enough to go out of window (and auto destroy)
	if cgame.CheckProbability("50%") {
		x = m.winArena.ClientRect().W - 1
		dist = -dist
	}
	a := cgame.NewAnimatorWaypoint(s.SpriteBase, cgame.AnimatorWaypointCfg{
		Waypoints: cgame.NewSimpleWaypoints([]cgame.Waypoint{
			{
				Type: cgame.WaypointRelative,
				X:    dist,
				Y:    0,
				T:    time.Duration((float64(abs(dist)) / float64(deltaHorizontalSpeed)) * float64(time.Second)),
			}})})
	s.AddAnimator(a)
	m.g.SpriteMgr.AddSprite(s)
}
