package main

import (
	"math/rand"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cutil"
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

func (b *spriteBoss) CollisionNotify(_ bool, _ []cgame.Sprite) cgame.CollisionResponseType {
	b.hpLeft--
	b.Update(cgame.UpdateArg{
		F: createBossFrameWithHP(b.hpLeft),
		// must off, or we turn into infinite recursion, because whatever hits the boss (alpha
		// or alpha bullet) is still there, thus the frame update would still generate a collision
		// event which would call this, so on and so forth.
		CD: cgame.CollisionDetectionOff,
	})
	if b.hpLeft <= 0 {
		cgame.CreateExplosion(b.SpriteBase, cgame.ExplosionCfg{
			MaxDuration: bossExplosionDuration,
			SpriteName:  bossExplosionName,
		})
	}
	return cgame.CollisionResponseJustDoIt

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
	if cutil.CheckProbability(bossBulletFiringProb) {
		b.fireBulletSquare(curR.X+leftGunX, curR.Y+leftGunY)
	}
	// right gun
	if cutil.CheckProbability(bossBulletFiringProb) {
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
		AnimatorCfgCommon: cgame.AnimatorCfgCommon{
			AfterUpdate: func() {
				s.fireWeapon()
			},
		},
	})
	s.AddAnimator(a)
	m.g.SpriteMgr.AsyncCreateSprite(s)
}

type bossWaypoints struct {
	s *spriteBoss
}

func (bw *bossWaypoints) Next() (cgame.Waypoint, bool) {
	curR := bw.s.Rect()
	parentClientR := bw.s.ParentRect()
	if ro, overlapped := curR.Overlap(parentClientR); !overlapped || curR != ro {
		dist := -curR.Y
		// this is when the boss is still fully or partially out of the arena
		return cgame.Waypoint{
			DX: 0,
			DY: dist,
			T:  time.Duration(cgame.CharPerSec(abs(dist))/bossSpeed) * time.Second,
		}, true
	}
	for {
		dist := rand.Int() % (bossMaxDistToGoBeforeDirChange - bossMinDistToGoBeforeDirChange + 1)
		dist += bossMinDistToGoBeforeDirChange
		dirIdx := cwin.Dir(rand.Int() % cwin.DirCount)
		newR := bw.s.Rect()
		newR.X += cwin.DirOffSetXY[dirIdx].X * dist
		newR.Y += cwin.DirOffSetXY[dirIdx].Y * dist
		if ro, overlapped := newR.Overlap(parentClientR); overlapped && ro == newR {
			return cgame.Waypoint{
				DX: newR.X - bw.s.Rect().X,
				DY: newR.Y - bw.s.Rect().Y,
				T:  time.Duration(cgame.CharPerSec(dist)/bossSpeed) * time.Second,
			}, true
		}
	}
}
