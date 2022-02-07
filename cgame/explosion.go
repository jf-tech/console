package cgame

import (
	"time"
)

type ExplosionCfg struct {
	// required
	MaxDuration time.Duration

	// optional
	AfterFinish func()
	SpriteName  string
}

func CreateExplosion(s Sprite, c ExplosionCfg) {
	s.Mgr().AddEvent(NewSpriteEventDelete(s))
	r := s.Win().Rect()
	name := s.Name() + "_explosion"
	if len(c.SpriteName) > 0 {
		name = c.SpriteName
	}
	s = NewSpriteBase(s.Mgr().g, s.Win().Parent(), name, FrameFromWin(s.Win()), r.X, r.Y)
	s.Mgr().AddEvent(NewSpriteEventCreate(s, NewAnimatorFrame(AnimatorFrameCfg{
		Frames: &explosionFrameProvider{
			s:          s,
			frameCount: int((c.MaxDuration / time.Second) * time.Duration(explosionFPS)),
		},
		AfterFinish: func(s Sprite) {
			s.Mgr().AddEvent(NewSpriteEventDelete(s))
			if c.AfterFinish != nil {
				c.AfterFinish()
			}
		},
	})))
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
	s          Sprite
	frameCount int
}

func (e *explosionFrameProvider) Next() (Frame, time.Duration, bool) {
	if e.frameCount <= 0 {
		return nil, -1, false
	}
	f := FrameFromWin(e.s.Win())
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
			if CheckProbability(explosionRuneFwdProb) {
				if idx >= len(explosionRunes)-1 {
					cellRemoval = true
				} else {
					f[i].Chx.Ch = explosionRunes[idx+1]
				}
			}
		} else if CheckProbability(changeIntoExplosionRunesProb) {
			f[i].Chx.Ch = explosionRunes[0]
		} else if CheckProbability(fireProb) {
			f[i].Chx.Ch = fire
		}
		if cellRemoval {
			for j := i; j < len(f)-1; j++ {
				f[j] = f[j+1]
			}
			f = f[:len(f)-1]
			i--
			continue
		}
	}
	e.frameCount--
	return f, time.Second / time.Duration(explosionFPS), true
}
