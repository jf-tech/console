package cwin

import "github.com/jf-tech/console/cterm"

type Dir int

const (
	DirUp = Dir(iota)
	DirUpRight
	DirRight
	DirDownRight
	DirDown
	DirDownLeft
	DirLeft
	DirUpLeft
	DirCount = int(DirUpLeft)
)

const (
	DirN  = DirUp
	DirNE = DirUpRight
	DirE  = DirRight
	DirSE = DirDownRight
	DirS  = DirDown
	DirSW = DirDownLeft
	DirW  = DirLeft
	DirNW = DirUpLeft
)

var (
	DirOffSetXY = map[Dir]Point{
		DirUp:        {X: 0, Y: -1},
		DirUpRight:   {X: 1, Y: -1},
		DirRight:     {X: 1, Y: 0},
		DirDownRight: {X: 1, Y: 1},
		DirDown:      {X: 0, Y: 1},
		DirDownLeft:  {X: -1, Y: 1},
		DirLeft:      {X: -1, Y: 0},
		DirUpLeft:    {X: -1, Y: -1},
	}
)

func OffsetXYToDir(dx, dy int) Dir {
	for dir, p := range DirOffSetXY {
		if p.X == dx && p.Y == dy {
			return dir
		}
	}
	panic("unable to match any direction offset")
}

var (
	DirRunes = map[Dir]rune{
		DirN:  '↑',
		DirNE: '↗',
		DirE:  '→',
		DirSE: '↘',
		DirS:  '↓',
		DirSW: '↙',
		DirW:  '←',
		DirNW: '↖',
	}
)

var (
	Key2Dir = map[cterm.Key]Dir{
		cterm.KeyArrowUp:    DirUp,
		cterm.KeyArrowRight: DirRight,
		cterm.KeyArrowDown:  DirDown,
		cterm.KeyArrowLeft:  DirLeft,
	}
)
