package main

import (
	"strings"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cwin"
	"github.com/nsf/termbox-go"
)

var (
	alphaName   = "alpha"
	alphaImgTxt = strings.Trim(`
  ┃
-█-█-
`, "\n")
	alphaAttr         = cwin.ChAttr{Fg: termbox.ColorLightYellow}
	alphaBullet1Name  = "alpha_bullet1"
	alphaBullet1Attr  = cwin.ChAttr{Fg: termbox.ColorLightYellow}
	alphaBullet1Speed = cgame.ActionPerSec(30)
)

type spriteAlpha struct {
	*cgame.SpriteBase
	m         *myGame
	betaKills int
}

func (a *spriteAlpha) SetPosRelative(dx, dy int) {
	newR := a.W.Rect()
	newR.X += dx
	newR.Y += dy
	if _, r := a.m.winArena.ClientRect().Overlap(newR); r == newR {
		a.W.SetPosRelative(dx, dy)
	}
}

func (a *spriteAlpha) fireWeapon() {
	x := a.W.Rect().X + a.W.Rect().W/2
	y := a.W.Rect().Y - 1
	a.Mgr.AddEvent(cgame.NewSpriteEventCreate(newSpriteBullet1(
		a.m.g, a.m.winArena, alphaBullet1Name, alphaBullet1Attr, 0, -1, alphaBullet1Speed, x, y)))
}

func (a *spriteAlpha) Collided(other cgame.Sprite) {
	if other.Cfg().Name == alphaBullet1Name || other.Cfg().Name == alphaName {
		return
	}
	a.m.g.Pause()
	a.m.g.GameOver()
}
