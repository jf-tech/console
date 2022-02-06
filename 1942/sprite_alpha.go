package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cwin"
	"github.com/nsf/termbox-go"
)

var (
	alphaName  = "alpha"
	alphaFrame = cgame.FrameFromString(strings.Trim(`
  ┃
-█-█-
`, "\n"), cwin.ChAttr{Fg: termbox.ColorLightYellow})
	alphaBulletName = "alpha_bullet"
)

type spriteAlpha struct {
	*cgame.SpriteBase
	m          *myGame
	betaKills  int
	gammaKills int
	deltaKills int
	gpWeapon   *giftPack
}

func (a *spriteAlpha) SetPosRelative(dx, dy int) {
	newR := a.Win().Rect()
	newR.X += dx
	newR.Y += dy
	if overlapped, r := a.m.winArena.ClientRect().ToOrigin().Overlap(newR); overlapped && r == newR {
		a.Win().SetPosRelative(dx, dy)
	}
}

func (a *spriteAlpha) fireWeapon() {
	x := a.Win().Rect().X + a.Win().Rect().W/2
	y := a.Win().Rect().Y - 1
	if a.gpWeapon == nil || a.gpWeapon.remainingLife() <= 0 {
		a.gpWeapon = nil
		createBullet(a.m, alphaBulletName, alphaBulletAttr, 0, -1, alphaBulletSpeed, x, y)
	} else {
		switch a.gpWeapon.name {
		case gpShotgunName, gpShotgun2Name:
			pellets := 3
			if a.gpWeapon.name == gpShotgun2Name {
				pellets = 5
			}
			for i := -pellets / 2; i <= pellets/2; i++ {
				createBullet(a.m, alphaBulletName, alphaBulletAttr, i, -1, alphaBulletSpeed, x, y)
			}
		default:
			panic(fmt.Sprintf("unknown weapon name: %s", a.gpWeapon.name))
		}
	}
}

func (a *spriteAlpha) weaponStats() (name, remaining string) {
	if a.gpWeapon != nil {
		return a.gpWeapon.name, (a.gpWeapon.remainingLife() / time.Second * time.Second).String()
	}
	return "Basic Gun", "Infinite"
}

func (a *spriteAlpha) killStats() map[string]int {
	return map[string]int{
		betaName:  a.betaKills,
		gammaName: a.gammaKills,
		deltaName: a.deltaKills,
	}
}

func (a *spriteAlpha) Collided(other cgame.Sprite) {
	a.Win().ToBottom()
	switch other.Name() {
	case alphaBulletName:
	case giftPackName:
		switch other.(*spriteGiftPack).gpSym {
		case gpShotgunSym:
			a.gpWeapon = newGiftPackShotgun(a.m.g.MasterClock)
		case gpShotgun2Sym:
			a.gpWeapon = newGiftPackShotgun2(a.m.g.MasterClock)
		}
	default:
		if !a.m.invincible {
			a.m.g.GameOver()
		}
	}
}

func createAlpha(m *myGame) {
	m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(&spriteAlpha{
		SpriteBase: cgame.NewSpriteBase(m.g, m.winArena, alphaName, alphaFrame,
			(m.winArena.ClientRect().W-cgame.FrameRect(alphaFrame).W)/2,
			m.winArena.ClientRect().H-cgame.FrameRect(alphaFrame).H),
		m: m}))
}
