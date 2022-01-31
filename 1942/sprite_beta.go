package main

import (
	"math/rand"
	"strings"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cwin"
	"github.com/nsf/termbox-go"
)

var (
	betaName   = "beta"
	betaImgTxt = strings.Trim(`
\┃┃/
 \/
 `, "\n")
	betaAttr  = cwin.ChAttr{Fg: termbox.ColorLightCyan}
	betaSpeed = cgame.ActionPerSec(4)

	betaGenProb    = 10000
	betaFiringProb = 20

	betaBullet1Name  = "beta_bullet1"
	betaBullet1Attr  = cwin.ChAttr{Fg: termbox.ColorLightCyan}
	betaBullet1Speed = cgame.ActionPerSec(10)
)

type spriteBeta struct {
	*cgame.SpriteAnimated
}

func (b *spriteBeta) Collided(other cgame.Sprite) {
	if other.Cfg().Name == betaBullet1Name || other.Cfg().Name == betaName {
		return
	}
	if other.Cfg().Name == alphaBullet1Name || other.Cfg().Name == alphaName {
		b.Mgr.AddEvent(cgame.NewSpriteEventDelete(b))
		b.Mgr.FindByName(alphaName).(*spriteAlpha).betaKills++
	}
}

func shouldGenBeta() bool {
	return rand.Int()%betaGenProb == 0
}

func shouldBetaFireShot() bool {
	return rand.Int()%betaFiringProb == 0
}

func newSpriteBeta(g *cgame.Game, parent *cwin.Win, x, y int) *spriteBeta {
	return &spriteBeta{
		cgame.NewSpriteAnimated(g, parent,
			cgame.SpriteAnimatedCfg{
				Name: betaName,
				Frames: [][]cgame.Cell{
					cgame.StringToCells(betaImgTxt, betaAttr), // single frame
				},
				DY:        1,
				MoveSpeed: betaSpeed,
				AfterMove: func(s cgame.Sprite) {
					if !shouldBetaFireShot() {
						return
					}
					b := s.(*cgame.SpriteAnimated)
					r := b.W.Rect()
					b.Mgr.AddEvent(cgame.NewSpriteEventCreate(newSpriteBullet1(
						g, parent, betaBullet1Name, betaBullet1Attr,
						-1, 1, betaBullet1Speed, r.X+r.W/2, r.Y+r.H)))
					b.Mgr.AddEvent(cgame.NewSpriteEventCreate(newSpriteBullet1(
						g, parent, betaBullet1Name, betaBullet1Attr,
						0, 1, betaBullet1Speed, r.X+r.W/2, r.Y+r.H)))
					b.Mgr.AddEvent(cgame.NewSpriteEventCreate(newSpriteBullet1(
						g, parent, betaBullet1Name, betaBullet1Attr,
						1, 1, betaBullet1Speed, r.X+r.W/2, r.Y+r.H)))

				},
			},
			x, y)}
}
