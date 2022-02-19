package cgame

import (
	"github.com/jf-tech/console/cwin"
)

type Sprite interface {
	Name() string
	UID() int64
	Mgr() *SpriteManager
	Game() *Game
	Rect() cwin.Rect
	ParentRect() cwin.Rect
	Frame() Frame
	Animators() []Animator
	Destroy()
}

type SpriteBase struct {
	name string
	uid  int64
	mgr  *SpriteManager
	g    *Game
	win  cwin.Win
	as   []Animator
}

func (sb *SpriteBase) Name() string {
	return sb.name
}

func (sb *SpriteBase) UID() int64 {
	return sb.uid
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
	sb.g.WinSys.RemoveWin(sb.win)
	sb.win = nil
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

func (sb *SpriteBase) ToBottom() {
	sb.win.SendToBottom(false)
}

func (sb *SpriteBase) ToTop() {
	sb.win.SendToTop(false)
}

type UpdateArg struct {
	DXY *cwin.Point            // update the sprite position by (dx,dy), if non nil.
	F   Frame                  // If nil, no frame update; If empty (len=0), frame will be wiped clean
	IBC InBoundsCheckType      // default to no in-bounds check.
	CD  CollisionDetectionType // default to collision detection
}

// returns true if update is carried out; false if not.
// Important, do not call Update with IBC/CD turned on from your InBoundsCheckResponse.Notify
// or CollisionResponse.Notify or it might cause infinite recursion.
func (sb *SpriteBase) Update(arg UpdateArg) bool {
	s := sb.toSprite()
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
		if r, ok := s.(InBoundsCheckResponse); ok {
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
		if r, ok := s.(CollisionResponse); ok {
			resp = r.CollisionNotify(true, collided)
		}
		if resp == CollisionResponseAbandon {
			return false
		}
		for _, c := range collided {
			if r, ok := c.(CollisionResponse); ok {
				r.CollisionNotify(false, []Sprite{s})
			}
		}
	}
	sb.win.SetPosRel(r.X-s.Rect().X, r.Y-s.Rect().Y)
	FrameToWin(f, sb.win)
	return true
}

func (sb *SpriteBase) toSprite() Sprite {
	// Most the time, game sprite uses *SpriteBase as an embedded field (so to have all
	// Sprite interface functionalities plus extra for free). But we sometimes need to
	// get the actual sprite object so we can interface type assertion. Since the sprites
	// are typically added to the SpriteManager using its object, not SpriteBase, thus we
	// can query that Sprite interface from SpriteManager and then we can do type assertion.
	if s, ok := sb.Mgr().TryFindByUID(sb.UID()); ok {
		return s
	}
	return sb
}

func NewSpriteBase(g *Game, parent cwin.Win, name string, f Frame, x, y int) *SpriteBase {
	r := FrameRect(f)
	return NewSpriteBaseR(g, parent, name, f, cwin.Rect{X: x, Y: y, W: r.W, H: r.H})
}

func NewSpriteBaseR(g *Game, parent cwin.Win, name string, f Frame, r cwin.Rect) *SpriteBase {
	sb := &SpriteBase{
		name: name,
		uid:  cwin.GenUID(),
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
