package cgame

import "github.com/jf-tech/console/cwin"

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
	dirCount
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
	DirOffSetXY = [dirCount]cwin.Point{
		{X: 0, Y: -1},  // up
		{X: 1, Y: -1},  // up right
		{X: 1, Y: 0},   // right
		{X: 1, Y: 1},   // down right
		{X: 0, Y: 1},   // down
		{X: -1, Y: 1},  // down left
		{X: -1, Y: 0},  // left
		{X: -1, Y: -1}, // up left
	}
	DirSymbols = [dirCount]rune{'↑', '↗', '→', '↘', '↓', '↙', '←', '↖'}
)
