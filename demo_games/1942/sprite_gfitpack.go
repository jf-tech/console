package main

import (
	"math/rand"
	"strings"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
)

// Each gift pack frame contains exactly '???' placeholder to be replaced
// at runtime with the actual gift pack symbol. It must be of 3 runes to fit in.
type giftPackSymbol [3]rune

var (
	giftPackName              = "gift_pack"
	giftPackSymbolPlaceholder = giftPackSymbol{'?', '?', '?'}
	// https://www.fileformat.info/info/unicode/char/25d0/index.htm
	giftPackFrameTxt = `⸨⸨???⸩⸩`
)

type spriteGiftPack struct {
	*cgame.SpriteBase
	gpSym giftPackSymbol
}

func (g *spriteGiftPack) CollisionNotify(_ bool, _ []cgame.Sprite) cgame.CollisionResponseType {
	g.Mgr().DeleteSprite(g)
	return cgame.CollisionResponseJustDoIt
}

func createGiftPack(m *myGame, sym giftPackSymbol, symAttr cwin.ChAttr) {
	frame := cgame.FrameFromString(strings.ReplaceAll(
		giftPackFrameTxt, string(giftPackSymbolPlaceholder[:]), string(sym[:])), symAttr)
	s := &spriteGiftPack{
		SpriteBase: cgame.NewSpriteBase(m.g, m.winArena, giftPackName, frame,
			rand.Int()%(m.winArena.ClientRect().W-cgame.FrameRect(frame).W), 0),
		gpSym: sym}
	dist := 1000 // large enough to go out of window (and auto destroy)
	a := cgame.NewAnimatorWaypoint(s.SpriteBase, cgame.AnimatorWaypointCfg{
		Waypoints: cgame.NewSimpleWaypoints([]cgame.Waypoint{
			{
				Type: cgame.WaypointRelative,
				X:    0,
				Y:    1 * dist,
				T:    time.Duration((float64(dist) / float64(giftPackMoveSpeed)) * float64(time.Second)),
			}})})
	s.AddAnimator(a)
	m.g.SpriteMgr.AddSprite(s)
}

type giftPack struct {
	name  string
	sym   giftPackSymbol
	life  time.Duration
	clock *cgame.Clock
	start time.Duration
}

func (gp *giftPack) remainingLife() time.Duration {
	elapsed := gp.clock.Now() - gp.start
	if gp.life <= elapsed {
		return 0
	}
	return gp.life - elapsed
}

var (
	gpShotgunName    = "Shotgun"
	gpShotgunSym     = giftPackSymbol{'-', 'S', '-'}
	gpShotgunSymAttr = cwin.ChAttr{Fg: cterm.ColorWhite, Bg: cterm.ColorBlack}

	gpShotgun2Name    = "Shotgun++"
	gpShotgun2Sym     = giftPackSymbol{'S', '+', '+'}
	gpShotgun2SymAttr = cwin.ChAttr{Fg: cterm.ColorLightYellow, Bg: cterm.ColorBlack}
)

func newGiftPackShotgun(clock *cgame.Clock) *giftPack {
	return &giftPack{
		name:  gpShotgunName,
		sym:   gpShotgunSym,
		life:  gpShotgunLife,
		clock: clock,
		start: clock.Now(),
	}
}

func newGiftPackShotgun2(clock *cgame.Clock) *giftPack {
	return &giftPack{
		name:  gpShotgun2Name,
		sym:   gpShotgun2Sym,
		life:  gpShotgun2Life,
		clock: clock,
		start: clock.Now(),
	}
}
