package cterm

type Provider int

const (
	TermBox Provider = iota
	TCell
)

var (
	provider = termProvider(&providerTermBox{})
)

func SetProvider(p Provider) {
	switch p {
	case TermBox:
		provider = &providerTermBox{}
	case TCell:
		provider = &providerTCell{}
	}
}

func Init() error {
	return provider.Init()
}

func Close() {
	provider.Close()
}

func Flush() error {
	return provider.Flush()
}

func Size() (int, int) {
	return provider.Size()
}

// Attribute affects the presentation of characters, such as color, boldness,
// and so forth.
type Attribute uint64

// Colors first.  The order here is significant.
const (
	ColorDefault Attribute = iota
	ColorBlack
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
	ColorDarkGray
	ColorLightRed
	ColorLightGreen
	ColorLightYellow
	ColorLightBlue
	ColorLightMagenta
	ColorLightCyan
	ColorLightGray
)

// SetCell sets the character cell at a given location to the given
// content (rune) and attributes.
func SetCell(x, y int, ch rune, fg, bg Attribute) {
	provider.SetCell(x, y, ch, fg, bg)
}

// EventType represents the type of event.
type EventType uint8

// Modifier represents the possible modifier keys.
type Modifier int16

// Key is a key press.
type Key int16

// Event represents an event like a key press, mouse action, or window resize.
type Event struct {
	Type EventType
	Key  Key
	Ch   rune
}

// Event types.
const (
	EventNone EventType = iota
	EventKey
)

// Keys codes.
const (
	KeyArrowUp = Key(iota + 1000)
	KeyArrowDown
	KeyArrowRight
	KeyArrowLeft
	KeyEnter
	KeyEsc
)

// PollEvent blocks until an event is ready, and then returns it.
func PollEvent() Event {
	return provider.PollEvent()
}

// Interrupt posts an interrupt event.
func Interrupt() {
	provider.Interrupt()
}
