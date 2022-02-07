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
	betaName  = "beta"
	betaFrame = cgame.FrameFromString(strings.Trim(`
\┃┃/
 \/
`, "\n"), cwin.ChAttr{Fg: termbox.ColorLightCyan})

	betaBulletName = "beta_bullet"
)

type spriteBeta struct {
	*cgame.SpriteBase
}

func (b *spriteBeta) Collided(other cgame.Sprite) {
	if other.Name() == alphaBulletName || other.Name() == alphaName {
		b.Mgr().AddEvent(cgame.NewSpriteEventDelete(b))
		cgame.CreateExplosion(b, cgame.ExplosionCfg{MaxDuration: betaExplosionDuration})
		b.Mgr().FindByName(alphaName).(*spriteAlpha).betaKills++
	}
}

func createBeta(m *myGame, stageIdx int) {
	dist := 1000 // large enough to go out of window (and auto destroy)
	a := cgame.NewAnimatorWaypoint(cgame.AnimatorWaypointCfg{
		Waypoints: cgame.NewSimpleWaypoints([]cgame.Waypoint{
			{
				Type: cgame.WaypointRelative,
				X:    0,
				Y:    1 * dist,
				T:    time.Duration((float64(dist) / float64(betaSpeed)) * float64(time.Second)),
			}}),
		AfterMove: func(s cgame.Sprite) {
			if !cgame.CheckProbability(betaFiringProbPerStage[stageIdx]) {
				return
			}
			x := s.Win().Rect().X + s.Win().Rect().W/2
			y := s.Win().Rect().Y + s.Win().Rect().H
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
	})
	s := &spriteBeta{cgame.NewSpriteBase(m.g, m.winArena, betaName, betaFrame,
		rand.Int()%(m.winArena.ClientRect().W-cgame.FrameRect(betaFrame).W), 0)}
	m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(s, a))
}

/*
var (
	betaDeathName    = "beta_death"
	betaDeathFrameTxts = []string{
		betaFrameTxt,
		strings.Trim(`
_┃┃_
 \/
`, "\n"),
		strings.Trim(`
_\/_
 \/
`, "\n"),
		strings.Trim(`
_\/_
 ||
`, "\n"),
		strings.Trim(`
_\/_
 /\
`, "\n"),
		strings.Trim(`
'  '
'  '
`, "\n"),
	}
	betaDeathSpeed = cgame.CharPerSec(5)
)

type spriteBetaDeath struct {
	*cgame.SpriteAnimated
}

func newSpriteBetaDeath(g *cgame.Game, parent *cwin.Win, x, y int) *spriteBetaDeath {
	return &spriteBetaDeath{
		cgame.NewSpriteAnimated(g, parent,
			cgame.SpriteAnimatedCfg{
				Name:       betaDeathName,
				Frames:     cgame.FramesFromString(betaDeathImgTxts, betaAttr),
				FrameSpeed: betaDeathSpeed,
			},
			x, y)}
}
*/
