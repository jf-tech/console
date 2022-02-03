package main

import (
	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cwin"
)

var (
	bulletImgTxt = '•'
)

type spriteBullet struct {
	*cgame.SpriteAnimated
}

func (b *spriteBullet) Collided(other cgame.Sprite) {
	if b.Name() == alphaBulletName {
		if other.Name() == betaName || other.Name() == gammaName {
			b.Mgr().AddEvent(cgame.NewSpriteEventDelete(b))
		}
	} else if other.Name() == alphaName {
		b.Mgr().AddEvent(cgame.NewSpriteEventDelete(b))
	}
}

func newSpriteBullet(g *cgame.Game, parent *cwin.Win, name string, attr cwin.ChAttr,
	dx, dy int, speed cgame.ActionPerSec, x, y int) *spriteBullet {
	return &spriteBullet{
		cgame.NewSpriteAnimated(g, parent,
			cgame.SpriteAnimatedCfg{
				Name: name,
				Frames: [][]cgame.Cell{
					cgame.StringToCells(string(bulletImgTxt), attr), // single frame
				},
				DX:        dx,
				DY:        dy,
				MoveSpeed: speed,
			},
			x, y)}
}
