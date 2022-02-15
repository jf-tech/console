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
	win  *cwin.Win
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

type InBoundsCheckType int

const (
	InBoundsCheckNone = InBoundsCheckType(iota)
	InBoundsCheckFullyVisible
	InBoundsCheckPartiallyVisible
)

type InBoundsCheckResult int

const (
	InBoundsCheckResultOK = InBoundsCheckResult(iota)
	InBoundsCheckResultN  // breanch to the north
	InBoundsCheckResultE  // breanch to the east
	InBoundsCheckResultS  // breanch to the south
	InBoundsCheckResultW  // breanch to the west
)

type PreUpdateNotifyResponseType int

const (
	PreUpdateNotifyResponseAbandon = PreUpdateNotifyResponseType(iota)
	PreUpdateNotifyResponseJustDoIt
)

type PreUpdateNotify func(
	inBoundsCheckResult InBoundsCheckResult, collided []Sprite) PreUpdateNotifyResponseType

type UpdateArg struct {
	DXY *cwin.Point            // update the sprite position by (dx,dy), if non nil.
	F   Frame                  // If nil, no frame update; If empty (len=0), frame will be wiped clean
	IBC InBoundsCheckType      // default to no in-bounds check.
	CD  CollisionDetectionType // default to collision detection
	// This is only invoked when either the in-bounds check fails or collision check fails.
	// If the return value is PreUpdateNotifyResponseJustDoIt, the update will be carried out;
	// If PreUpdateNotifyResponseAbandon, update will be abandoned. If no Notify is given
	// for the Update call, then update will be abandoned.
	// !!! IMPORTANT !!! If the Notify implementation needs to call Update, make sure
	// UpdateArg.IBC and UpdateArg.CD are turned off, or we would get into infinite recursion!
	Notify PreUpdateNotify
}

// return true if update is carried out; false if not.
func (sb *SpriteBase) Update(arg UpdateArg) bool {
	r := sb.Rect()
	if arg.DXY != nil {
		r.X += arg.DXY.X
		r.Y += arg.DXY.Y
	}
	f := arg.F
	if f == nil {
		f = sb.Frame()
	}
	inBoundsCheckResult := InBoundsCheckResultOK
	if arg.IBC != InBoundsCheckNone {
		inBoundsCheckResult = sb.inBoundsCheck(arg.IBC, r, f)
	}
	var collided []Sprite
	if arg.CD == CollisionDetectionOn {
		collided = sb.Mgr().CheckCollision(sb, r, f)
	}
	if (inBoundsCheckResult == InBoundsCheckResultOK && len(collided) <= 0) ||
		(arg.Notify != nil &&
			arg.Notify(inBoundsCheckResult, collided) == PreUpdateNotifyResponseJustDoIt) {
		sb.win.SetPosAbs(r.X, r.Y)
		FrameToWin(f, sb.win)
		return true
	}
	return false
}

func (sb *SpriteBase) inBoundsCheck(
	checkType InBoundsCheckType, newR cwin.Rect, newF Frame) InBoundsCheckResult {
	// Note the parent rect is really parent windows' client rect shifted to origin - i.e.
	// from POV of the sprite. So we only use it's W/H components, as its X/Y are always 0.
	parentR := sb.ParentRect()
	xRegion := func(x int) int {
		if x < 0 {
			return -1
		} else if x < parentR.W {
			return 0
		}
		return 1
	}
	yRegion := func(y int) int {
		if y < 0 {
			return -1
		} else if y < parentR.H {
			return 0
		}
		return 1
	}
	totalCells := 0
	result := map[InBoundsCheckResult]int{}
	for i := 0; i < len(newF); i++ {
		if newF[i].Chx == cwin.TransparentChx() {
			continue
		}
		totalCells++
		x := newR.X + newF[i].X
		y := newR.Y + newF[i].Y
		xReg, yReg := xRegion(x), yRegion(y)
		switch xReg {
		case -1:
			switch yReg {
			case -1:
				if -x > -y {
					result[InBoundsCheckResultW]++
				} else {
					result[InBoundsCheckResultN]++
				}
			case 0:
				result[InBoundsCheckResultW]++
			case 1:
				if -x > y-parentR.H+1 {
					result[InBoundsCheckResultW]++
				} else {
					result[InBoundsCheckResultS]++
				}
			}
		case 0:
			switch yReg {
			case -1:
				result[InBoundsCheckResultN]++
			case 0:
				result[InBoundsCheckResultOK]++
			case 1:
				result[InBoundsCheckResultS]++
			}
		case 1:
			switch yReg {
			case -1:
				if x-parentR.W+1 > -y {
					result[InBoundsCheckResultE]++
				} else {
					result[InBoundsCheckResultN]++
				}
			case 0:
				result[InBoundsCheckResultE]++
			case 1:
				if x-parentR.W+1 > y-parentR.H+1 {
					result[InBoundsCheckResultE]++
				} else {
					result[InBoundsCheckResultS]++
				}
			}
		}
	}
	var maxNonOkResult InBoundsCheckResult
	maxNonOkResultCount := 0
	for k, v := range result {
		if k != InBoundsCheckResultOK && v > maxNonOkResultCount {
			maxNonOkResult, maxNonOkResultCount = k, v
		}
	}
	switch checkType {
	case InBoundsCheckFullyVisible:
		if result[InBoundsCheckResultOK] == totalCells {
			return InBoundsCheckResultOK
		}
		return maxNonOkResult
	case InBoundsCheckPartiallyVisible:
		if result[InBoundsCheckResultOK] > 0 {
			return InBoundsCheckResultOK
		}
		return maxNonOkResult
	}
	return InBoundsCheckResultOK
}

func NewSpriteBase(g *Game, parent *cwin.Win, name string, f Frame, x, y int) *SpriteBase {
	r := FrameRect(f)
	return NewSpriteBaseR(g, parent, name, f, cwin.Rect{X: x, Y: y, W: r.W, H: r.H})
}

func NewSpriteBaseR(g *Game, parent *cwin.Win, name string, f Frame, r cwin.Rect) *SpriteBase {
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
