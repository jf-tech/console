package main

import (
	"fmt"
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
	alphaAttr        = cwin.ChAttr{Fg: termbox.ColorLightYellow}
	alphaBulletName  = "alpha_bullet"
	alphaBulletAttr  = cwin.ChAttr{Fg: termbox.ColorLightYellow}
	alphaBulletSpeed = cgame.ActionPerSec(30)
)

type spriteAlpha struct {
	*cgame.SpriteBase
	m          *myGame
	betaKills  int
	gammaKills int
	gpWeapon   *giftPack
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
	if a.gpWeapon == nil || a.gpWeapon.remainingLife() <= 0 {
		a.gpWeapon = nil
		a.Mgr.AddEvent(cgame.NewSpriteEventCreate(newSpriteBullet(
			a.m.g, a.m.winArena, alphaBulletName, alphaBulletAttr, 0, -1, alphaBulletSpeed, x, y)))
	} else {
		switch a.gpWeapon.name {
		case gpShotgunName, gpShotgun2Name:
			pellets := 3
			if a.gpWeapon.name == gpShotgun2Name {
				pellets = 5
			}
			for i := -pellets / 2; i <= pellets/2; i++ {
				a.Mgr.AddEvent(cgame.NewSpriteEventCreate(newSpriteBullet(
					a.m.g, a.m.winArena, alphaBulletName, alphaBulletAttr, i, -1, alphaBulletSpeed, x, y)))
			}
		default:
			panic(fmt.Sprintf("unknown weapon name: %s", a.gpWeapon.name))
		}
	}
}

func (a *spriteAlpha) displayWeaponInfo() {
	name := "Basic"
	remaining := "Infinite"
	if a.gpWeapon != nil && a.gpWeapon.remainingLife() > 0 {
		name = a.gpWeapon.name
		remaining = a.gpWeapon.lifeRemaining.String()
	}
	a.m.winWeapon.SetText("WEAPON: %s  REMAINING TIME: %s", name, remaining)
}

func (a *spriteAlpha) Collided(other cgame.Sprite) {
	switch other.Cfg().Name {
	case alphaBulletName:
	case giftPackSpriteName:
		switch other.(*spriteGiftPack).gpSym {
		case gpShotgunSym:
			a.gpWeapon = newGiftPackShotgun(a.m.g.Clock)
		case gpShotgun2Sym:
			a.gpWeapon = newGiftPackShotgun2(a.m.g.Clock)
		}
	default:
		a.m.g.Pause()
		a.m.g.GameOver()
	}
}

func (a *spriteAlpha) displayKills() {
	a.m.winKills.SetText("KILLS: Beta:%d  Gamma:%d", a.betaKills, a.gammaKills)
}
