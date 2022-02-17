package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
)

type stage struct {
	m        *myGame
	stageIdx int

	stageStartTime time.Duration
	stageSkipped   bool
	bossCreated    bool
}

func (s *stage) Run() {
	s.init()
	s.runStageIntroBanner()
	s.m.g.Run(gameOverKeys, pauseGameKeys,
		func(ev cterm.Event) bool {
			if s.checkStageDone() {
				return true
			}
			alpha := s.m.g.SpriteMgr.FindByName(alphaName).(*spriteAlpha)
			alpha.ToTop()
			if ev.Type == cterm.EventKey {
				// due to console aspect ration, make left/right move a bit faster.
				// also let retreat (down) a bit faster than up to make the game exp
				// better.
				if ev.Key == cterm.KeyArrowUp {
					alpha.move(0, -1)
				} else if ev.Key == cterm.KeyArrowDown {
					alpha.move(0, 2)
				} else if ev.Key == cterm.KeyArrowLeft {
					alpha.move(-3, 0)
				} else if ev.Key == cterm.KeyArrowRight {
					alpha.move(3, 0)
				} else if ev.Ch == ' ' {
					alpha.fireWeapon()
				} else if cwin.FindKey(skipStageKeys, ev) {
					s.stageSkipped = true
				} else if cwin.FindKey(invincibleModeKeys, ev) {
					s.m.invincible = !s.m.invincible
				}
			}
			s.genSprites()
			s.displayStats(alpha)
			return false
		})
	if !s.m.g.IsGameOver() && s.stageIdx != totalStages-1 {
		s.runStagePassedBanner()
	}
}

func (s *stage) init() {
	bgStarGenProb.Reset(s.m.g.MasterClock)
	gpShotgunProb.Reset(s.m.g.MasterClock)
	gpShotgun2Prob.Reset(s.m.g.MasterClock)
	gpShotgun2ProbEasy.Reset(s.m.g.MasterClock)

	betaGenProbPerStage[s.stageIdx].Reset(s.m.g.MasterClock)
	gammaGenProbPerStage[s.stageIdx].Reset(s.m.g.MasterClock)
	deltaGenProbPerStage[s.stageIdx].Reset(s.m.g.MasterClock)
}

func (s *stage) runStageIntroBanner() {
	bannerDone := false
	createStageIntroBanner(s.m, s.stageIdx, func() {
		s.stageStartTime = s.m.g.MasterClock.Now()
		createAlpha(s.m, s)
		bannerDone = true
	})
	s.m.g.Run(gameOverKeys, pauseGameKeys, func(cterm.Event) bool { return bannerDone })
}

func (s *stage) runStagePassedBanner() {
	bannerDone := false
	createStagePassedBanner(s.m, func() {
		bannerDone = true
	})
	s.m.g.Run(gameOverKeys, pauseGameKeys, func(cterm.Event) bool { return bannerDone })
}

func (s *stage) genSprites() {
	if s.checkStageWindingDown() {
		if s.stageIdx == totalStages-1 && !s.bossCreated {
			s.genBoss()
			s.bossCreated = true
		}
		return
	}
	s.genBackgroundStar()
	s.genBeta()
	s.genGamma()
	s.genDelta()
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

func (s *stage) genGamma() {
	if !gammaGenProbPerStage[s.stageIdx].Check() {
		return
	}
	createGamma(s.m, s.stageIdx)
}

func (s *stage) genDelta() {
	if !deltaGenProbPerStage[s.stageIdx].Check() {
		return
	}
	createDelta(s.m)
}

func (s *stage) genBoss() {
	createBoss(s.m)
	s.m.g.SoundMgr.PlayMP3(sfxBossIsHereFile, sfxClipVol, 1)
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
	return s.m.g.MasterClock.Now()-s.stageStartTime > stageDurations[s.stageIdx] || s.stageSkipped
}

func (s *stage) checkStageDone() bool {
	if !s.checkStageWindingDown() {
		return false
	}
	if s.stageIdx == totalStages-1 && !s.bossCreated {
		return false
	}
	// waiting for all the enemy  sprites to be done/out 'coz they might still kill
	// our player during the process :)
	for _, name := range []string{
		betaName,
		gammaName,
		deltaName,
		bossName,
		bossExplosionName,
	} {
		if _, found := s.m.g.SpriteMgr.TryFindByName(name); found {
			return false
		}
	}
	// we're truly done. remove all the non enemy sprites
	s.m.g.SpriteMgr.AsyncDeleteAll()
	return true
}

func (s *stage) displayStats(alpha *spriteAlpha) {
	const spacer = " -- "
	var headerSB strings.Builder

	weaponName, weaponLife := alpha.weaponStats()
	killStats := alpha.killStats()
	headerSB.WriteString(fmt.Sprintf("WEAPON: %s (%s)", weaponName, weaponLife))
	headerSB.WriteString(spacer)

	killStat := func(name string) string {
		if n, ok := killStats[name]; ok {
			return fmt.Sprint(n)
		}
		return "N/A"
	}
	killStatsText := fmt.Sprintf("KILLS: Beta: %s", killStat(betaName))
	if s.stageIdx > 0 {
		killStatsText += fmt.Sprintf(" | Gamma: %s", killStat(gammaName))
	}
	if s.stageIdx > 1 {
		killStatsText += fmt.Sprintf(" | Delta: %s", killStat(deltaName))
	}
	headerSB.WriteString(killStatsText)
	headerSB.WriteString(spacer)
	headerSB.WriteString(fmt.Sprintf("STAGE TIME LEFT: %s", func() string {
		timeLeft := stageDurations[s.stageIdx] - (s.m.g.MasterClock.Now() - s.stageStartTime)
		if timeLeft < 0 {
			timeLeft = time.Duration(0)
		}
		if s.stageIdx < totalStages-1 || timeLeft > 0 {
			return (timeLeft / time.Second * time.Second).String()
		}
		return "Until BOSS dies!"
	}()))
	s.m.winHeader.SetText(headerSB.String())

	s.m.winStats.SetText(fmt.Sprintf(`
Master clock: %s
Stage index: %d
%sArena Rect: %s
Alpha Rect: %s
FPS: %.0f
Total "pixels" rendered: %s
Memory usage: %s
%s`,
		time.Duration(s.m.g.MasterClock.Now()/(time.Second))*(time.Second),
		s.stageIdx+1,
		func() string {
			var ss []string
			if s.m.easyMode {
				ss = append(ss, "Easy Mode: On")
			}
			if s.m.invincible {
				ss = append(ss, "Invincible Mode: On")
				ss = append(ss, fmt.Sprintf("Hits on Alpha: %d", alpha.hits))
			}
			if len(ss) <= 0 {
				return ""
			}
			return strings.Join(ss, "\n") + "\n"
		}(),
		s.m.winArena.Rect(),
		func() string {
			if as, ok := s.m.g.SpriteMgr.TryFindByName(alphaName); ok {
				return as.Rect().String()
			}
			return "N/A"
		}(),
		s.m.g.FPS(),
		cwin.ByteSizeStr(s.m.g.WinSys.TotalChxRendered()),
		cwin.ByteSizeStr(s.m.g.HeapUsageInBytes()),
		s.m.g.SpriteMgr.DbgStats()))
}

func newStage(m *myGame, idx int) *stage {
	return &stage{m: m, stageIdx: idx}
}
