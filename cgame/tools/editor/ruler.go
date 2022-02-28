package main

import (
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
)

var (
	rulerAttr = cwin.Attr{Fg: cterm.ColorDarkGray}
)

type rulers struct {
	*cwin.WinBase
}

func int2rune(x int) rune {
	if x < 0 {
		panic("invalid value")
	}
	return rune(int('0') + x%10)
}

func (r *rulers) drawRulers() {
	for x := 1; x < r.Rect().W-1; x++ {
		r.PutClientCh(x, 0, int2rune(x-1))
		r.PutClientCh(x, r.Rect().H-1, int2rune(x-1))
	}
	for y := 1; y < r.Rect().H-1; y++ {
		r.PutClientCh(0, y, int2rune(y-1))
		r.PutClientCh(r.Rect().W-1, y, int2rune(y-1))
	}
}

func (r *rulers) NonRulerRect() cwin.Rect {
	return cwin.Rect{X: 1, Y: 1, W: r.Rect().W - 2, H: r.Rect().H - 2}
}

func createRulers(sys *cwin.Sys, parent cwin.Win) *rulers {
	r := sys.CreateWinEx(parent, func() cwin.Win {
		return &rulers{
			WinBase: cwin.NewWinBase(sys, parent, cwin.WinCfg{
				R:          parent.ClientRect().ToOrigin(),
				Name:       "_rulers",
				NoBorder:   true,
				ClientAttr: rulerAttr,
				NoTitle:    true,
			}),
		}
	}).(*rulers)
	r.drawRulers()
	return r
}
