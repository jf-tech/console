package main

import (
	"math/rand"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
	"github.com/jf-tech/go-corelib/maths"
)

var (
	bossName           = "boss"
	bossExplosionName  = "boss_explosion"
	bossFrameWithoutHP = cgame.FrameFromString(`
               |||      |||
               | |  __  | |
|-|_____-----/   |_|  |_|   \-----_____|-|
|_|_________{   }|  (^) |{  }__________|_|
 ||          |_| |   ^  | |_|          ||
 |              \|  /\  |/              |
 |               \ |--| /               |
 =               \ |__| /               =
 +               \      /               +
                  \    /
                  \    /
                   \  /
                   \  /
                   \  /
                    \/
`, cwin.ChAttr{Fg: cterm.ColorWhite})
	bossHPAttr     = cwin.ChAttr{Fg: cterm.ColorRed}
	bossBulletName = "boss_bullet"

	leftGunX  = 1
	leftGunY  = 10
	rightGunX = 40
	rightGunY = 10
)

type spriteBoss struct {
	*cgame.SpriteBase
	m      *myGame
	hpLeft int
}

func (b *spriteBoss) IsCollidableWith(other cgame.Collidable) bool {
	return other.Name() == alphaBulletName || other.Name() == alphaName
}

func (b *spriteBoss) ResolveCollision(other cgame.Collidable) cgame.CollisionResolution {
	b.hpLeft--
	// TODO/BUG, infinite recursion.
	b.SetFrame(createBossFrameWithHP(b.hpLeft))
	// a bullet hits me, so during collision resolution i'm setting up the new frame.
	// since the bullet still here, the action of setting up new frame causes collision
	// so on and so forth.
	if b.hpLeft <= 0 {
		cgame.CreateExplosion(b.SpriteBase, cgame.ExplosionCfg{
			MaxDuration: bossExplosionDuration,
			SpriteName:  bossExplosionName,
		})
	}
	return cgame.CollisionAllowed
}

func createBossFrameWithHP(hpLeft int) cgame.Frame {
	r := cgame.FrameRect(bossFrameWithoutHP)
	f := cgame.CopyFrame(bossFrameWithoutHP)
	for i := 0; i < len(f); i++ {
		f[i].Y += 2
	}
	f = append(f, cgame.Cell{X: 0, Y: 0, Chx: cwin.Chx{Ch: '[', Attr: bossHPAttr}})
	f = append(f, cgame.Cell{X: r.W - 1, Y: 0, Chx: cwin.Chx{Ch: ']', Attr: bossHPAttr}})
	hpStartX := 1
	hpEndX := r.W - 2
	hpFullLength := hpEndX - hpStartX + 1
	hpLen := int(float64(hpLeft) / float64(bossHP) * float64(hpFullLength))
	if hpLeft > 0 {
		hpLen = maths.MaxInt(1, hpLen)
	}
	for i := 0; i < hpLen; i++ {
		f = append(f, cgame.Cell{X: hpStartX + i, Y: 0, Chx: cwin.Chx{Ch: '=', Attr: bossHPAttr}})
	}
	return f
}

func (b *spriteBoss) fireWeapon() {
	curR := b.Rect()
	// left gun
	if cgame.CheckProbability(bossBulletFiringProb) {
		b.fireBulletSquare(curR.X+leftGunX, curR.Y+leftGunY)
	}
	// right gun
	if cgame.CheckProbability(bossBulletFiringProb) {
		b.fireBulletSquare(curR.X+rightGunX, curR.Y+rightGunY)
	}
}

func (b *spriteBoss) fireBulletSquare(x, y int) {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			if b.m.easyMode && abs(dx)+abs(dy) == 1 {
				continue
			}
			createBullet(b.m, bossBulletName, enemyBulletAttr,
				dx, dy, bossBulletSpeed, x, y)
		}
	}
}

func createBoss(m *myGame) {
	f := createBossFrameWithHP(bossHP)
	s := &spriteBoss{
		SpriteBase: cgame.NewSpriteBase(m.g, m.winArena, bossName, f,
			rand.Int()%(m.winArena.ClientRect().W-cgame.FrameRect(f).W),
			-cgame.FrameRect(f).H+1),
		m:      m,
		hpLeft: bossHP}
	a := cgame.NewAnimatorWaypoint(s.SpriteBase, cgame.AnimatorWaypointCfg{
		Waypoints: &bossWaypoints{s: s},
		AfterMove: func() {
			s.fireWeapon()
		},
	})
	s.AddAnimator(a)
	m.g.SpriteMgr.AddSprite(s)
}

type bossWaypoints struct {
	s *spriteBoss
}

func (bw *bossWaypoints) Next() (cgame.Waypoint, bool) {
	curR := bw.s.Rect()
	parentClientR := bw.s.m.winArena.ClientRect().ToOrigin()
	if overlapped, ro := curR.Overlap(parentClientR); !overlapped || curR != ro {
		dist := -curR.Y
		// this is when the boss is still fully or partially out of the arena
		return cgame.Waypoint{
			Type: cgame.WaypointRelative,
			X:    0,
			Y:    dist,
			T:    time.Duration(cgame.CharPerSec(abs(dist))/bossSpeed) * time.Second,
		}, true
	}
	for {
		dist := rand.Int() % (bossMaxDistToGoBeforeDirChange - bossMinDistToGoBeforeDirChange + 1)
		dist += bossMinDistToGoBeforeDirChange
		dirIdx := rand.Int() % len(cgame.DirOffSetXY)
		newR := bw.s.Rect()
		newR.X += cgame.DirOffSetXY[dirIdx].X * dist
		newR.Y += cgame.DirOffSetXY[dirIdx].Y * dist
		if overlapped, ro := newR.Overlap(parentClientR); overlapped && ro == newR {
			return cgame.Waypoint{
				Type: cgame.WaypointAbs,
				X:    newR.X,
				Y:    newR.Y,
				T:    time.Duration(cgame.CharPerSec(dist)/bossSpeed) * time.Second,
			}, true
		}
	}
}
