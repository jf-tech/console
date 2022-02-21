package cgame

import (
	"fmt"

	"github.com/jf-tech/console/cwin"
)

type Sprite interface {
	Name() string

	// Base returns the embedded SpriteBase pointer for accessing its functionalities as well
	// as serving as an identity (by using its pointer address)
	Base() *SpriteBase
	// This returns the actual object that implements Sprite interface that is registered
	// with SpriteManager. Because SpriteBase implements Sprite interface, sometimes we go
	// into situation where a SpriteBase pointer is getting passed around but eventually
	// when deverloper wants to cast back to their own object (which embeds SpriteBase) they
	// get type assertion failure. As long as the object (that implements Sprite) passed into
	// SpriteManager.AddSprite is the "top-level" object, This() will always return that
	// registered object. Note if calling This() on a non SpriteManager managed (i.e. not
	// added with SpriteManager.AddSprite) object, it will panic.
	This() Sprite

	Mgr() *SpriteManager
	Game() *Game

	Rect() cwin.Rect
	ParentRect() cwin.Rect
	Frame() Frame
	Animators() []Animator

	Destroy()

	fmt.Stringer
}
