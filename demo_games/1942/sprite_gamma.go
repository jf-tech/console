package main

import (
	"math/rand"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
)

var (
	gammaName  = "gamma"
	gammaFrame = cgame.FrameFromString(`
/^#^\
\vvv/
`, cwin.ChAttr{Fg: cterm.ColorLightBlue})

	gammaBulletName = "gamma_bullet"
)

type spriteGamma struct {
	*cgame.SpriteBase
}

func (g *spriteGamma) Collided(other cgame.Sprite) {
	if other.Name() == alphaBulletName || other.Name() == alphaName {
		g.Mgr().AddEvent(cgame.NewSpriteEventDelete(g))
		cgame.CreateExplosion(g, cgame.ExplosionCfg{MaxDuration: gammaExplosionDuration})
		g.Mgr().FindByName(alphaName).(*spriteAlpha).gammaKills++
	}
}

func createGamma(m *myGame, stageIdx int) {
	dist := 1000 // large enough to go out of window (and auto destroy)
	a := cgame.NewAnimatorWaypoint(cgame.AnimatorWaypointCfg{
		Waypoints: cgame.NewSimpleWaypoints([]cgame.Waypoint{
			{
				Type: cgame.WaypointRelative,
				X:    0,
				Y:    1 * dist,
				T:    time.Duration((float64(dist) / float64(gammaSpeed)) * float64(time.Second)),
			}}),
		AfterMove: func(s cgame.Sprite) {
			if !cgame.CheckProbability(gammaFiringProbPerStage[stageIdx]) {
				return
			}
			centerX := s.Win().Rect().X + s.Win().Rect().W/2
			centerY := s.Win().Rect().Y + s.Win().Rect().H/2
			for y := -1; y <= 1; y++ {
				for x := -1; x <= 1; x++ {
					if x == 0 && y == 0 {
						continue
					}
					if m.easyMode && abs(x)+abs(y) == 1 {
						continue
					}
					createBullet(m, gammaBulletName, enemyBulletAttr,
						x, y, gammaBulletSpeed, centerX, centerY)
				}
			}
		},
	})
	s := &spriteGamma{cgame.NewSpriteBase(m.g, m.winArena, gammaName, gammaFrame,
		rand.Int()%(m.winArena.ClientRect().W-cgame.FrameRect(gammaFrame).W), 0)}
	m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(s, a))
}
