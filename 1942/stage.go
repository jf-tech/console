package main

import (
	"fmt"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cwin"
	"github.com/nsf/termbox-go"
)

var (
	gameOverKeys  = cgame.Keys(termbox.KeyEsc, 'q')
	pauseGameKeys = cgame.Keys('p')

	stageDurations = []time.Duration{time.Minute, time.Minute}

	betaGenProbPerStage = []*cgame.PeriodicProbabilityChecker{
		cgame.NewPeriodicProbabilityChecker("0.6%", 10*time.Millisecond),
		cgame.NewPeriodicProbabilityChecker("1%", 10*time.Millisecond),
	}
	betaFiringProbPerStage    = []string{"10%", "12%"}
	betaFiringPelletsPerStage = []int{3, 5}
)

type stage struct {
	m              *myGame
	stageIdx       int
	stageStartTime time.Duration
}

func (s *stage) Run() {
	s.init()
	s.runStageIntroBanner()
	s.m.g.Run(gameOverKeys, pauseGameKeys,
		func(ev termbox.Event) bool {
			alpha := s.m.g.SpriteMgr.FindByName(alphaName).(*spriteAlpha)
			if ev.Type == termbox.EventKey {
				if !s.m.g.IsPaused() {
					if ev.Key == termbox.KeyArrowUp {
						s.m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventSetPosRelative(alpha, 0, -1))
					} else if ev.Key == termbox.KeyArrowDown {
						s.m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventSetPosRelative(alpha, 0, 2))
					} else if ev.Key == termbox.KeyArrowLeft {
						s.m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventSetPosRelative(alpha, -3, 0))
					} else if ev.Key == termbox.KeyArrowRight {
						s.m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventSetPosRelative(alpha, 3, 0))
					} else if ev.Key == termbox.KeySpace {
						alpha.fireWeapon()
					}
				}
			}
			s.genSprites()
			s.displayStats(alpha)
			return s.checkStageDone()
		})
	if !s.m.g.IsGameOver() {
		s.runStagePassedBanner()
	}
}

func (s *stage) init() {
	bgStarGenProb.Reset(s.m.g.MasterClock)
	gpShotgunProb.Reset(s.m.g.MasterClock)
	gpShotgun2Prob.Reset(s.m.g.MasterClock)
	gpShotgun2ProbEasy.Reset(s.m.g.MasterClock)

	betaGenProbPerStage[s.stageIdx].Reset(s.m.g.MasterClock)
}

func (s *stage) runStageIntroBanner() {
	bannerDone := false
	createStageIntroBanner(s.m, s.stageIdx, func() {
		s.stageStartTime = s.m.g.MasterClock.Now()
		createAlpha(s.m)
		bannerDone = true
	})
	s.m.g.Run(gameOverKeys, pauseGameKeys, func(termbox.Event) bool { return bannerDone })
}

func (s *stage) runStagePassedBanner() {
	bannerDone := false
	createStagePassedBanner(s.m, func() {
		s.m.g.SpriteMgr.AddEvent(
			cgame.NewSpriteEventDelete(s.m.g.SpriteMgr.FindByName(alphaName)))
		bannerDone = true
	})
	s.m.g.Run(gameOverKeys, pauseGameKeys, func(termbox.Event) bool { return bannerDone })
}

func (s *stage) genSprites() {
	if s.m.g.IsPaused() {
		return
	}
	if s.checkStageWindingDown() {
		return
	}
	s.genBackgroundStar()
	s.genBeta()
	// s.genGamma()
	s.genGiftPack()
}

func (s *stage) genBackgroundStar() {
	if !bgStarGenProb.Check() {
		return
	}
	createBackgroundStar(s.m)
}

func (s *stage) genBeta() {
	if !betaGenProbPerStage[s.stageIdx].Check() {
		return
	}
	createBeta(s.m, s.stageIdx)
}

func (s *stage) genGiftPack() {
	if gpShotgunProb.Check() {
		createGiftPack(s.m, gpShotgunSym, gpShotgunSymAttr)
	}
	if s.m.easyMode {
		if gpShotgun2ProbEasy.Check() {
			createGiftPack(s.m, gpShotgun2Sym, gpShotgun2SymAttr)
		}
	} else {
		if gpShotgun2Prob.Check() {
			createGiftPack(s.m, gpShotgun2Sym, gpShotgun2SymAttr)
		}
	}
}

func (s *stage) checkStageWindingDown() bool {
	return s.m.g.MasterClock.Now()-s.stageStartTime > stageDurations[s.stageIdx]
}

func (s *stage) checkStageDone() bool {
	if !s.checkStageWindingDown() {
		return false
	}
	for _, name := range []string{
		betaName,
		betaBulletName,
		// TODO add gamma, etc.
	} {
		if _, found := s.m.g.SpriteMgr.TryFindByName(name); found {
			return false
		}
	}
	return true
}

func (s *stage) displayStats(alpha *spriteAlpha) {
	weaponName, weaponLife := alpha.weaponStats()
	killStats := alpha.killStats()
	s.m.winWeapon.SetText("WEAPON: %s (%s)", weaponName, weaponLife)
	s.m.winKills.SetText("KILLS: Beta: %s",
		func() string {
			if n, ok := killStats["Beta"]; ok {
				return fmt.Sprint(n)
			}
			return "N/A"
		}())

	s.m.winStats.SetText(fmt.Sprintf(`
Game stats:
----------------------------
Time: %s %s
%s
Internals:
----------------------------
Arena Rect: %s
FPS: %.0f
Total "pixels" rendered: %s
Memory usage: %s
%s`,
		time.Duration(s.m.g.MasterClock.Now()/(time.Second))*(time.Second),
		func() string {
			if s.m.g.IsPaused() {
				return "(paused)"
			}
			return ""
		}(),
		func() string {
			if s.m.easyMode {
				return "Easy Mode: On\n"
			}
			return ""
		}(),
		s.m.winArena.Rect(),
		s.m.g.FPS(),
		cwin.ByteSizeStr(s.m.g.WinSys.TotalChxRendered()),
		cwin.ByteSizeStr(s.m.g.HeapUsageInBytes()),
		s.m.g.SpriteMgr.DbgStats()))
}

func newStage(m *myGame, idx int) *stage {
	return &stage{m: m, stageIdx: idx}
}
