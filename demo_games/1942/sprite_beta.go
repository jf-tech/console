package main

import (
	"math/rand"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
)

var (
	betaName  = "beta"
	betaFrame = cgame.FrameFromString(`
\┃┃/
 \/
`, cwin.ChAttr{Fg: cterm.ColorLightCyan})

	betaBulletName = "beta_bullet"
)

type spriteBeta struct {
	*cgame.SpriteBase
}

// cgame.CollisionResponse
func (b *spriteBeta) CollisionNotify(_ bool, _ []cgame.Sprite) cgame.CollisionResponseType {
	cgame.CreateExplosion(b.SpriteBase, cgame.ExplosionCfg{MaxDuration: betaExplosionDuration})
	b.Mgr().FindByName(alphaName).(*spriteAlpha).betaKills++
	return cgame.CollisionResponseJustDoIt
}

func createBeta(m *myGame, stageIdx int) {
	s := &spriteBeta{cgame.NewSpriteBase(m.g, m.winArena, betaName, betaFrame,
		rand.Int()%(m.winArena.ClientRect().W-cgame.FrameRect(betaFrame).W), 0)}
	dist := 1000 // large enough to go out of window (and auto destroy)
	a := cgame.NewAnimatorWaypoint(s.SpriteBase, cgame.AnimatorWaypointCfg{
		Waypoints: cgame.NewSimpleWaypoints([]cgame.Waypoint{
			{
				DX: 0,
				DY: 1 * dist,
				T:  time.Duration((float64(dist) / float64(betaSpeed)) * float64(time.Second)),
			}}),
		AnimatorCfgCommon: cgame.AnimatorCfgCommon{
			AfterUpdate: func() {
				if !cgame.CheckProbability(betaFiringProbPerStage[stageIdx]) {
					return
				}
				x := s.Rect().X + s.Rect().W/2
				y := s.Rect().Y + s.Rect().H
				pellets := betaFiringPelletsPerStage[stageIdx]
				if m.easyMode {
					pellets /= 2
				}
				for dx := -(pellets / 2); dx <= pellets/2; dx++ {
					if pellets%2 == 0 && dx == 0 {
						continue
					}
					createBullet(m, betaBulletName, enemyBulletAttr, dx, 1, betaBulletSpeed, x, y)
				}
			},
		},
	})
	s.AddAnimator(a)
	m.g.SpriteMgr.AddSprite(s)
}
