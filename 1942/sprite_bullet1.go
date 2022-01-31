package main

import (
	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cwin"
)

var (
	bullet1ImgTxt = "•"
)

type spriteBullet1 struct {
	*cgame.SpriteAnimated
}

func (b1 *spriteBullet1) Collided(other cgame.Sprite) {
	if b1.Cfg().Name == alphaBullet1Name {
		if other.Cfg().Name == betaName {
			b1.Mgr.AddEvent(cgame.NewSpriteEventDelete(b1))
		}
	} else if b1.Cfg().Name == betaBullet1Name {
		if other.Cfg().Name == alphaName {
			b1.Mgr.AddEvent(cgame.NewSpriteEventDelete(b1))
		}
	}
}

func newSpriteBullet1(g *cgame.Game, parent *cwin.Win,
	name string, attr cwin.ChAttr,
	dx, dy int, speed cgame.ActionPerSec, x, y int) *spriteBullet1 {
	return &spriteBullet1{
		cgame.NewSpriteAnimated(g, parent,
			cgame.SpriteAnimatedCfg{
				Name: name,
				Frames: [][]cgame.Cell{
					cgame.StringToCells(bullet1ImgTxt, attr), // single frame
				},
				DX:        dx,
				DY:        dy,
				MoveSpeed: speed,
			},
			x, y)}
}
