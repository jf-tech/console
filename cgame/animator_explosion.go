package cgame

import (
	"time"

	"github.com/jf-tech/console/cwin"
	"github.com/jf-tech/go-corelib/maths"
)

var (
	explosionRunes = []rune(`|/-\-~. `)
)

type animatorExplosion struct {
	region      cwin.Rect
	t           time.Duration
	afterFinish func(Sprite)

	clock     *Clock
	start     time.Duration
	layers    int // number of layers of concentric rects of the region.
	expansion int // number of cells of expansion to reach the edge of the sprite rect.
}

func (ase *animatorExplosion) Animate(s Sprite) AnimatorState {
	ase.checkToInit(s)
	pct := maths.MaxInt(int((ase.clock.Now()-ase.start)*100/ase.t), 100)
	if pct == 100 {
		s.Mgr().AddEvent(NewSpriteEventDelete(s))
		if ase.afterFinish != nil {
			ase.afterFinish(s)
		}
		return AnimatorCompleted
	}
	return AnimatorRunning
}

func (ase *animatorExplosion) checkToInit(s Sprite) {
	if ase.clock != nil {
		return
	}
	ase.clock = s.Game().MasterClock
	ase.start = ase.clock.Now()
	ase.layers = maths.MinInt(ase.region.W+1, ase.region.H+1) / 2
	sr := s.Win().Rect()
	ase.expansion = maths.MinInt((sr.W-ase.region.W)/2, (sr.H-ase.region.H)/2)
}

type ExplosionCfg struct {
	Scale       float64
	T           time.Duration
	AfterFinish func(Sprite)
}

func CreateExplosion(mgr *SpriteManager, s Sprite, c ExplosionCfg) {
	if c.Scale < 1 {
		panic("explosion scale must be >= 1.0, no implosion please")
	}
	mgr.AddEvent(NewSpriteEventDelete(s))
	f := FrameFromWin(s.Win())
	r := s.Win().Rect()
	newW := maths.MaxInt(r.W, int(float64(r.W)*c.Scale))
	newH := maths.MaxInt(r.H, int(float64(r.H)*c.Scale))
	offsetX, offsetY := (newW-r.W)/2, (newH-r.H)/2
	for i := 0; i < len(f); i++ {
		f[i].X += offsetX
		f[i].Y += offsetY
	}
	s = NewSpriteBaseR(mgr.g, s.Win().Parent(),
		s.Name()+"_explosion", f,
		cwin.Rect{X: r.X - offsetX, Y: r.Y - offsetY, W: newW, H: newH})
	mgr.AddEvent(NewSpriteEventCreate(s, &animatorExplosion{
		region:      cwin.Rect{X: offsetX, Y: offsetY, W: r.W, H: r.H},
		t:           c.T,
		afterFinish: c.AfterFinish,
	}))
}
