package main

import (
	"math/rand"
	"strings"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
	"github.com/jf-tech/go-corelib/maths"
)

var (
	bossName          = "boss"
	bossExplosionName = "boss_explosion"
	bossFrame         = cgame.FrameFromString(strings.Trim(`
[========================================]

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
`, "\n"), cwin.ChAttr{Fg: cterm.ColorWhite})
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

func (b *spriteBoss) Collided(other cgame.Sprite) {
	if other.Name() == alphaBulletName || other.Name() == alphaName {
		b.hpLeft--
		b.drawHP()
		if b.hpLeft <= 0 {
			b.Mgr().AddEvent(cgame.NewSpriteEventDelete(b))
			cgame.CreateExplosion(b, cgame.ExplosionCfg{
				MaxDuration: bossExplosionDuration,
				SpriteName:  bossExplosionName,
			})
		}
	}
}

func (b *spriteBoss) drawHP() {
	r := b.Win().Rect()
	hpStartX := 1
	hpEndX := r.W - 2
	hpFullLength := hpEndX - hpStartX + 1
	hpLen := int(float64(b.hpLeft) / float64(bossHP) * float64(hpFullLength))
	if b.hpLeft > 0 {
		hpLen = maths.MaxInt(1, hpLen)
	}
	b.Win().PutClient(0, 0, cwin.Chx{Ch: '[', Attr: bossHPAttr})
	for x := hpStartX; x <= hpEndX; x++ {
		if x <= hpLen {
			b.Win().PutClient(x, 0, cwin.Chx{Ch: '=', Attr: bossHPAttr})
		} else {
			b.Win().PutClient(x, 0, cwin.TransparentChx())
		}
	}
	b.Win().PutClient(r.W-1, 0, cwin.Chx{Ch: ']', Attr: bossHPAttr})
}

func (b *spriteBoss) fireWeapon() {
	curR := b.Win().Rect()
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
	s := &spriteBoss{
		SpriteBase: cgame.NewSpriteBase(m.g, m.winArena, bossName, bossFrame,
			rand.Int()%(m.winArena.ClientRect().W-cgame.FrameRect(bossFrame).W),
			-cgame.FrameRect(bossFrame).H+1),
		m:      m,
		hpLeft: bossHP}
	s.drawHP()
	a := cgame.NewAnimatorWaypoint(cgame.AnimatorWaypointCfg{
		Waypoints: &bossWaypoints{s: s},
		AfterMove: func(s cgame.Sprite) {
			s.(*spriteBoss).fireWeapon()
		},
	})
	m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(s, a))
}

type bossWaypoints struct {
	s *spriteBoss
}

func (bw *bossWaypoints) Next() (cgame.Waypoint, bool) {
	curR := bw.s.Win().Rect()
	parentClientR := bw.s.Win().Parent().ClientRect().ToOrigin()
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
		w := bw.s.Win()
		newR := w.Rect()
		newR.X += cgame.DirOffSetXY[dirIdx].A * dist
		newR.Y += cgame.DirOffSetXY[dirIdx].B * dist
		if overlapped, ro := newR.Overlap(w.Parent().ClientRect().ToOrigin()); overlapped && ro == newR {
			return cgame.Waypoint{
				Type: cgame.WaypointAbs,
				X:    newR.X,
				Y:    newR.Y,
				T:    time.Duration(cgame.CharPerSec(dist)/bossSpeed) * time.Second,
			}, true
		}
	}
}
