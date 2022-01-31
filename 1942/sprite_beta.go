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
	g *cgame.Game
}

func (b *spriteBeta) Collided(other cgame.Sprite) {
	if other.Cfg().Name == betaBullet1Name || other.Cfg().Name == betaName {
		return
	}
	if other.Cfg().Name == alphaBullet1Name || other.Cfg().Name == alphaName {
		b.Mgr.AddEvent(cgame.NewSpriteEventDelete(b))
		b.Mgr.AddEvent(cgame.NewSpriteEventCreate(
			newSpriteBetaDeath(b.g, b.W.Parent(), b.W.Rect().X, b.W.Rect().Y)))
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
			x, y),
		g}
}

var (
	betaDeathName    = "beta_death"
	betaDeathImgTxts = []string{
		betaImgTxt,
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
	betaDeathSpeed = cgame.ActionPerSec(5)
)

type spriteBetaDeath struct {
	*cgame.SpriteAnimated
}

func newSpriteBetaDeath(g *cgame.Game, parent *cwin.Win, x, y int) *spriteBetaDeath {
	return &spriteBetaDeath{
		cgame.NewSpriteAnimated(g, parent,
			cgame.SpriteAnimatedCfg{
				Name:       betaDeathName,
				Frames:     cgame.StringsToFrames(betaDeathImgTxts, betaAttr),
				FrameSpeed: betaDeathSpeed,
			},
			x, y)}
}
