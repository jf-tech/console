package cgame

import (
	"math"
	"strings"

	"github.com/jf-tech/console/cwin"
	"github.com/jf-tech/go-corelib/maths"
)

// - Merge Sprite and PositionSettable - SpriteBase always support setting position
// - SpriteMgr directly set position on Sprite.
// - Add pre-/post event hook for SpriteEventSetPosRelative in SpriteCfg
// - Add built-in Pre set pos event hook to prevent out of bound
// - Add built-in post set pos event hook to do auto delete.
// - Rename Animated to Dynamic; Act to Update -- scratch that. shouldn't we just remove Animated interface?
// - Instead, add utility (animator) to help SpriteMgr to animate sprites via the Sprite.SetPosRelative()
//    - LinearAnimator(Sprite, init x/y, target x/y, speed)
//    - FrameAnimator
// - cgame.Run(func())
// - cgame.FPS(), Mem() etc

type Sprite interface {
	Name() string
	UID() int64
	Win() *cwin.Win
	Clock() *Clock
	Mgr() *SpriteManager
	Game() *Game
}

type PositionSettable interface {
	SetPosRelative(dx, dy int)
}

type Animated interface {
	Act()
}

type Collidable interface {
	Collided(other Sprite)
}

type Cell struct {
	X, Y int
	Chx  cwin.Chx
}

func StringToCells(s string, attr cwin.ChAttr) []Cell {
	s = strings.Trim(s, "\n")
	rect := cwin.TextDimension(s)
	cells := make([]Cell, rect.W*rect.H)
	for y := 0; y < rect.H; y++ {
		for x := 0; x < rect.W; x++ {
			idx := y*rect.W + x
			cells[idx].X = x
			cells[idx].Y = y
			cells[idx].Chx = cwin.TransparentChx()
		}
	}
	rs := []rune(s)
	rsLen := len(rs)
	for i, x, y := 0, 0, 0; i < rsLen; i++ {
		if rs[i] == '\n' {
			x = 0
			y++
			continue
		}
		if rs[i] != ' ' {
			idx := y*rect.W + x
			cells[idx].Chx = cwin.Chx{Ch: rs[i], Attr: attr}
		}
		x++
	}
	return cells
}

type SpriteCfg struct {
	Name  string
	Cells []Cell
}

type SpriteBase struct {
	cfg   SpriteCfg
	uid   int64
	win   *cwin.Win
	clock *Clock
	mgr   *SpriteManager
	g     *Game
}

func (sb *SpriteBase) Name() string {
	return sb.cfg.Name
}

func (sb *SpriteBase) UID() int64 {
	return sb.uid
}

func (sb *SpriteBase) Win() *cwin.Win {
	return sb.win
}

func (sb *SpriteBase) Clock() *Clock {
	return sb.clock
}

func (sb *SpriteBase) Mgr() *SpriteManager {
	return sb.mgr
}

func (sb *SpriteBase) Game() *Game {
	return sb.g
}

func NewSpriteBase(g *Game, parent *cwin.Win, cfg SpriteCfg, x, y int) *SpriteBase {
	if len(cfg.Cells) <= 0 {
		panic("cells cannot be empty")
	}
	sb := &SpriteBase{
		cfg:   cfg,
		uid:   cwin.GenUID(),
		clock: g.SpriteMgr.g.SpriteMgr.clockMgr.createClock(),
		mgr:   g.SpriteMgr,
		g:     g,
	}
	w, h := normalizeCells(sb.cfg.Cells)
	winCfg := cwin.WinCfg{
		R:        cwin.Rect{X: x, Y: y, W: w, H: h},
		Name:     sb.cfg.Name,
		NoBorder: true,
	}
	sb.win = g.WinSys.CreateWin(parent, winCfg)
	putNormalizedCellsToWin(sb.cfg.Cells, sb.win)
	return sb
}

func normalizeCells(cells []Cell) (w, h int) {
	minX, maxX := math.MaxInt32, math.MinInt32
	minY, maxY := math.MaxInt32, math.MinInt32
	for i := 0; i < len(cells); i++ {
		minX = maths.MinInt(minX, cells[i].X)
		maxX = maths.MaxInt(maxX, cells[i].X)
		minY = maths.MinInt(minY, cells[i].Y)
		maxY = maths.MaxInt(maxY, cells[i].Y)
	}
	for i := 0; i < len(cells); i++ {
		cells[i].X -= minX
		cells[i].Y -= minY
	}
	return maxX - minX + 1, maxY - minY + 1
}

func putNormalizedCellsToWin(normalizedCells []Cell, w *cwin.Win) {
	w.FillClient(w.ClientRect().ToOrigin(), cwin.TransparentChx())
	for i := 0; i < len(normalizedCells); i++ {
		w.PutClient(normalizedCells[i].X, normalizedCells[i].Y, normalizedCells[i].Chx)
	}
}
