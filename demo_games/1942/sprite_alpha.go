package main

import (
	"fmt"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
)

var (
	alphaName  = "alpha"
	alphaFrame = cgame.FrameFromString(`
  ┃
-█-█-
`, cwin.ChAttr{Fg: cterm.ColorLightYellow})
	alphaBulletName = "alpha_bullet"
)

type spriteAlpha struct {
	*cgame.SpriteBase
	m          *myGame
	betaKills  int
	gammaKills int
	deltaKills int
	hit        int
	gpWeapon   *giftPack
}

func (a *spriteAlpha) fireWeapon() {
	x := a.Rect().X + a.Rect().W/2
	y := a.Rect().Y - 1
	if a.gpWeapon == nil || a.gpWeapon.remainingLife() <= 0 {
		a.gpWeapon = nil
		delete(a.m.g.Exchange.StringData, exchangeGiftPackWeapon)
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
	a.m.g.SoundMgr.PlayMP3(sfxPewFile, sfxClipVol, 1)
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

func (a *spriteAlpha) IsCollidableWith(other cgame.Collidable) bool {
	return other.Name() != alphaBulletName
}

func (a *spriteAlpha) ResolveCollision(other cgame.Collidable) cgame.CollisionResolution {
	ret := cgame.CollisionAllowed
	switch other.Name() {
	case giftPackName:
		switch other.(*spriteGiftPack).gpSym {
		case gpShotgunSym:
			a.gpWeapon = newGiftPackShotgun(a.m.g.MasterClock)
		case gpShotgun2Sym:
			a.gpWeapon = newGiftPackShotgun2(a.m.g.MasterClock)
		}
		a.m.g.Exchange.GenericData[exchangeGiftPackWeapon] = a.gpWeapon
		a.m.g.SoundMgr.PlayMP3(sfxWeaponUpgradedFile, sfxClipVol, 1)
	default:
		a.hit++
		if !a.m.invincible {
			a.m.g.GameOver()
			return ret
		}
		a.m.g.SoundMgr.PlayMP3(sfxOuchFile, sfxClipVol, 1)
	}
	return ret
}

func createAlpha(m *myGame, stage *stage) {
	gp := (*giftPack)(nil)
	if d, ok := m.g.Exchange.GenericData[exchangeGiftPackWeapon]; ok {
		gp = d.(*giftPack)
	}
	alpha := &spriteAlpha{
		SpriteBase: cgame.NewSpriteBase(m.g, m.winArena, alphaName, alphaFrame,
			(m.winArena.ClientRect().W-cgame.FrameRect(alphaFrame).W)/2,
			m.winArena.ClientRect().H-cgame.FrameRect(alphaFrame).H),
		m:        m,
		gpWeapon: gp}
	alpha.SetParentRectBound(true)
	m.g.SpriteMgr.AddSprite(alpha)
}
