package main

import (
	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cwin"
	"github.com/nsf/termbox-go"
)

var (
	bgStarName    = "background_star"
	bgStarImgTxt  = "."
	bgStarAttr    = cwin.ChAttr{Fg: termbox.ColorDarkGray}
	bgStarSpeed   = cgame.ActionPerSec(15)
	bgStarGenProb = 400
)

func newSpriteBackgroundStar(g *cgame.Game, parent *cwin.Win, x, y int) *cgame.SpriteAnimated {
	s := cgame.NewSpriteAnimated(g, parent,
		cgame.SpriteAnimatedCfg{
			Name: bgStarName,
			Frames: [][]cgame.Cell{
				cgame.StringToCells(bgStarImgTxt, bgStarAttr),
			},
			DY:        1,
			MoveSpeed: bgStarSpeed,
		},
		x, y)
	s.Win().ToBottom()
	return s
}
