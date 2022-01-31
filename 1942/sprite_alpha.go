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
	alphaAttr = cwin.ChAttr{Fg: termbox.ColorLightYellow}

	// https://textkool.com/en/ascii-art-generator?hl=default&vl=default&font=Colossal&text=Game%20Over%20!
	gameOverTxt = `

 .d8888b.                                        .d88888b.                                 888
d88P  Y88b                                      d88P" "Y88b                                888
888    888                                      888     888                                888
888         8888b.  88888b.d88b.   .d88b.       888     888 888  888  .d88b.  888d888      888
888  88888     "88b 888 "888 "88b d8P  Y8b      888     888 888  888 d8P  Y8b 888P"        888
888    888 .d888888 888  888  888 88888888      888     888 Y88  88P 88888888 888          Y8P
Y88b  d88P 888  888 888  888  888 Y8b.          Y88b. .d88P  Y8bd8P  Y8b.     888           "
 "Y8888P88 "Y888888 888  888  888  "Y8888        "Y88888P"    Y88P    "Y8888  888          888


                                     Press Enter to exit.
`
	alphaBullet1Name  = "alpha_bullet1"
	alphaBullet1Attr  = cwin.ChAttr{Fg: termbox.ColorLightYellow}
	alphaBullet1Speed = cgame.ActionPerSec(25)
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
	a.m.g.WinSys.MessageBox(nil, "Uh oh...", gameOverTxt)
	a.m.g.GameOver()
}
