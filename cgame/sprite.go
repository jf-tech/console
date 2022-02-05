package cgame

import (
	"github.com/jf-tech/console/cwin"
)

type Sprite interface {
	Name() string
	UID() int64
	Win() *cwin.Win
	Mgr() *SpriteManager
	Game() *Game
}

type SpriteBase struct {
	name     string
	uid      int64
	win      *cwin.Win
	mgr      *SpriteManager
	g        *Game
	curFrame SpriteFrame
}

func (sb *SpriteBase) Name() string {
	return sb.name
}

func (sb *SpriteBase) UID() int64 {
	return sb.uid
}

func (sb *SpriteBase) Win() *cwin.Win {
	return sb.win
}

func (sb *SpriteBase) Mgr() *SpriteManager {
	return sb.mgr
}

func (sb *SpriteBase) Game() *Game {
	return sb.g
}

func (sb *SpriteBase) CurFrame() SpriteFrame {
	return sb.curFrame
}

func (sb *SpriteBase) SetFrame(f SpriteFrame) {
	FrameToWin(f, sb.win)
}

func NewSpriteBase(g *Game, parent *cwin.Win, name string, f SpriteFrame, x, y int) *SpriteBase {
	sb := &SpriteBase{
		name:     name,
		uid:      cwin.GenUID(),
		mgr:      g.SpriteMgr,
		g:        g,
		curFrame: f,
	}
	r := FrameRect(f)
	winCfg := cwin.WinCfg{
		R:        cwin.Rect{X: x, Y: y, W: r.W, H: r.H},
		Name:     name,
		NoBorder: true,
	}
	sb.win = g.WinSys.CreateWin(parent, winCfg)
	sb.SetFrame(f)
	return sb
}
