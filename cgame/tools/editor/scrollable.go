package main

import "github.com/jf-tech/console/cwin"

type ScrollableRenderer func(w cwin.Win, cfg ScrollableCfg, cur int)

type ScrollableCfg struct {
	cwin.WinCfg
	TotalRange   int
	VisibleRange int
	Renderer     ScrollableRenderer
}

type Scrollable struct {
	*cwin.WinBase
	cfg ScrollableCfg
	cur int
}

func (s *Scrollable) SetCur(cur int) {
	if cur < 0 || cur >= s.cfg.TotalRange {
		panic("invalid cur value")
	}
	s.cur = cur
	s.cfg.Renderer(s.This(), s.cfg, s.cur)
}

func (s *Scrollable) GetCur() int {
	return s.cur
}

func createScrollable(sys *cwin.Sys, parent cwin.Win, cfg ScrollableCfg) *Scrollable {
	s := sys.CreateWinEx(parent, func() cwin.Win {
		cfg.WinCfg.NoBorder = true
		cfg.WinCfg.NoTitle = true
		return &Scrollable{WinBase: cwin.NewWinBase(sys, parent, cfg.WinCfg), cfg: cfg}
	}).(*Scrollable)
	s.SetCur(0)
	return s
}

func CreateHScrollable(sys *cwin.Sys, parent cwin.Win, cfg ScrollableCfg) *Scrollable {
	cfg.R.H = 1
	if cfg.R.W%cfg.VisibleRange != 0 {
		panic("W must be a multiple of VisibleRange")
	}
	return createScrollable(sys, parent, cfg)
}

func CreateVScrollable(sys *cwin.Sys, parent cwin.Win, cfg ScrollableCfg) *Scrollable {
	cfg.R.W = 1
	if cfg.R.H%cfg.VisibleRange != 0 {
		panic("H must be a multiple of VisibleRange")
	}
	return createScrollable(sys, parent, cfg)
}

func ScrollableRendererDigits(w cwin.Win, cfg ScrollableCfg, cur int) {
}
