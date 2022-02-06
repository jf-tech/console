package cgame

import (
	"time"

	"github.com/jf-tech/console/cwin"
)

type AnimatorExplosionCfg struct {
	Region cwin.Rect
	T      time.Duration
}

type AnimatorExplosion struct {
	cfg AnimatorExplosionCfg

	clock *Clock
}

func (ase *AnimatorExplosion) Animate(s Sprite) AnimatorState {
	ase.checkToInit(s)
	return AnimatorCompleted
}

func (ase *AnimatorExplosion) checkToInit(s Sprite) {
	if ase.clock != nil {
		return
	}
	ase.clock = s.Game().MasterClock
}

func NewAnimatorExplosion(c AnimatorExplosionCfg) *AnimatorExplosion {
	return &AnimatorExplosion{cfg: c}
}
