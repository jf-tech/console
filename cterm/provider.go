package cterm

type termProvider interface {
	Init() error
	Close()
	Flush() error
	Sync() error
	Size() (int, int)
	SetCell(x, y int, ch rune, fg, bg Attribute)
	PollEvent() Event
	Interrupt()
}
