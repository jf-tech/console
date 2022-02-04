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
	gammaName   = "gamma"
	gammaImgTxt = strings.Trim(`
/^#^\
\vvv/
`, "\n")

	gammaAttr  = cwin.ChAttr{Fg: termbox.ColorLightBlue}
	gammaSpeed = betaSpeed

	gammaGenProb       = betaGenProb * 3
	gammaFiringMinProb = 30
	gammaFiringCurProb = gammaFiringMinProb
	gammaFiringMaxProb = 15

	gammaBulletName  = "gamma_bullet"
	gammaBulletAttr  = betaBulletAttr
	gammaBulletSpeed = betaBulletSpeed
)

type spriteGamma struct {
	*cgame.SpriteAnimated
}

func (g *spriteGamma) Collided(other cgame.Sprite) {
	if other.Name() == alphaBulletName || other.Name() == alphaName {
		g.Mgr().AddEvent(cgame.NewSpriteEventDelete(g))
		g.Mgr().AddEvent(cgame.NewSpriteEventCreate(
			newSpriteGammaDeath(g.Game(), g.Win().Parent(), g.Win().Rect().X, g.Win().Rect().Y)))
		a, _ := g.Mgr().FindByName(alphaName)
		a.(*spriteAlpha).gammaKills++
	}
}

func newSpriteGamma(g *cgame.Game, parent *cwin.Win, x, y int) *spriteGamma {
	s := &spriteGamma{}
	s.SpriteAnimated = cgame.NewSpriteAnimated(g, parent,
		cgame.SpriteAnimatedCfg{
			Name: gammaName,
			Frames: [][]cgame.Cell{
				cgame.StringToCells(gammaImgTxt, gammaAttr),
			},
			DY:        1,
			MoveSpeed: gammaSpeed,
			AfterMove: func(cgame.Sprite) {
				newFiringProb := maths.MaxInt(
					gammaFiringMinProb-int(s.Clock().Now()/(5*time.Second)),
					gammaFiringMaxProb)
				if newFiringProb < gammaFiringCurProb {
					gammaFiringCurProb = newFiringProb
				}
				if !testProb(gammaFiringCurProb) {
					return
				}
				r := s.Win().Rect()
				for y := -1; y <= 1; y++ {
					for x := -1; x <= 1; x++ {
						if x == 0 && y == 0 {
							continue
						}
						s.Mgr().AddEvent(cgame.NewSpriteEventCreate(newSpriteBullet(
							g, parent, gammaBulletName, gammaBulletAttr,
							x, y, gammaBulletSpeed, r.X+r.W/2, r.Y+r.H/2)))
					}
				}
			},
		},
		x, y)
	return s
}

var (
	gammaDeathName    = "gamma_death"
	gammaDeathImgTxts = []string{
		gammaImgTxt,
		strings.Trim(`
/^#^\
\vvv/
`, "\n"),
		strings.Trim(`
|^.^|
|.v.|
`, "\n"),
		strings.Trim(`
\'.'/
/ . \
`, "\n"),
		strings.Trim(`
~  '~
~'  ~
`, "\n"),
		strings.Trim(`
'   '
'   '
`, "\n"),
	}
	gammaDeathSpeed = betaDeathSpeed
)

type spriteGammaDeath struct {
	*cgame.SpriteAnimated
}

func newSpriteGammaDeath(g *cgame.Game, parent *cwin.Win, x, y int) *spriteGammaDeath {
	return &spriteGammaDeath{
		cgame.NewSpriteAnimated(g, parent,
			cgame.SpriteAnimatedCfg{
				Name:       gammaDeathName,
				Frames:     cgame.StringsToFrames(gammaDeathImgTxts, gammaAttr),
				FrameSpeed: gammaDeathSpeed,
			},
			x, y)}
}
