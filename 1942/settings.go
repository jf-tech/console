package main

import (
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cwin"
	"github.com/nsf/termbox-go"
)

var (
	gameOverKeys  = cwin.Keys(termbox.KeyEsc, 'q')
	pauseGameKeys = cwin.Keys('p')

	stageDurations                 = []time.Duration{time.Minute, time.Minute, time.Minute}
	stageIntroBannerInOutDuration  = 300 * time.Millisecond
	stageIntroBannerStayDuration   = 1 * time.Second
	stagePassedBannerInOutDuration = stageIntroBannerInOutDuration
	stagePassedBannerStayDuration  = 2 * time.Second

	bgStarSpeed   = cgame.CharPerSec(25)
	bgStarGenProb = cgame.NewPeriodicProbabilityChecker("50%", 100*time.Millisecond)

	alphaBulletSpeed = cgame.CharPerSec(30)

	betaSpeed           = cgame.CharPerSec(4)
	betaBulletSpeed     = cgame.CharPerSec(10)
	betaGenProbPerStage = []*cgame.PeriodicProbabilityChecker{
		cgame.NewPeriodicProbabilityChecker("0.6%", 10*time.Millisecond),
		cgame.NewPeriodicProbabilityChecker("0.5%", 10*time.Millisecond),
		cgame.NewPeriodicProbabilityChecker("0.5%", 10*time.Millisecond),
	}
	betaFiringProbPerStage    = []string{"10%", "10%", "10%"}
	betaFiringPelletsPerStage = []int{3, 4, 5}

	gammaSpeed           = cgame.CharPerSec(4)
	gammaBulletSpeed     = cgame.CharPerSec(10)
	gammaGenProbPerStage = []*cgame.PeriodicProbabilityChecker{
		cgame.NewPeriodicProbabilityChecker("0%", 10*time.Millisecond),
		cgame.NewPeriodicProbabilityChecker("0.2%", 10*time.Millisecond),
		cgame.NewPeriodicProbabilityChecker("0.3%", 10*time.Millisecond),
	}
	gammaFiringProbPerStage = []string{"0%", "10%", "10%"}

	deltaSpeed           = cgame.CharPerSec(35)
	deltaGenProbPerStage = []*cgame.PeriodicProbabilityChecker{
		cgame.NewPeriodicProbabilityChecker("0%", 10*time.Millisecond),
		cgame.NewPeriodicProbabilityChecker("0%", 10*time.Millisecond),
		cgame.NewPeriodicProbabilityChecker("0.3%", 10*time.Millisecond),
	}

	giftPackMoveSpeed  = cgame.CharPerSec(5)
	gpShotgunLife      = time.Minute
	gpShotgunProb      = cgame.NewPeriodicProbabilityChecker("8%", time.Second)
	gpShotgun2Life     = time.Minute
	gpShotgun2Prob     = cgame.NewPeriodicProbabilityChecker("2%", time.Second)
	gpShotgun2ProbEasy = cgame.NewPeriodicProbabilityChecker("10%", time.Second)
)
