package main

import (
	"path"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cutil"
	"github.com/jf-tech/console/cwin"
)

var (
	gameOverKeys       = cwin.Keys(cterm.KeyEsc, 'q')
	pauseGameKeys      = cwin.Keys('p')
	replayGameKeys     = cwin.Keys('r')
	skipStageKeys      = cwin.Keys('s')
	easyModeKeys       = cwin.Keys('e')
	invincibleModeKeys = cwin.Keys('i')

	sfxBackgroundVol = float64(-4)
	sfxClipVol       = float64(-1)

	sfxAllFiles = []string{}
	sfxFile     = func(relpath string) string {
		sfxAllFiles = append(sfxAllFiles, path.Join(cutil.GetCurFileDir(), relpath))
		return relpath
	}
	sfxBackgroundFile     = sfxFile("resources/background.mp3")
	sfxGameStartFile      = sfxFile("resources/game-start.mp3")
	sfxPewFile            = sfxFile("resources/pew.mp3")
	sfxOuchFile           = sfxFile("resources/ouch.mp3")
	sfxBossIsHereFile     = sfxFile("resources/boss-is-here.mp3")
	sfxGameOverFile       = sfxFile("resources/game-over.mp3")
	sfxWeaponUpgradedFile = sfxFile("resources/weapon-upgraded.mp3")
	sfxYouWonFile         = sfxFile("resources/you-won.mp3")

	totalStages                    = 3
	stageDurations                 = []time.Duration{time.Minute, time.Minute, time.Minute}
	stageIntroBannerInOutDuration  = 300 * time.Millisecond
	stageIntroBannerStayDuration   = 1 * time.Second
	stagePassedBannerInOutDuration = stageIntroBannerInOutDuration
	stagePassedBannerStayDuration  = 2 * time.Second

	alphaBulletAttr = cwin.ChAttr{Fg: cterm.ColorLightYellow}
	enemyBulletAttr = cwin.ChAttr{Fg: cterm.ColorLightCyan}

	bgStarSpeed   = cgame.CharPerSec(25)
	bgStarGenProb = cutil.NewPeriodicProbabilityChecker("50%", 100*time.Millisecond)

	alphaBulletSpeed = cgame.CharPerSec(30)

	betaSpeed           = cgame.CharPerSec(4)
	betaBulletSpeed     = cgame.CharPerSec(10)
	betaGenProbPerStage = []*cutil.PeriodicProbabilityChecker{
		cutil.NewPeriodicProbabilityChecker("0.6%", 10*time.Millisecond),
		cutil.NewPeriodicProbabilityChecker("0.5%", 10*time.Millisecond),
		cutil.NewPeriodicProbabilityChecker("0.4%", 10*time.Millisecond),
	}
	betaFiringProbPerStage    = []string{"10%", "10%", "10%"}
	betaFiringPelletsPerStage = []int{2, 3, 3}
	betaExplosionDuration     = 2 * time.Second

	gammaSpeed           = cgame.CharPerSec(4)
	gammaBulletSpeed     = cgame.CharPerSec(10)
	gammaGenProbPerStage = []*cutil.PeriodicProbabilityChecker{
		cutil.NewPeriodicProbabilityChecker("0%", 10*time.Millisecond),
		cutil.NewPeriodicProbabilityChecker("0.2%", 10*time.Millisecond),
		cutil.NewPeriodicProbabilityChecker("0.2%", 10*time.Millisecond),
	}
	gammaFiringProbPerStage = []string{"0%", "10%", "10%"}
	gammaExplosionDuration  = 2 * time.Second

	deltaVerticalSpeed     = cgame.CharPerSec(35)
	deltaHorizontalSpeed   = deltaVerticalSpeed * 2
	deltaSpeedDiscountEasy = 0.6
	deltaGenProbPerStage   = []*cutil.PeriodicProbabilityChecker{
		cutil.NewPeriodicProbabilityChecker("0%", 10*time.Millisecond),
		cutil.NewPeriodicProbabilityChecker("0%", 10*time.Millisecond),
		cutil.NewPeriodicProbabilityChecker("0.4%", 10*time.Millisecond),
	}
	deltaVerticalProb      = "50%"
	deltaExplosionDuration = 1 * time.Second

	bossSpeed                      = cgame.CharPerSec(2)
	bossMinDistToGoBeforeDirChange = 8
	bossMaxDistToGoBeforeDirChange = 20
	bossHP                         = 200
	bossBulletFiringProb           = "20%"
	bossBulletSpeed                = cgame.CharPerSec(10)
	bossExplosionDuration          = 4 * time.Second

	giftPackMoveSpeed  = cgame.CharPerSec(5)
	gpShotgunLife      = time.Minute
	gpShotgunProb      = cutil.NewPeriodicProbabilityChecker("8%", time.Second)
	gpShotgun2Life     = time.Minute
	gpShotgun2Prob     = cutil.NewPeriodicProbabilityChecker("2%", time.Second)
	gpShotgun2ProbEasy = cutil.NewPeriodicProbabilityChecker("10%", time.Second)

	exchangeGiftPackWeapon = "giftpack_weapon"
)
