package cgame

import (
	"fmt"
	"unsafe"

	"github.com/jf-tech/console/cwin"
)

type SpriteBase struct {
	name string
	mgr  *SpriteManager
	g    *Game
	win  cwin.Win
	as   []Animator
}

// Sprite interface methods

func (sb *SpriteBase) Name() string {
	return sb.name
}

func (sb *SpriteBase) Base() *SpriteBase {
	return sb
}

func (sb *SpriteBase) This() Sprite {
	return sb.Mgr().Find(sb)
}

func (sb *SpriteBase) Mgr() *SpriteManager {
	return sb.mgr
}

func (sb *SpriteBase) Game() *Game {
	return sb.g
}

func (sb *SpriteBase) Rect() cwin.Rect {
	return sb.win.Rect()
}

func (sb *SpriteBase) ParentRect() cwin.Rect {
	return sb.win.Parent().ClientRect().ToOrigin()
}

func (sb *SpriteBase) Frame() Frame {
	return FrameFromWin(sb.win)
}

func (sb *SpriteBase) Animators() []Animator {
	// make a snapshot to return so that caller won't get into potentially
	// changing slice.
	var cp []Animator
	cp = append(cp, sb.as...)
	return cp
}

func (sb *SpriteBase) Destroy() {
	if !sb.IsDestroyed() {
		sb.Game().WinSys.RemoveWin(sb.win)
	}
	sb.win = nil
}

func (sb *SpriteBase) String() string {
	return fmt.Sprintf("SpriteBase['%s'|0x%X]", sb.Name(), uintptr(unsafe.Pointer(sb)))
}

// Non Sprite interface methods

func (sb *SpriteBase) IsDestroyed() bool {
	return sb.win == nil
}

func (sb *SpriteBase) AddAnimator(as ...Animator) {
	sb.as = append(sb.as, as...)
}

func (sb *SpriteBase) DeleteAnimator(as ...Animator) {
	for _, a := range as {
		for i := 0; i < len(sb.as); i++ {
			if sb.as[i] == a {
				copy(sb.as[i:], sb.as[i+1:])
				sb.as = sb.as[:len(sb.as)-1]
				break
			}
		}
	}
}

func (sb *SpriteBase) SendToBottom() {
	sb.win.SendToBottom(false)
}

func (sb *SpriteBase) SendToTop() {
	sb.win.SendToTop(false)
}

type UpdateArg struct {
	DXY      *cwin.Point            // update the sprite position by (dx,dy), if non nil.
	F        Frame                  // If nil, no frame update; If empty (len=0), frame will be wiped clean
	IBC      InBoundsCheckType      // default to no in-bounds check.
	CD       CollisionDetectionType // default to collision detection
	TestOnly bool                   // if true, do the tests only, no actual update or notifications
}

// returns false if any of the tests (bounds, collision) fails; true otherwise.
// Important, do not call Update with IBC/CD turned on from your InBoundsCheckResponse.Notify
// or CollisionResponse.Notify or it might cause infinite recursion.
func (sb *SpriteBase) Update(arg UpdateArg) bool {
	// If bounds check or collision check results in notifying caller, we want to try our best
	// to notify them with the "top-level" sprite implementer objects. However this is only
	// possible if the sprite is already registered with the SpriteManager. If an Update call
	// comes before the sprite is AddSprite into SpriteManager (which is sometimes necessary: i.e.
	// caller wants to do bounds or collision check before fully creating/registering the sprite)
	// we have to fallback to the SpriteBase as the notification arguments.
	var s Sprite
	var ok bool
	if s, ok = sb.Mgr().TryFind(sb); !ok {
		s = sb
	}
	r := s.Rect()
	if arg.DXY != nil {
		r.X += arg.DXY.X
		r.Y += arg.DXY.Y
	}
	f := arg.F
	if f == nil {
		f = s.Frame()
	}
	inBoundsCheckResult := InBoundsCheckResultOK
	if arg.IBC != InBoundsCheckNone {
		inBoundsCheckResult = InBoundsCheck(arg.IBC, r, f, s.ParentRect())
	}
	if inBoundsCheckResult != InBoundsCheckResultOK {
		resp := InBoundsCheckResponseAbandon
		if r, ok := s.(InBoundsCheckResponse); ok && !arg.TestOnly {
			resp = r.InBoundsCheckNotify(inBoundsCheckResult)
		}
		if resp == InBoundsCheckResponseAbandon {
			return false
		}
	}
	var collided []Sprite
	if arg.CD == CollisionDetectionOn {
		collided = s.Mgr().CheckCollision(s, r, f)
	}
	if len(collided) > 0 {
		resp := CollisionResponseAbandon
		if r, ok := s.(CollisionResponse); ok && !arg.TestOnly {
			resp = r.CollisionNotify(true, collided)
		}
		if resp == CollisionResponseAbandon {
			return false
		}
		for _, c := range collided {
			if r, ok := c.(CollisionResponse); ok && !arg.TestOnly {
				r.CollisionNotify(false, []Sprite{s})
			}
		}
	}
	// It is possible that during the bounds check or collision detection notification
	// callback, the sprite decides to SpriteManager.DeleteSprite itself (e.g. a bullet
	// sprite hits a target, and it destroys itself during collision callback). Thus we
	// have to guard the update against that.
	if !arg.TestOnly && !sb.IsDestroyed() {
		sb.win.SetPosRel(r.X-s.Rect().X, r.Y-s.Rect().Y)
		FrameToWin(f, sb.win)
	}
	return true
}

func NewSpriteBase(g *Game, parent cwin.Win, name string, f Frame, x, y int) *SpriteBase {
	r := FrameRect(f)
	return NewSpriteBaseR(g, parent, name, f, cwin.Rect{X: x, Y: y, W: r.W, H: r.H})
}

func NewSpriteBaseR(g *Game, parent cwin.Win, name string, f Frame, r cwin.Rect) *SpriteBase {
	sb := &SpriteBase{
		name: name,
		mgr:  g.SpriteMgr,
		g:    g,
	}
	winCfg := cwin.WinCfg{
		R:        r,
		Name:     name,
		NoBorder: true,
	}
	sb.win = g.WinSys.CreateWin(parent, winCfg)
	sb.Update(UpdateArg{
		F:   f,
		IBC: InBoundsCheckNone,
		CD:  CollisionDetectionOff,
	})
	return sb
}
