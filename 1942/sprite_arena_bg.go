package main

import (
	"strings"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cwin"
	"github.com/nsf/termbox-go"
)

var (
	bgWWStaticName   = "background_ww_static"
	bgWWAnimatedName = "background_ww_animated"
	// https://textkool.com/en/ascii-art-generator?hl=default&vl=default&font=Colossal&text=WW%20II
	bgWWImgTxt = strings.Trim(`
888       888 888       888      8888888 8888888
888   o   888 888   o   888        888     888
888  d8b  888 888  d8b  888        888     888
888 d888b 888 888 d888b 888        888     888
888d88888b888 888d88888b888        888     888
88888P Y88888 88888P Y88888        888     888
8888P   Y8888 8888P   Y8888        888     888
888P     Y888 888P     Y888      8888888 8888888
`, "\n")

	bg1942StaticName   = "background_1942_static"
	bg1942AnimatedName = "background_1942_animated"
	// https://textkool.com/en/ascii-art-generator?hl=default&vl=default&font=Colossal&text=1942
	bg1942ImgTxt = strings.Trim(`
 d888   .d8888b.      d8888   .d8888b.
d8888  d88P  Y88b    d8P888  d88P  Y88b
  888  888    888   d8P 888         888
  888  Y88b. d888  d8P  888       .d88P
  888   "Y888P888 d88   888   .od888P"
  888         888 8888888888 d88P"
  888  Y88b  d88P       888  888"
8888888 "Y8888P"        888  888888888`, "\n")

	bgAttr        = cwin.ChAttr{Fg: termbox.ColorDarkGray}
	bgSpeed       = cgame.ActionPerSec(100)
	bgInitialWait = 2 * time.Second
)

func newSpriteBackgroundStatic(g *cgame.Game, parent *cwin.Win, x, y int, name, imgTxt string) *cgame.SpriteBase {
	return cgame.NewSpriteBase(g, parent,
		cgame.SpriteCfg{
			Name:  name,
			Cells: cgame.StringToCells(imgTxt, bgAttr),
		},
		x, y)
}

func newSpriteBackgroundAnimated(g *cgame.Game, parent *cwin.Win, x, y int, name, imgTxt string, dx int) *cgame.SpriteAnimated {
	return cgame.NewSpriteAnimated(g, parent,
		cgame.SpriteAnimatedCfg{
			Name: name,
			Frames: [][]cgame.Cell{
				cgame.StringToCells(imgTxt, bgAttr), // single frame
			},
			DX:        dx,
			MoveSpeed: bgSpeed,
		},
		x, y)
}
