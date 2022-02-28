package cwin

const (
	RuneSpace rune = ' '
)

const (
	BorderRuneUL int = iota
	BorderRuneUR
	BorderRuneLR
	BorderRuneLL
	BorderRuneV
	BorderRuneH
	BorderRuneCount
)

type BorderRunes [BorderRuneCount]rune

var (
	SingleLineBorderRunes = BorderRunes{'┏', '┓', '┛', '┗', '┃', '━'}
	DoubleLineBorderRunes = BorderRunes{'╔', '╗', '╝', '╚', '║', '═'}
)

// these vars not intended to be used, just here for library purpose.
var (
	boxDrawingRunes1 = `
┏━━┓
┃  ┃
┗━━┛`
	boxDrawingRunes2 = `
⎡‾‾⎤
⎢  ⎥
⎣__⎦`
	boxDrawingRunes3 = `
┌──┐
│  │
└──┘`
	boxDrawingRunes4 = `
╔══╗
║  ║
╚══╝`
	boxDrawingRunes5 = `
╭──╮
│  │
╰──╯`
)

var (
	blocks = `░`
)
