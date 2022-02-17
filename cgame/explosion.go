package cgame

import (
	"time"

	"github.com/jf-tech/console/cutil"
)

type ExplosionCfg struct {
	// required
	MaxDuration time.Duration

	// optional
	AfterFinish func()
	SpriteName  string
}

func CreateExplosion(s *SpriteBase, c ExplosionCfg) {
	r := s.win.Rect()
	name := s.Name() + "_explosion"
	if len(c.SpriteName) > 0 {
		name = c.SpriteName
	}
	newS := NewSpriteBase(s.Mgr().g, s.win.Parent(), name, FrameFromWin(s.win), r.X, r.Y)
	newS.AddAnimator(NewAnimatorFrame(newS, AnimatorFrameCfg{
		Frames: &explosionFrameProvider{
			s:          newS,
			frameCount: int((c.MaxDuration / time.Second) * time.Duration(explosionFPS)),
		},
		AnimatorCfgCommon: AnimatorCfgCommon{
			AfterFinish: c.AfterFinish,
		},
	}))
	newS.Mgr().AsyncCreateSprite(newS)
	s.Mgr().AsyncDeleteSprite(s)
}

var (
	explosionRunes               = []rune("\"~'`.")
	fire                         = '🔥'
	fireProb                     = "33%"
	explosionFPS                 = 8
	changeIntoExplosionRunesProb = "50%"
	explosionRuneFwdProb         = "80%"
)

type explosionFrameProvider struct {
	s          *SpriteBase
	frameCount int
}

func (e *explosionFrameProvider) Next() (Frame, time.Duration, bool) {
	if e.frameCount <= 0 {
		return nil, -1, false
	}
	f := FrameFromWin(e.s.win)
	if len(f) <= 0 {
		return nil, -1, false
	}
	indexRune := func(rs []rune, r rune) int {
		for i := 0; i < len(rs); i++ {
			if rs[i] == r {
				return i
			}
		}
		return -1
	}
	for i := 0; i < len(f); i++ {
		cellRemoval := false
		if idx := indexRune(explosionRunes, f[i].Chx.Ch); idx >= 0 {
			if cutil.CheckProbability(explosionRuneFwdProb) {
				if idx >= len(explosionRunes)-1 {
					cellRemoval = true
				} else {
					f[i].Chx.Ch = explosionRunes[idx+1]
				}
			}
		} else if cutil.CheckProbability(changeIntoExplosionRunesProb) {
			f[i].Chx.Ch = explosionRunes[0]
		} else if cutil.CheckProbability(fireProb) {
			f[i].Chx.Ch = fire
		}
		if cellRemoval {
			copy(f[i:], f[i+1:])
			f = f[:len(f)-1]
			i--
			continue
		}
	}
	e.frameCount--
	return f, time.Second / time.Duration(explosionFPS), true
}
