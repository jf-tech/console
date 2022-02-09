package cgame

var (
	DirOffSetXY = []PairInt{
		{A: 0, B: -1},  // up
		{A: 1, B: -1},  // up right
		{A: 1, B: 0},   // right
		{A: 1, B: 1},   // down right
		{A: 0, B: 1},   // down
		{A: -1, B: 1},  // down left
		{A: -1, B: 0},  // left
		{A: -1, B: -1}, // up left
	}
	DirSymbols = []rune("↑↗→↘↓↙←↖")
)
