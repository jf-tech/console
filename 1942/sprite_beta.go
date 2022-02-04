package main

import (
	"strings"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cwin"
	"github.com/jf-tech/go-corelib/maths"
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

	betaGenProb       = 10000
	betaFiringMinProb = 20
	betaFiringCurProb = betaFiringMinProb
	betaFiringMaxProb = 5

	betaBulletName  = "beta_bullet"
	betaBulletAttr  = cwin.ChAttr{Fg: termbox.ColorLightCyan}
	betaBulletSpeed = cgame.ActionPerSec(10)
)

type spriteBeta struct {
	*cgame.SpriteAnimated
}

func (b *spriteBeta) Collided(other cgame.Sprite) {
	if other.Name() == alphaBulletName || other.Name() == alphaName {
		b.Mgr().AddEvent(cgame.NewSpriteEventDelete(b))
		b.Mgr().AddEvent(cgame.NewSpriteEventCreate(
			newSpriteBetaDeath(b.Game(), b.Win().Parent(), b.Win().Rect().X, b.Win().Rect().Y)))
		b.Mgr().FindByName(alphaName).(*spriteAlpha).betaKills++
	}
}

func newSpriteBeta(g *cgame.Game, parent *cwin.Win, x, y int) *spriteBeta {
	s := &spriteBeta{}
	s.SpriteAnimated = cgame.NewSpriteAnimated(g, parent,
		cgame.SpriteAnimatedCfg{
			Name: betaName,
			Frames: [][]cgame.Cell{
				cgame.StringToCells(betaImgTxt, betaAttr), // single frame
			},
			DY:        1,
			MoveSpeed: betaSpeed,
			AfterMove: func(cgame.Sprite) {
				newFiringProb := maths.MaxInt(
					betaFiringMinProb-int(s.Clock().Now()/(5*time.Second)),
					betaFiringMaxProb)
				if newFiringProb < betaFiringCurProb {
					betaFiringCurProb = newFiringProb
				}
				if !testProb(betaFiringCurProb) {
					return
				}
				r := s.Win().Rect()
				for i := -1; i <= 1; i++ {
					s.Mgr().AddEvent(cgame.NewSpriteEventCreate(newSpriteBullet(
						g, parent, betaBulletName, betaBulletAttr,
						i, 1, betaBulletSpeed, r.X+r.W/2, r.Y+r.H)))
				}
			},
		},
		x, y)
	return s
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
