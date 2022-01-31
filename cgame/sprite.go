package cgame

import (
	"math"
	"strings"

	"github.com/jf-tech/console/cwin"
	"github.com/jf-tech/go-corelib/maths"
)

type Sprite interface {
	Cfg() SpriteCfg
	UID() int64
	Win() *cwin.Win
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

// SpriteBase represents a basic sprite that
//  - possess a region (rectangle)
//  - always visible
//  - not PositionSettable
// It forms a basis for all user defined sprites that can have additional
// capabilities such as PositionSettable, Collidable etc.
type SpriteBase struct {
	Config   SpriteCfg
	UniqueID int64
	Mgr      *SpriteManager
	W        *cwin.Win
}

func (sb *SpriteBase) Cfg() SpriteCfg {
	return sb.Config
}

func (sb *SpriteBase) UID() int64 {
	return sb.UniqueID
}

func (sb *SpriteBase) Win() *cwin.Win {
	return sb.W
}

func NewSpriteBase(g *Game, parent *cwin.Win, c SpriteCfg, x, y int) *SpriteBase {
	sb := &SpriteBase{
		Config:   c,
		UniqueID: cwin.GenUID(),
		Mgr:      g.SpriteMgr,
	}
	if len(sb.Config.Cells) <= 0 {
		panic("cells cannot be empty")
	}
	w, h := normalizeCells(sb.Config.Cells)
	winCfg := cwin.WinCfg{
		R:        cwin.Rect{X: x, Y: y, W: w, H: h},
		Name:     sb.Config.Name,
		NoBorder: true,
	}
	sb.W = g.WinSys.CreateWin(parent, winCfg)
	putNormalizedCellsToWin(sb.Config.Cells, sb.W)
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
