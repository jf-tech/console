package main

import (
	"strings"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cwin"
	"github.com/nsf/termbox-go"
)

// Each gift pack ImgTxt frame contains exactly '???' placeholder to be replaced
// at runtime with the actual gift pack symbol. It must be of 3 runes to fit in.
type giftPackSymbol [3]rune

var (
	giftPackSpriteName = "gift_pack"
	// https://www.fileformat.info/info/unicode/char/25d0/index.htm
	giftPackImgTxts = []string{
		strings.Trim(`
⸨⸨???⸩⸩
`, "\n"),
		strings.Trim(`
⸨⸨???⸩⸩
`, "\n"),
		strings.Trim(`
⸨⸨???⸩⸩
`, "\n"),
		strings.Trim(`
⸨⸨???⸩⸩
`, "\n"),
	}
	giftPackBracketColors = []termbox.Attribute{
		termbox.ColorRed,
		termbox.ColorYellow,
		termbox.ColorGreen,
		termbox.ColorCyan,
	}
	giftPackAnimationSpeed = cgame.ActionPerSec(8)
	giftPackMoveSpeed      = cgame.ActionPerSec(5)
)

type spriteGiftPack struct {
	*cgame.SpriteAnimated
	gpSym giftPackSymbol
}

func (g *spriteGiftPack) Collided(other cgame.Sprite) {
	if other.Name() == alphaName {
		g.Mgr().AddEvent(cgame.NewSpriteEventDelete(g))
	}
}

func newSpriteGiftPack(g *cgame.Game, parent *cwin.Win, x, y int,
	sym giftPackSymbol, symAttr cwin.ChAttr) *spriteGiftPack {
	s := &spriteGiftPack{
		cgame.NewSpriteAnimated(g, parent,
			cgame.SpriteAnimatedCfg{
				Name:       giftPackSpriteName,
				Frames:     cgame.StringsToFrames(giftPackImgTxts, cwin.ChAttr{Fg: termbox.ColorWhite}),
				FrameSpeed: giftPackAnimationSpeed,
				Loop:       true,
				DY:         1,
				MoveSpeed:  giftPackMoveSpeed,
			},
			x, y),
		sym,
	}
	// Introduce rainbow colors to the brackets in the animation frames.
	// And replace ? with the gift pack symbol character.
	for i := 0; i < len(s.Config.Frames); i++ {
		for bracketIdx, symIdx, j := 0, 0, 0; j < len(s.Config.Frames[i]); j++ {
			switch s.Config.Frames[i][j].Chx.Ch {
			case '⸨', '⸩':
				s.Config.Frames[i][j].Chx.Attr.Fg =
					giftPackBracketColors[(i+bracketIdx)%len(giftPackBracketColors)]
				bracketIdx++
			case '?':
				s.Config.Frames[i][j].Chx.Ch = sym[symIdx]
				s.Config.Frames[i][j].Chx.Attr = symAttr
				symIdx++
			}
		}
	}
	return s
}

type giftPack struct {
	name          string
	sym           giftPackSymbol
	lifeTicker    *cgame.ActionPerSecTicker
	lifeRemaining time.Duration
}

func (gp *giftPack) remainingLife() time.Duration {
	gp.lifeRemaining -= time.Duration(gp.lifeTicker.HowMany()) * time.Second
	return gp.lifeRemaining
}

var (
	gpNoneSym     = giftPackSymbol{}
	gpNoneSymAttr = cwin.ChAttr{}

	gpShotgunName    = "Shotgun"
	gpShotgunSym     = giftPackSymbol{' ', 'S', ' '}
	gpShotgunSymAttr = cwin.ChAttr{Fg: termbox.ColorBlue, Bg: termbox.ColorWhite}
	gpShotgunLife    = time.Minute
	gpShotgunProb    = 100000

	gpShotgun2Name    = "Shotgun++"
	gpShotgun2Sym     = giftPackSymbol{'S', '+', '+'}
	gpShotgun2SymAttr = cwin.ChAttr{Fg: termbox.ColorBlue, Bg: termbox.ColorWhite}
	gpShotgun2Life    = time.Minute
	gpShotgun2Prob    = gpShotgunProb * 2
)

func newGiftPackShotgun(clock *cgame.Clock) *giftPack {
	return &giftPack{
		name:          gpShotgunName,
		sym:           gpShotgunSym,
		lifeTicker:    cgame.NewActionPerSecTicker(clock, 1, true),
		lifeRemaining: gpShotgunLife,
	}
}

func newGiftPackShotgun2(clock *cgame.Clock) *giftPack {
	return &giftPack{
		name:          gpShotgun2Name,
		sym:           gpShotgun2Sym,
		lifeTicker:    cgame.NewActionPerSecTicker(clock, 1, true),
		lifeRemaining: gpShotgun2Life,
	}
}

func genGiftPack() (giftPackSymbol, cwin.ChAttr, bool) {
	if testProb(gpShotgunProb) {
		return gpShotgunSym, gpShotgunSymAttr, true
	}
	if testProb(gpShotgun2Prob) {
		return gpShotgun2Sym, gpShotgun2SymAttr, true
	}
	return gpNoneSym, gpNoneSymAttr, false
}
