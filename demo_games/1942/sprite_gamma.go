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

// cgame.CollisionResponse
func (g *spriteGamma) CollisionNotify(_ bool, _ []cgame.Sprite) cgame.CollisionResponseType {
	cgame.CreateExplosion(g.SpriteBase, cgame.ExplosionCfg{MaxDuration: gammaExplosionDuration})
	g.Mgr().FindByName(alphaName).(*spriteAlpha).gammaKills++
	return cgame.CollisionResponseJustDoIt
}

func createGamma(m *myGame, stageIdx int) {
	s := &spriteGamma{cgame.NewSpriteBase(m.g, m.winArena, gammaName, gammaFrame,
		rand.Int()%(m.winArena.ClientRect().W-cgame.FrameRect(gammaFrame).W), 0)}
	dist := 1000 // large enough to go out of window (and auto destroy)
	a := cgame.NewAnimatorWaypoint(s.SpriteBase, cgame.AnimatorWaypointCfg{
		Waypoints: cgame.NewSimpleWaypoints([]cgame.Waypoint{
			{
				DX: 0,
				DY: 1 * dist,
				T:  time.Duration((float64(dist) / float64(gammaSpeed)) * float64(time.Second)),
			}}),
		AnimatorCfgCommon: cgame.AnimatorCfgCommon{
			AfterUpdate: func() {
				if !cgame.CheckProbability(gammaFiringProbPerStage[stageIdx]) {
					return
				}
				centerX := s.Rect().X + s.Rect().W/2
				centerY := s.Rect().Y + s.Rect().H/2
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
		},
	})
	s.AddAnimator(a)
	m.g.SpriteMgr.AsyncCreateSprite(s)
}
