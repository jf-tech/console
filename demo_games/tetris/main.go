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
	g        *cgame.Game
	winArena *cwin.Win
	winNexts []*cwin.Win
	winStats *cwin.Win
	next     []*spritePiece
	s        *spritePiece
	shadow   *spritePiece
	lw, lh   int
	board    [][]cgame.Sprite
	speed    time.Duration
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
	m.g.Run(nil, nil, func(ev cterm.Event) bool {
		_, ok := m.g.Exchange.BData["countdown_done"]
		return ok
	})

	m.s = m.genRandomSpritePiece(m.winArena, cwin.Point{X: 4, Y: 0})
	m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(m.s, m.s.createNaturalDropAnimator()))
	m.s.setupShadow(m.s)
	m.g.Run(cwin.Keys(cterm.KeyEsc, 'q'), cwin.Keys('p'), func(ev cterm.Event) bool {
		if ev.Type == cterm.EventKey {
			if ev.Key == cterm.KeyArrowUp {
				m.s.rotate()
			} else if ev.Key == cterm.KeyArrowLeft {
				m.s.move(-1)
			} else if ev.Key == cterm.KeyArrowRight {
				m.s.move(1)
			} else if ev.Key == cterm.KeyArrowDown {
				m.s.down()
			}
		}
		return false
	})

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
	winInstrH      = 8
	winGameW       = winArenaFrameW + 1 /*space*/ + winNextFrameW
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

var (
	gameOverKeys   = cwin.Keys(cterm.KeyEsc, 'q')
	replayGameKeys = cwin.Keys('r')
)

func (m *myGame) gameSetup() {
	winSysClientR := m.g.WinSys.GetSysWin().ClientRect()
	winGame := m.g.WinSys.CreateWin(nil, cwin.WinCfg{
		R: cwin.Rect{
			X: (winSysClientR.W - winGameW) / 2,
			Y: 0,
			W: winGameW,
			H: winSysClientR.H},
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
	m.board = make([][]cgame.Sprite, m.lh)
	for i := 0; i < m.lh; i++ {
		m.board[i] = make([]cgame.Sprite, m.lw)
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
		Name:    "next_frame",
		NoTitle: true,
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
%c / %c   : move
%c       : rotation
%c       : speed up
Space   : direct drop
'p'     : pause/resume
ESC,'q' : quit
`, cgame.DirSymbols[cgame.DirLeft], cgame.DirSymbols[cgame.DirRight],
		cgame.DirSymbols[cgame.DirUp], cgame.DirSymbols[cgame.DirDown]), "\n"))

	for i := 0; i < nextN; i++ {
		m.next = append(m.next, m.genRandomSpritePiece(m.winNexts[i], cwin.Point{X: 0, Y: 0}))
	}

	m.speed = 500 * time.Millisecond

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
	a := cgame.NewAnimatorFrame(cgame.AnimatorFrameCfg{
		Frames: cgame.NewSimpleFrameProvider(frames, 800*time.Millisecond, false),
		AfterFinish: func(s cgame.Sprite) {
			m.g.Exchange.BData["countdown_done"] = true
		},
	})
	framesR.X = (m.winArena.ClientRect().W - framesR.W) / 2
	framesR.Y = (m.winArena.ClientRect().H - framesR.H) / 2
	s := cgame.NewSpriteBaseR(m.g, m.winArena, "count_down", frames[0], framesR)
	m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(s, a))
}

func (m *myGame) createSpritePiece(
	name string, idx int, color cterm.Attribute, w *cwin.Win, lxy cwin.Point) *spritePiece {
	f := cgame.SetAttrInFrame(cgame.CopyFrame(pieceLib[name][idx].f), cwin.ChAttr{Bg: color})
	s := &spritePiece{
		SpriteBase: cgame.NewSpriteBase(
			m.g, w, name, f,
			LX2X(lxy.X),
			LY2Y(lxy.Y)),
		color: color,
		idx:   idx,
		m:     m,
	}
	return s
}

func (m *myGame) genRandomSpritePiece(w *cwin.Win, lxy cwin.Point) *spritePiece {
	name := pieceNames[rand.Int()%len(pieceNames)]
	return m.createSpritePiece(name, 0, pieceActiveColor, w, lxy)
}

// Make one block of frame cells, without the color attributes.
func mkBlock(lx, ly int) cgame.Frame {
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

// Represents a tetris piece in a given rotation.
type piece struct {
	f cgame.Frame
	// f is the frame that contains the shape of the piece in the given rotation.
	// In order for it to have a natural look when it's rotated, we might need to
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

func mkPiece(lxy []cwin.Point, offsetlxy cwin.Point) *piece {
	p := &piece{}
	for _, e := range lxy {
		p.f = append(p.f, mkBlock(e.X, e.Y)...)
	}
	p.offsetFromPrev = offsetlxy
	return p
}

var (
	pieceZ1    = "Z1"
	pieceNames = []string{pieceZ1}
	pieceLib   = map[string][]*piece{
		pieceZ1: {
			mkPiece(
				[]cwin.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 1}},
				cwin.Point{X: 0, Y: 1},
			),
			mkPiece(
				[]cwin.Point{{X: 1, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}, {X: 0, Y: 2}},
				cwin.Point{X: 0, Y: -1},
			),
		},
	}
	pieceActiveColor  = cterm.ColorBlue
	pieceShadowColor  = cterm.ColorBlack
	pieceSettledColor = cterm.ColorDarkGray
)

type spritePiece struct {
	*cgame.SpriteBase
	color cterm.Attribute
	idx   int
	m     *myGame
}

func (s *spritePiece) destroy() {
	if _, found := s.m.g.SpriteMgr.TryFindByUID(s.UID()); found {
		s.m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventDelete(s))
		return
	}
	s.m.g.WinSys.RemoveWin(s.Win())
}

func (*spritePiece) check(s *spritePiece, dxy cwin.Point) bool {
	lx, ly := X2LX(s.Win().Rect().X)+dxy.X, Y2LY(s.Win().Rect().Y)+dxy.Y
	lw, lh := X2LX(s.Win().Rect().W), Y2LY(s.Win().Rect().H)
	if lx < 0 || lx+lw > s.m.lw ||
		ly < 0 || ly+lh > s.m.lh {
		return false
	}
	f := cgame.FrameFromWin(s.Win())
	for i := 0; i < len(f); i++ {
		flx, fly := X2LX(s.Win().Rect().X+f[i].X)+dxy.X, Y2LY(s.Win().Rect().Y+f[i].Y)+dxy.Y
		if s.m.board[fly][flx] != nil {
			return false
		}
	}
	return true
}

// looks like we can't just create a new piece for each rotation
// because new piece will have its own brand new drop waypoint, thus if you're
// rotating fast enough, the piece will never actually drop.
func (s *spritePiece) rotate() {
	lx, ly := X2LX(s.Win().Rect().X), Y2LY(s.Win().Rect().Y)
	ps := pieceLib[s.Name()]
	nextIdx := (s.idx + 1) % len(ps)
	p := ps[nextIdx]
	lx += p.offsetFromPrev.X
	ly += p.offsetFromPrev.Y
	newS := s.m.createSpritePiece(s.Name(), nextIdx, s.color,
		s.Win().Parent(), cwin.Point{X: lx, Y: ly})
	if !s.check(newS, cwin.Point{}) {
		newS.destroy()
		return
	}
	s.Mgr().AddEvent(cgame.NewSpriteEventDelete(s))
	s.m.s = newS
	s.Mgr().AddEvent(cgame.NewSpriteEventCreate(newS, newS.createNaturalDropAnimator()))
	s.setupShadow(newS)
}

func (s *spritePiece) move(dlx int) {
	if !s.check(s, cwin.Point{X: dlx}) {
		return
	}
	s.Win().SetPosRelative(LX2X(dlx), 0)
	s.setupShadow(s)
}

func (*spritePiece) setupShadow(s *spritePiece) {
	if s.m.shadow != nil {
		s.m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventDelete(s.m.shadow))
		s.m.shadow = nil
	}
	for y := s.m.lh - 1; y >= Y2LY(s.Win().Rect().Y+s.Win().Rect().H); y-- {
		shadow := s.m.createSpritePiece(s.Name(),
			s.idx, pieceShadowColor,
			s.Win().Parent(), cwin.Point{X: X2LX(s.Win().Rect().X), Y: y})
		if !s.check(shadow, cwin.Point{}) {
			shadow.destroy()
			continue
		}
		s.m.shadow = shadow
		s.m.g.SpriteMgr.AddEvent(cgame.NewSpriteEventCreate(s.m.shadow))
		break
	}
}

func (s *spritePiece) down() {
	if !s.check(s, cwin.Point{Y: 1}) {
		return
	}
	s.Win().SetPosRelative(0, LY2Y(1))
	s.setupShadow(s)
}

func (s *spritePiece) createNaturalDropAnimator() cgame.Animator {
	return cgame.NewAnimatorWaypoint(cgame.AnimatorWaypointCfg{Waypoints: s})
}

func (s *spritePiece) Next() (cgame.Waypoint, bool) {
	ly := Y2LY(s.Win().Rect().Y)
	if ly > 15 {
		return cgame.Waypoint{}, false
	}
	return cgame.Waypoint{
		Type: cgame.WaypointRelative,
		X:    0,
		Y:    LY2Y(1),
		T:    s.m.speed,
	}, true
}
