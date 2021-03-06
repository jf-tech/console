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
`, cwin.Attr{Fg: cterm.ColorLightYellow})
	alphaBulletName = "alpha_bullet"
)

type spriteAlpha struct {
	*cgame.SpriteBase
	m          *myGame
	betaKills  int
	gammaKills int
	deltaKills int
	hits       int
	gpWeapon   *giftPack
}

// cgame.CollisionResponse
func (a *spriteAlpha) CollisionNotify(_ bool, collidees []cgame.Sprite) cgame.CollisionResponseType {
	hits := 0
	gp := cgame.Sprite(nil)
	for i := len(collidees) - 1; i >= 0; i-- {
		if collidees[i].Name() == giftPackName {
			if gp == nil {
				gp = collidees[i]
			}
		} else {
			hits++
		}
	}
	if hits > 0 {
		if !a.m.invincible {
			a.m.g.GameOver()
			return cgame.CollisionResponseJustDoIt
		}
		a.m.g.SoundMgr.PlayMP3(sfxOuchFile, sfxClipVol, 1)
		a.hits += hits
	}
	if gp != nil {
		switch gp.(*spriteGiftPack).gpSym {
		case gpShotgunSym:
			a.gpWeapon = newGiftPackShotgun(a.m.g.MasterClock)
		case gpShotgun2Sym:
			a.gpWeapon = newGiftPackShotgun2(a.m.g.MasterClock)
		}
		a.m.g.Exchange.GenericData[exchangeGiftPackWeapon] = a.gpWeapon
		a.m.g.SoundMgr.PlayMP3(sfxWeaponUpgradedFile, sfxClipVol, 1)
	}
	return cgame.CollisionResponseJustDoIt
}

func (a *spriteAlpha) move(dx, dy int) {
	if (dx == 0 && dy == 0) || (dx != 0 && dy != 0) {
		panic("one and only one of dx, dy be non-zero")
	}
	// Why not a simple Update call? Why so much trouble and copmlexity here?
	// For four directions, we don't always use 1 (or -1) as delta since in certain direction
	// we want our alpha to move a bit faster. This leads to a situation where the alpha can't
	// always reach the very edge of the arena. A hack here is to "nudge" the alpha back a little
	// bit if an update position fails (that is of course, only when the game is still not over).
	calcStep := func(delta int) int {
		step := 0
		if delta < 0 {
			step = 1
		} else if delta > 0 {
			step = -1
		}
		return step
	}
	xstep, ystep := calcStep(dx), calcStep(dy)
	for !a.Game().IsGameOver() && !a.Update(
		cgame.UpdateArg{DXY: &cwin.Point{X: dx, Y: dy}, IBC: cgame.InBoundsCheckFullyVisible}) {
		dx += xstep
		dy += ystep
	}
}

func (a *spriteAlpha) fireWeapon() {
	x := a.Rect().X + a.Rect().W/2
	y := a.Rect().Y - 1
	if a.gpWeapon == nil || a.gpWeapon.remainingLife() <= 0 {
		a.gpWeapon = nil
		delete(a.m.g.Exchange.GenericData, exchangeGiftPackWeapon)
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

func createAlpha(m *myGame, stage *stage) {
	gp := (*giftPack)(nil)
	if d, ok := m.g.Exchange.GenericData[exchangeGiftPackWeapon]; ok {
		gp = d.(*giftPack)
	}
	m.g.SpriteMgr.AddSprite(&spriteAlpha{
		SpriteBase: cgame.NewSpriteBase(m.g, m.winArena, alphaName, alphaFrame,
			(m.winArena.ClientRect().W-cgame.FrameRect(alphaFrame).W)/2,
			m.winArena.ClientRect().H-cgame.FrameRect(alphaFrame).H),
		m:        m,
		gpWeapon: gp})
}
