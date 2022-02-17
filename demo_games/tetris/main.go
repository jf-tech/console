package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
)

const (
	codeQuit int = iota
	codeReplay
	codeGameInitFailure

	// https://textkool.com/en/ascii-art-generator?hl=default&vl=default&font=Colossal&text=Game%20Over%20!
	gameOverTxt = `

 .d8888b.                                        .d88888b.                                 888
d88P  Y88b                                      d88P" "Y88b                                888
888    888                                      888     888                                888
888         8888b.  88888b.d88b.   .d88b.       888     888 888  888  .d88b.  888d888      888
888  88888     "88b 888 "888 "88b d8P  Y8b      888     888 888  888 d8P  Y8b 888P"        888
888    888 .d888888 888  888  888 88888888      888     888 Y88  88P 88888888 888          Y8P
Y88b  d88P 888  888 888  888  888 Y8b.          Y88b. .d88P  Y8bd8P  Y8b.     888           "
 "Y8888P88 "Y888888 888  888  888  "Y8888        "Y88888P"    Y88P    "Y8888  888          888


                            Press ESC or 'q' to quit, 'r' to replay.`
)

func main() {
	code := codeReplay
	for code == codeReplay {
		code = (&myGame{}).main()
	}
	os.Exit(code)
}

type myGame struct {
	g                   *cgame.Game
	winArena            *cwin.Win
	winNexts            []*cwin.Win
	winStats            *cwin.Win
	nexts               []*spritePiece
	s                   *spritePiece
	shadow              *spritePiece
	lw, lh              int
	board               [][]*spriteSettled
	shadowOff           bool
	naturalDropDelay    time.Duration // aka game speed
	kbSuspended         bool          // keyboard suspended for animation e.g piece dropping, etc
	totalRowsEliminated int
}

func (m *myGame) main() int {
	var err error
	m.g, err = cgame.Init(cterm.TCell)
	if err != nil {
		return codeGameInitFailure
	}
	defer m.g.Close()

	m.gameSetup()

	m.createCountDown()
	m.g.Run(cwin.Keys(cterm.KeyEsc, 'q'), nil, func(ev cterm.Event) cwin.MsgLoopResponseType {
		_, ok := m.g.Exchange.BoolData["countdown_done"]
		return cwin.TrueForMsgLoopStop(ok)
	})
	if m.g.IsGameOver() {
		goto game_over
	}

	m.g.SpriteMgr.CollidableRegistry().
		Register(pieceName, settledName).
		Register(shadowName, settledName)

	for i := 0; i < nextN; i++ {
		m.nexts = append(m.nexts,
			m.newSpritePiece(
				pieceName, pieceID(rand.Int()%int(pieceCount)), 0, cwin.ChAttr{Bg: pieceColor},
				m.winNexts[i], cwin.Point{X: 0, Y: 0}))
	}
	m.readyNextPieceForPlay()

	m.g.Run(cwin.Keys(cterm.KeyEsc, 'q'), cwin.Keys('p'), func(ev cterm.Event) cwin.MsgLoopResponseType {
		m.winStats.SetText(m.stats())
		// It's possible that m.s (current piece) is nil, so we must guard for that:
		// when a piece is settled, we call SpriteMgr.AsyncCallback to schedule
		// next piece readiness. And it's possible (or is it?) that before that queue
		// based callback event is processed, a keyboard press comes in, thus we
		// need this nil guard.
		if ev.Type != cterm.EventKey || m.s == nil || m.kbSuspended {
			return cwin.MsgLoopContinue
		}
		if ev.Key == cterm.KeyArrowUp {
			m.s.rotate()
		} else if ev.Key == cterm.KeyArrowLeft {
			m.s.move(-1)
		} else if ev.Key == cterm.KeyArrowRight {
			m.s.move(1)
		} else if ev.Key == cterm.KeyArrowDown {
			m.s.down()
		} else if ev.Ch == ' ' {
			m.s.directDrop()
		} else if ev.Ch == 's' {
			m.shadowOff = !m.shadowOff
			m.s.setupShadow()
		}
		return cwin.MsgLoopContinue
	})

game_over:
	ev := m.g.WinSys.MessageBoxEx(
		nil, append(gameOverKeys, replayGameKeys...), "uh oh...", gameOverTxt)
	if ev.Ch == 'r' {
		return codeReplay
	}
	return codeQuit
}

var (
	// due to console character aspect ratio, W:H 2:1 is a good ratio to achieve
	// (somewhat) square shape. Each block of the tetris pieces be 8x4.
	xscale         = 4
	yscale         = 2
	winArenaW      = LX2X(10)                     // standard tetris game is 10 block wide
	winArenaFrameW = 1 /*border*/ + winArenaW + 1 /*border*/
	winNextW       = LX2X(4)                      // max width of tetris pieces is 4 blocks (the bar)
	nextN          = 3
	winNextH       = nextN*LY2Y(2) + (nextN+1)*LY2Y(1)/2                        // next window has nextN pieces separated by single space.
	winNextFrameW  = 1 /*border*/ + LX2X(1) + winNextW + LX2X(1) /*space */ + 1 /*border*/
	winNextFrameH  = 1 /*border*/ + winNextH + 1                                /*border*/
	winInstrH      = 9
	winGameW       = winArenaFrameW + 1 /*space*/ + winNextFrameW

	gameOverKeys   = cwin.Keys(cterm.KeyEsc, 'q')
	replayGameKeys = cwin.Keys('r')

	directDropDelay      = 10 * time.Millisecond
	baseNaturalDropDelay = 500 * time.Millisecond
	settledFlyDelay      = 100 * time.Millisecond
)

// all x and y prefixed by 'l/L' means it's logical x/y representing units of tetris blocks.
// x / y withou 'l/L' prefix means physical terminal screen x's and y's.
func LX2X(lx int) int {
	return lx * xscale
}

func LY2Y(ly int) int {
	return ly * yscale
}

func X2LX(x int) int {
	return x / xscale
}

func Y2LY(y int) int {
	return y / yscale
}

func (m *myGame) gameSetup() {
	winSysClientR := m.g.WinSys.GetSysWin().ClientRect()
	h := int((winSysClientR.H-2)/yscale)*yscale + 2 // make sure we have no fractional rows
	winGame := m.g.WinSys.CreateWin(nil, cwin.WinCfg{
		R: cwin.Rect{
			X: (winSysClientR.W - winGameW) / 2,
			Y: (winSysClientR.H - h + 1) / 2,
			W: winGameW,
			H: h,
		},
		Name:     "game",
		NoBorder: true,
	})
	winArenaFrame := m.g.WinSys.CreateWin(winGame, cwin.WinCfg{
		R: cwin.Rect{
			X: 0,
			Y: 0,
			W: winArenaFrameW,
			H: winGame.ClientRect().H,
		},
		Name: fmt.Sprintf("Tetris (%dx%d)", X2LX(winArenaW), Y2LY(winGame.ClientRect().H)),
	})
	m.winArena = m.g.WinSys.CreateWin(winArenaFrame, cwin.WinCfg{
		R: cwin.Rect{
			X: 0,
			Y: 0,
			W: winArenaW,
			H: winArenaFrame.ClientRect().H},
		Name:     "arena",
		NoBorder: true,
	})
	m.lw, m.lh = X2LX(m.winArena.ClientRect().W), Y2LY(m.winArena.ClientRect().H)
	m.board = make([][]*spriteSettled, m.lh)
	for i := 0; i < m.lh; i++ {
		m.board[i] = make([]*spriteSettled, m.lw)
	}
	for i := 0; i < m.lh; i++ {
		for j := 0; j < m.lw; j++ {
			m.board[i][j] = nil
		}
	}
	winNextFrame := m.g.WinSys.CreateWin(winGame, cwin.WinCfg{
		R: cwin.Rect{
			X: winArenaFrame.Rect().W + 1,
			Y: 0,
			W: winNextFrameW,
			H: winNextFrameH,
		},
		Name: "Next",
	})
	for i := 0; i < nextN; i++ {
		m.winNexts = append(m.winNexts, m.g.WinSys.CreateWin(winNextFrame, cwin.WinCfg{
			R: cwin.Rect{
				X: LX2X(1),
				Y: LY2Y(1)/2 + i*(5*LY2Y(1)/2),
				W: winNextW,
				H: LY2Y(2),
			},
			Name:     fmt.Sprintf("next_%d", i),
			NoBorder: true,
			NoTitle:  true,
		}))
	}
	m.winStats = m.g.WinSys.CreateWin(winGame, cwin.WinCfg{
		R: cwin.Rect{
			X: winNextFrame.Rect().X,
			Y: winNextFrame.Rect().H,
			W: winNextFrame.Rect().W,
			H: winArenaFrame.Rect().H - winNextFrame.Rect().H - winInstrH,
		},
		Name: "Stats",
	})
	winInstr := m.g.WinSys.CreateWin(winGame, cwin.WinCfg{
		R: cwin.Rect{
			X: winNextFrame.Rect().X,
			Y: m.winStats.Rect().Y + m.winStats.Rect().H,
			W: winNextFrame.Rect().W,
			H: winInstrH,
		},
		Name:       "Keyboard",
		ClientAttr: cwin.ChAttr{Bg: cterm.ColorBlue},
	})
	winInstr.SetText(strings.Trim(fmt.Sprintf(`
%c / %c  : move
%c      : rotation
%c      : speed up
Space  : direct drop
's'    : shadow on/off
'p'    : pause/resume
ESC,'q': quit
`, cwin.DirRunes[cwin.DirLeft], cwin.DirRunes[cwin.DirRight],
		cwin.DirRunes[cwin.DirUp], cwin.DirRunes[cwin.DirDown]), "\n"))

	m.naturalDropDelay = baseNaturalDropDelay

	m.g.WinSys.Update()
	m.g.Resume()
}

func (m *myGame) createCountDown() {
	frames, framesR := cgame.FramesFromString([]string{
		`
 .d8888b.
d88P  Y88b
     .d88P
    8888"
     "Y8b.
888    888
Y88b  d88P
 "Y8888P"
`,
		`
 .d8888b.
d88P  Y88b
       888
     .d88P
 .od888P"
d88P"
888"
888888888
`,
		`
 d888
d8888
  888
  888
  888
  888
  888
8888888
`,
		`
 .d8888b.                888
d88P  Y88b               888
888    888               888
888         .d88b.       888
888  88888 d88""88b      888
888    888 888  888      Y8P
Y88b  d88P Y88..88P       "
 "Y8888P88  "Y88P"       888
`,
	}, cwin.ChAttr{Fg: cterm.ColorLightYellow})
	framesR.X = (m.winArena.ClientRect().W - framesR.W) / 2
	framesR.Y = (m.winArena.ClientRect().H - framesR.H) / 2
	s := cgame.NewSpriteBaseR(m.g, m.winArena, "count_down", frames[0], framesR)
	a := cgame.NewAnimatorFrame(s, cgame.AnimatorFrameCfg{
		Frames: cgame.NewSimpleFrameProvider(frames, 800*time.Millisecond, false),
		AnimatorCfgCommon: cgame.AnimatorCfgCommon{
			AfterFinish: func() {
				m.g.Exchange.BoolData["countdown_done"] = true
			}}})
	s.AddAnimator(a)
	m.g.SpriteMgr.AsyncCreateSprite(s)
}

func (m *myGame) newSpritePiece(spriteName string, pieceID pieceID, rotationIdx int,
	color cwin.ChAttr, parentW *cwin.Win, lxy cwin.Point) *spritePiece {
	f := cgame.SetAttrInFrame(mkFrame(pieceLibrary[pieceID][rotationIdx]), color)
	s := &spritePiece{
		SpriteBase:  cgame.NewSpriteBase(m.g, parentW, spriteName, f, LX2X(lxy.X), LY2Y(lxy.Y)),
		m:           m,
		pieceID:     pieceID,
		rotationIdx: rotationIdx,
		color:       color,
	}
	return s
}

// Represents a tetris piece in a given rotation.
type piece struct {
	blocks []cwin.Point
	// In order for a piece to have a natural look when it's rotated, we might need to
	// move it by a certain offset. E.g. z1 rotation 0 deg looks like this:
	//   12
	//    34
	// Now its rotation 90 deg looks like this:
	//    1
	//   32
	//   4
	// If we didn't do anything, when player rotates the piece, they would see a "sudden"
	// y-axis drop. Now imagine, if we anchor the rotation around block 2, then the rotation
	// looks natural. Which means, from frame 0 deg to frame 90 deg, we need a -1 y-axis offset.
	// offsetFromPrev records both axes' offsets from previous rotation.
	// Note the piece's first rotation has an offsetFromPrev based on the last rotation, as
	// rotation comes in circle.
	offsetFromPrev cwin.Point
}

type pieceID int

const (
	pieceID_Z1 = pieceID(iota)
	pieceID_Z2
	pieceID_L1
	pieceID_L2
	pieceID_Bar
	pieceID_T
	pieceID_SQ
	pieceCount
)

var (
	pieceLibrary = map[pieceID][]*piece{
		pieceID_Z1: {
			{
				blocks:         []cwin.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 1}},
				offsetFromPrev: cwin.Point{X: 0, Y: 1},
			},
			{
				blocks:         []cwin.Point{{X: 1, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}, {X: 0, Y: 2}},
				offsetFromPrev: cwin.Point{X: 0, Y: -1},
			},
		},
		pieceID_Z2: {
			{
				blocks:         []cwin.Point{{X: 1, Y: 0}, {X: 2, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}},
				offsetFromPrev: cwin.Point{X: -1, Y: 1},
			},
			{
				blocks:         []cwin.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 2}},
				offsetFromPrev: cwin.Point{X: 1, Y: -1},
			},
		},
		pieceID_L1: {
			{
				blocks:         []cwin.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 0, Y: 1}},
				offsetFromPrev: cwin.Point{X: 0, Y: 1},
			},
			{
				blocks:         []cwin.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}},
				offsetFromPrev: cwin.Point{X: 1, Y: -1},
			},
			{
				blocks:         []cwin.Point{{X: 2, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}},
				offsetFromPrev: cwin.Point{X: -1, Y: 1},
			},
			{
				blocks:         []cwin.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 1, Y: 2}},
				offsetFromPrev: cwin.Point{X: 0, Y: -1},
			},
		},
		pieceID_L2: {
			{
				blocks:         []cwin.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 1}},
				offsetFromPrev: cwin.Point{X: 0, Y: 1},
			},
			{
				blocks:         []cwin.Point{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}, {X: 0, Y: 2}},
				offsetFromPrev: cwin.Point{X: 1, Y: -1},
			},
			{
				blocks:         []cwin.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}},
				offsetFromPrev: cwin.Point{X: -1, Y: 1},
			},
			{
				blocks:         []cwin.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}},
				offsetFromPrev: cwin.Point{X: 0, Y: -1},
			},
		},
		pieceID_Bar: {
			{
				blocks:         []cwin.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}},
				offsetFromPrev: cwin.Point{X: -1, Y: 2},
			},
			{
				blocks:         []cwin.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}},
				offsetFromPrev: cwin.Point{X: 1, Y: -2},
			},
		},
		pieceID_T: {
			{
				blocks:         []cwin.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 1, Y: 1}},
				offsetFromPrev: cwin.Point{X: -1, Y: 1},
			},
			{
				blocks:         []cwin.Point{{X: 1, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 2}},
				offsetFromPrev: cwin.Point{X: 0, Y: -1},
			},
			{
				blocks:         []cwin.Point{{X: 1, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}},
				offsetFromPrev: cwin.Point{X: 0, Y: 0},
			},
			{
				blocks:         []cwin.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}, {X: 0, Y: 2}},
				offsetFromPrev: cwin.Point{X: 1, Y: 0},
			},
		},
		pieceID_SQ: {
			{
				blocks:         []cwin.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}},
				offsetFromPrev: cwin.Point{X: 0, Y: 0},
			},
		},
	}

	pieceName   = "active"
	shadowName  = "shadow"
	settledName = "settled"

	pieceColor   = cterm.ColorBlue
	shadowColor  = cterm.ColorDarkGray
	settledColor = cterm.ColorDarkGray
)

// Make one block of frame cells, without the color attributes.
func mkBlockFrameCell(lx, ly int) cgame.Frame {
	var f cgame.Frame
	f = append(f, cgame.Cell{X: 0, Y: 0, Chx: cwin.Chx{Ch: '⎡'}})
	f = append(f, cgame.Cell{X: xscale - 1, Y: 0, Chx: cwin.Chx{Ch: '⎤'}})
	f = append(f, cgame.Cell{X: xscale - 1, Y: yscale - 1, Chx: cwin.Chx{Ch: '⎦'}})
	f = append(f, cgame.Cell{X: 0, Y: yscale - 1, Chx: cwin.Chx{Ch: '⎣'}})
	for x := 1; x < xscale-1; x++ {
		f = append(f, cgame.Cell{X: x, Y: 0, Chx: cwin.Chx{Ch: '‾'}})
		f = append(f, cgame.Cell{X: x, Y: yscale - 1, Chx: cwin.Chx{Ch: '_'}})
	}
	for y := 1; y < yscale-1; y++ {
		f = append(f, cgame.Cell{X: 0, Y: y, Chx: cwin.Chx{Ch: '⎢'}})
		f = append(f, cgame.Cell{X: xscale - 1, Y: y, Chx: cwin.Chx{Ch: '⎥'}})
	}
	for i := 0; i < len(f); i++ {
		f[i].X += lx * xscale
		f[i].Y += ly * yscale
	}
	return f
}

func mkFrame(p *piece) cgame.Frame {
	var f cgame.Frame
	for _, b := range p.blocks {
		f = append(f, mkBlockFrameCell(b.X, b.Y)...)
	}
	return f
}

type spritePiece struct {
	*cgame.SpriteBase
	m           *myGame
	pieceID     pieceID
	rotationIdx int
	color       cwin.ChAttr
}

func (s *spritePiece) checkOrUpdatePosition(dlx, dly int) bool {
	return s.Update(cgame.UpdateArg{
		DXY: &cwin.Point{X: LX2X(dlx), Y: LY2Y(dly)},
		IBC: cgame.InBoundsCheckFullyVisible,
		CD:  cgame.CollisionDetectionOn,
	})
}

// looks like we can't just create a new piece for each rotation
// because new piece will have its own brand new drop waypoint, thus if you're
// rotating fast enough, the piece will never actually drop.
func (s *spritePiece) rotate() {
	lx, ly := X2LX(s.Rect().X), Y2LY(s.Rect().Y)
	ps := pieceLibrary[s.pieceID]
	nextIdx := (s.rotationIdx + 1) % len(ps)
	p := ps[nextIdx]
	lx += p.offsetFromPrev.X
	ly += p.offsetFromPrev.Y
	newS := s.m.newSpritePiece(
		pieceName, s.pieceID, nextIdx, s.color, s.m.winArena, cwin.Point{X: lx, Y: ly})
	if !newS.checkOrUpdatePosition(0, 0) {
		newS.Destroy()
		return
	}
	newS.setupShadow()
	newS.addDropAnimator(s.m.naturalDropDelay)
	s.Mgr().AsyncDeleteSprite(s)
	s.Mgr().AsyncCreateSprite(newS)
	s.m.s = newS
}

func (s *spritePiece) move(dlx int) {
	if s.checkOrUpdatePosition(dlx, 0) {
		s.setupShadow()
	}
}

func (s *spritePiece) setupShadow() {
	if s.m.shadow != nil {
		s.m.g.SpriteMgr.AsyncDeleteSprite(s.m.shadow)
		s.m.shadow = nil
	}
	if s.m.shadowOff {
		return
	}
	dly := 0
	for s.checkOrUpdatePosition(0, 1) {
		dly++
	}
	if dly <= 0 {
		return
	}
	s.checkOrUpdatePosition(0, -dly)
	shadow := s.m.newSpritePiece(shadowName, s.pieceID, s.rotationIdx, cwin.ChAttr{Fg: shadowColor},
		s.m.winArena, cwin.Point{X: X2LX(s.Rect().X), Y: Y2LY(s.Rect().Y) + dly})
	if !cgame.DetectCollision(s.Rect(), s.Frame(), shadow.Rect(), shadow.Frame()) {
		s.Mgr().AsyncCreateSprite(shadow)
		s.m.shadow = shadow
	} else {
		shadow.Destroy()
	}
}

func (s *spritePiece) down() {
	if s.checkOrUpdatePosition(0, 1) {
		// you would think going down shouldn't change the shadow?
		// but we might need to delete the shadow if the active piece is going down enough.
		s.setupShadow()
	}
}

func (s *spritePiece) directDrop() {
	s.addDropAnimator(directDropDelay)
	s.m.kbSuspended = true
}

func (s *spritePiece) addDropAnimator(dropDelay time.Duration) {
	for _, a := range s.Animators() {
		s.DeleteAnimator(a)
	}
	s.AddAnimator(cgame.NewAnimatorWaypoint(s.SpriteBase, cgame.AnimatorWaypointCfg{
		Waypoints: cgame.NewSimpleLoopWaypoints([]cgame.Waypoint{
			{
				DX: 0,
				DY: LY2Y(1),
				T:  dropDelay,
			},
		}),
		SingleMovePerWaypoint: true,
		AnimatorCfgCommon: cgame.AnimatorCfgCommon{
			InBoundsCheckType:      cgame.InBoundsCheckFullyVisible,
			CollisionDetectionType: cgame.CollisionDetectionOn,
			AfterUpdate: func() {
				s.setupShadow()
			},
			AfterFinish: func() {
				s.m.settlePiece(s)
				s.m.kbSuspended = false
			},
		}}))
}

// cgame.WaypointProvider
func (s *spritePiece) Next() (cgame.Waypoint, bool) {
	return cgame.Waypoint{
		DX: 0,
		DY: LY2Y(1),
		T:  s.m.naturalDropDelay,
	}, true
}

type spriteSettled struct {
	*cgame.SpriteBase
	m *myGame
}

func (ss *spriteSettled) addFlyOutAnimator(afterFinish func()) {
	ss.AddAnimator(cgame.NewAnimatorWaypoint(ss.SpriteBase, cgame.AnimatorWaypointCfg{
		Waypoints: cgame.NewSimpleWaypoints([]cgame.Waypoint{
			{
				DX: LX2X(ss.m.lw),
				DY: 0,
				T:  settledFlyDelay,
			},
		}),
		AnimatorCfgCommon: cgame.AnimatorCfgCommon{
			AfterFinish: afterFinish,
		},
	}))
}

func (ss *spriteSettled) addFlyDownAnimator(dly int) {
	ss.AddAnimator(cgame.NewAnimatorWaypoint(ss.SpriteBase, cgame.AnimatorWaypointCfg{
		Waypoints: cgame.NewSimpleWaypoints([]cgame.Waypoint{
			{
				DX: 0,
				DY: dly,
				T:  settledFlyDelay,
			},
		}),
		AnimatorCfgCommon: cgame.AnimatorCfgCommon{}}))
}

func (m *myGame) settlePiece(s *spritePiece) {
	r := s.Rect()
	for _, b := range pieceLibrary[s.pieceID][s.rotationIdx].blocks {
		lx, ly := X2LX(r.X)+b.X, Y2LY(r.Y)+b.Y
		settled := &spriteSettled{
			SpriteBase: cgame.NewSpriteBase(
				m.g, m.winArena, settledName,
				cgame.SetAttrInFrame(mkBlockFrameCell(0, 0), cwin.ChAttr{Bg: settledColor}),
				LX2X(lx), LY2Y(ly)),
			m: m,
		}
		m.g.SpriteMgr.AsyncCreateSprite(settled)
		m.board[ly][lx] = settled
	}
	// No need to destroy s (or m.s, the same thing) because the settlement coming from
	// the current piece's AfterFinish call from its waypoint animator, which would auto
	// destroy the sprite given we didn't use KeepAliveWhenFinished.
	// Simply clean up m.s for some sanity.
	m.s = nil
	m.g.SpriteMgr.AsyncFunc(func() {
		// the row elimination must take place asynchronously at the end of the event Q
		// so that the dropped piece will have finished turning into settled blocks.
		m.rowElimination()
	})
}

func (m *myGame) rowFull(ly int) bool {
	for lx := 0; lx < m.lw; lx++ {
		if m.board[ly][lx] == nil {
			return false
		}
	}
	return true
}

func (m *myGame) rowEmpty(ly int) bool {
	for lx := 0; lx < m.lw; lx++ {
		if m.board[ly][lx] != nil {
			return false
		}
	}
	return true
}

func (m *myGame) rowElimination() {
	rowsToDelete := 0
	for ly := 0; ly < m.lh; ly++ {
		if m.rowFull(ly) {
			rowsToDelete++
		}
	}
	if rowsToDelete <= 0 {
		m.readyNextPieceForPlay()
		return
	}
	m.totalRowsEliminated += rowsToDelete
	m.kbSuspended = true
	totalSettledFlyOut := rowsToDelete * m.lw
	// row elimination stage 1: first let all the full rows fly out.
	for ly := m.lh - 1; ly >= 0; ly-- {
		if !m.rowFull(ly) {
			continue
		}
		for lx := m.lw - 1; lx >= 0; lx-- {
			m.board[ly][lx].addFlyOutAnimator(func() {
				totalSettledFlyOut--
				if totalSettledFlyOut > 0 {
					return
				}
				m.rowCompact()
			})
			m.board[ly][lx] = nil
		}
	}
}

func (m *myGame) rowCompact() {
	for ly := m.lh - 1; ly >= 0; ly-- {
		if !m.rowEmpty(ly) {
			continue
		}
		anySettled := false
		for ly2 := ly; ly2 >= 1; ly2-- {
			for lx := 0; lx < m.lw; lx++ {
				if m.board[ly2-1][lx] != nil {
					m.board[ly2-1][lx].Update(cgame.UpdateArg{
						DXY: &cwin.Point{X: 0, Y: LY2Y(1)}})
					anySettled = true
				}
				m.board[ly2][lx] = m.board[ly2-1][lx]
			}
		}
		for lx := 0; lx < m.lw; lx++ {
			m.board[0][lx] = nil
		}
		if !anySettled {
			break
		}
		ly++
	}
	m.kbSuspended = false
	m.readyNextPieceForPlay()
}

func (m *myGame) readyNextPieceForPlay() {
	id := m.nexts[0].pieceID
	r := cgame.FrameRect(mkFrame(pieceLibrary[id][0]))
	m.s = m.newSpritePiece(
		pieceName, id, 0, cwin.ChAttr{Bg: pieceColor}, m.winArena,
		cwin.Point{X: X2LX((m.winArena.ClientRect().W - r.W) / 2), Y: 0})
	if !m.s.checkOrUpdatePosition(0, 0) {
		m.g.GameOver()
		return
	}
	m.s.setupShadow()
	m.s.addDropAnimator(m.naturalDropDelay)
	m.g.SpriteMgr.AsyncCreateSprite(m.s)

	for i := 0; i < len(m.nexts); i++ {
		m.nexts[i].Destroy()
		id := pieceID(rand.Int() % int(pieceCount))
		if i < len(m.nexts)-1 {
			id = m.nexts[i+1].pieceID
		}
		m.nexts[i] = m.newSpritePiece(
			pieceName, id, 0, cwin.ChAttr{Bg: pieceColor}, m.winNexts[i], cwin.Point{X: 0, Y: 0})
	}
}

func (m *myGame) stats() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Time: %s\n", m.g.MasterClock.Now().Round(time.Millisecond)))
	sb.WriteString(fmt.Sprintf("Rows completed: %d\n", m.totalRowsEliminated))
	sb.WriteString("\nDebug:\n")
	sb.WriteString("-----------------\n")
	sb.WriteString(fmt.Sprintf("FPS: %.0f\n", m.g.FPS()))
	sb.WriteString(fmt.Sprintf("Mem: %s\n", cwin.ByteSizeStr(m.g.HeapUsageInBytes())))
	sb.WriteString(m.g.SpriteMgr.DbgStats())
	return sb.String()
}
