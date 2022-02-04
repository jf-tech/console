package cgame

import (
	"fmt"

	"github.com/jf-tech/console/cwin"
)

func StringsToFrames(ss []string, attr cwin.ChAttr) [][]Cell {
	var ret [][]Cell
	for _, s := range ss {
		ret = append(ret, StringToCells(s, attr))
	}
	return ret
}

type SpriteAnimatedCfg struct {
	Name       string
	Frames     [][]Cell
	FrameSpeed ActionPerSec
	Loop       bool
	DX, DY     int
	MoveSpeed  ActionPerSec
	AfterMove  func(Sprite)
}

type SpriteAnimated struct {
	*SpriteBase
	Config          SpriteAnimatedCfg
	AnimationTicker *ActionPerSecTicker
	MoveTicker      *ActionPerSecTicker

	curFrame int
}

func (sa *SpriteAnimated) Act() {
	// frame animation
	sa.setFrame(sa.curFrame + int(sa.AnimationTicker.HowMany()))

	// movement
	moves := int(sa.MoveTicker.HowMany())
	if moves <= 0 {
		return
	}
	dx, dy := sa.Config.DX*moves, sa.Config.DY*moves
	newR := sa.win.Rect().MoveDelta(dx, dy)
	if overlapped, _ := sa.win.Parent().ClientRect().Overlap(newR); !overlapped {
		sa.mgr.AddEvent(NewSpriteEventDelete(sa))
		return
	}
	sa.win.SetPosRelative(dx, dy)
	if sa.Config.AfterMove != nil {
		sa.Config.AfterMove(sa)
	}
}

func (sa *SpriteAnimated) setFrame(frame int) {
	if frame == sa.curFrame {
		return
	}
	if frame >= len(sa.Config.Frames) && !sa.Config.Loop {
		sa.mgr.AddEvent(NewSpriteEventDelete(sa))
		return
	}
	sa.curFrame = frame % len(sa.Config.Frames)
	putNormalizedCellsToWin(sa.Config.Frames[sa.curFrame], sa.win)
}

func NewSpriteAnimated(g *Game, parent *cwin.Win, c SpriteAnimatedCfg, x, y int) *SpriteAnimated {
	normalizeFrames(c.Frames)
	baseCfg := &SpriteCfg{
		Name:  c.Name,
		Cells: c.Frames[0],
	}
	sb := NewSpriteBase(g, parent, *baseCfg, x, y)
	sa := &SpriteAnimated{
		sb,
		c,
		NewActionPerSecTicker(sb.clock, c.FrameSpeed, true),
		NewActionPerSecTicker(sb.clock, c.MoveSpeed, true),
		0,
	}
	return sa
}

func normalizeFrames(frames [][]Cell) {
	if len(frames) <= 0 {
		panic("frames cannot be empty")
	}
	w, h := -1, -1
	for i := 0; i < len(frames); i++ {
		nw, nh := normalizeCells(frames[i])
		if w == -1 {
			w, h = nw, nh
		}
		if nw != w || nh != h {
			panic(fmt.Sprintf("frame[%d] has a different WxH(%dx%d) than frame[0] WxH(%dx%d)",
				i, nw, nh, w, h))
		}
	}
}
