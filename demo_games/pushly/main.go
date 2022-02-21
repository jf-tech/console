package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cgame/assets"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cutil"
	"github.com/jf-tech/console/cwin"
	"github.com/jf-tech/console/cwin/ccomp"
	"github.com/jf-tech/go-corelib/ios"
)

func main() {
	for (&myGame{}).main() == assets.GameResultReplay {
	}
}

type myGame struct {
	g                *cgame.Game
	winArenaFrame    cwin.Win
	winArena         cwin.Win
	winLevelLB       cwin.Win
	arenaTitlePrefix string
	lw, lh           int

	lvlCount int

	lvl          int
	levelChanged bool
	board        [][]*sprite // note target sprites not on board.
	targets      []*sprite
	steps        int
}

func (m *myGame) main() assets.GameResult {
	var err error
	m.g, err = cgame.Init(cterm.TCell)
	if err != nil {
		return assets.GameResultSystemFailure
	}
	defer m.g.Close()

	m.gameSetup()

	m.g.SpriteMgr.CollidableRegistry().
		RegisterBulk(alphaName, []string{brickName, concreteName}).
		RegisterBulk(brickName, []string{alphaName, brickName, concreteName})

	m.lvl = 0
	m.levelChanged = true
	for {
		m.loadLevel(m.lvl)
		m.levelChanged = false
		m.g.WinSys.Update()
		m.winLevelLB.(*ccomp.ListBox).SetSelected(m.lvl)
		m.winArenaFrame.SetTitle("%s - Level %d", m.arenaTitlePrefix, m.lvl+1)
		m.g.WinSys.MessageBox(
			nil, "Info", "Level %d ready! Press Enter to start...", m.lvl+1)
		replay := false
		m.g.Run(assets.GameOverKeys, nil, func(ev cterm.Event) cwin.EventResponse {
			if ev.Type == cterm.EventKey {
				switch ev.Key {
				case cterm.KeyArrowUp:
					m.alpha().push(cwin.DirUp)
				case cterm.KeyArrowRight:
					m.alpha().push(cwin.DirRight)
				case cterm.KeyArrowDown:
					m.alpha().push(cwin.DirDown)
				case cterm.KeyArrowLeft:
					m.alpha().push(cwin.DirLeft)
				}
				switch ev.Ch {
				case 'r':
					replay = true
					return cwin.EventLoopStop
				case 's':
					m.g.WinSys.SetFocus(m.winLevelLB)
					return cwin.EventHandled
				}
				if m.checkLevelClear() {
					return cwin.EventLoopStop
				}
			}
			if m.levelChanged {
				return cwin.EventLoopStop
			}
			return cwin.EventHandled
		})
		if m.g.IsGameOver() {
			break
		}
		if m.levelChanged || replay {
			continue
		}
		if m.lvl < m.lvlCount-1 {
			m.g.SoundMgr.PlayMP3(sfxLevelClear, 0, 1)
		}
		m.g.WinSys.MessageBox(
			nil, "Good job!", "Level %d complete! Press Enter for next level...", m.lvl+1)
		if m.lvl == m.lvlCount-1 {
			break
		}
		m.lvl++
		m.levelChanged = true
	}
	m.g.SoundMgr.PlayMP3(sfxGameOver, 0, 1)
	return assets.DisplayGameOverDialog(m.g)
}

var (
	lvlKeySpace    = '.'
	lvlKeyAlpha    = 'A'
	lvlKeyBrick    = 'B'
	lvlKeyConcrete = 'C'
	lvlKeyTarget   = 'X'

	xscale         = 8
	yscale         = 3
	arenaLW        = 8
	arenaLH        = 8
	winArenaW      = LX2X(arenaLW) - 1
	winArenaH      = LY2Y(arenaLH) - 1
	winArenaFrameW = 1 /*border*/ + winArenaW + 1 /*border*/
	winArenaFrameH = 1 /*border*/ + winArenaH + 1 /*border*/
	winInstrW      = 34
	winInstrH      = 12
	winLevelLBW    = winInstrW
	winLevelLBH    = winArenaFrameH - winInstrH
	winGameW       = winArenaFrameW + 1 /*space*/ + winLevelLBW
	winGameH       = winArenaFrameH

	sfxFile = func(relpath string) string {
		return path.Join(cutil.GetCurFileDir(), relpath)
	}
	sfxClick      = sfxFile("resources/click.mp3")
	sfxLevelClear = sfxFile("resources/level_clear.mp3")
	sfxGameOver   = sfxFile("resources/gameover.mp3")
)

func (m *myGame) gameSetup() {
	winSysClientR := m.g.WinSys.SysWin().ClientRect()
	winGame := m.g.WinSys.CreateWin(nil, cwin.WinCfg{
		R: cwin.Rect{
			X: (winSysClientR.W - winGameW) / 2,
			Y: (winSysClientR.H - winGameH) / 2,
			W: winGameW,
			H: winGameH,
		},
		Name:     "game",
		NoBorder: true,
	})
	m.lw, m.lh = arenaLW, arenaLH
	m.arenaTitlePrefix = fmt.Sprintf("Pushly (%dx%d)", m.lw, m.lh)
	m.winArenaFrame = m.g.WinSys.CreateWin(winGame, cwin.WinCfg{
		R: cwin.Rect{
			X: 0,
			Y: 0,
			W: winArenaFrameW,
			H: winArenaFrameH,
		},
		Name: m.arenaTitlePrefix,
	})
	m.g.WinSys.SetFocus(m.winArenaFrame)
	m.winArena = m.g.WinSys.CreateWin(m.winArenaFrame, cwin.WinCfg{
		R: cwin.Rect{
			X: 0,
			Y: 0,
			W: winArenaW,
			H: winArenaH},
		Name:     "arena",
		NoBorder: true,
	})
	for ly := 1; ly < m.lh; ly++ {
		m.winArena.FillClient(
			cwin.Rect{X: 0, Y: ly*(yscale+1) - 1, W: winArenaW, H: 1},
			cwin.Chx{Ch: '─', Attr: cwin.Attr{Fg: cterm.ColorDarkGray}})
	}
	for lx := 1; lx < m.lw; lx++ {
		m.winArena.FillClient(
			cwin.Rect{X: lx*(xscale+1) - 1, Y: 0, W: 1, H: winArenaH},
			cwin.Chx{Ch: '│', Attr: cwin.Attr{Fg: cterm.ColorDarkGray}})
	}
	m.lvlCount = m.detectLevelCount()
	m.winLevelLB = ccomp.CreateListBox(m.g.WinSys, winGame, ccomp.ListBoxCfg{
		WinCfg: cwin.WinCfg{
			R: cwin.Rect{
				X: winArenaFrameW + 1,
				Y: 0,
				W: winLevelLBW,
				H: winLevelLBH,
			},
			Name: "Available Levels",
		},
		Items: func() []string {
			var lvls []string
			for i := 0; i < m.lvlCount; i++ {
				lvls = append(lvls, fmt.Sprintf("Level %d", i+1))
			}
			return lvls
		}(),
		EnterKeyToSelect: true,
		OnSelect: func(idx int, selected string) {
			m.lvl = idx
			m.levelChanged = true
			m.g.WinSys.SetFocus(m.winArenaFrame)
		},
	})
	winInstr := m.g.WinSys.CreateWin(winGame, cwin.WinCfg{
		R: cwin.Rect{
			X: m.winLevelLB.Rect().X,
			Y: m.winLevelLB.Rect().Y + m.winLevelLB.Rect().H,
			W: winInstrW,
			H: winInstrH,
		},
		Name:       "Keyboard",
		ClientAttr: cwin.Attr{Bg: cterm.ColorBlue},
	})
	winInstr.SetText(fmt.Sprintf(`
  %c
%c %c %c   : to push a brick

 's'    : select a level

 'r'    : replay current level

ESC,'q' : quit

`,
		cwin.DirRunes[cwin.DirUp], cwin.DirRunes[cwin.DirLeft],
		cwin.DirRunes[cwin.DirDown], cwin.DirRunes[cwin.DirRight]))

	m.board = make([][]*sprite, m.lh)
	for i := 0; i < m.lh; i++ {
		m.board[i] = make([]*sprite, m.lw)
	}

	m.g.WinSys.Update()
	m.g.Resume()
}

func (m *myGame) createAlpha(lx, ly int) {
	m.board[ly][lx] = &sprite{
		SpriteBase: cgame.NewSpriteBase(
			m.g, m.winArena, alphaName, alpahFrameUp, LX2X(lx), LY2Y(ly)),
		m: m,
	}
	m.g.SpriteMgr.AddSprite(m.alpha())
}

func (m *myGame) createBrick(lx, ly int) {
	m.board[ly][lx] = &sprite{
		SpriteBase: cgame.NewSpriteBase(
			m.g, m.winArena, brickName, brickFrame, LX2X(lx), LY2Y(ly)),
		m: m,
	}
	m.g.SpriteMgr.AddSprite(m.board[ly][lx])
}

func (m *myGame) createConcrete(lx, ly int) {
	m.board[ly][lx] = &sprite{
		SpriteBase: cgame.NewSpriteBase(
			m.g, m.winArena, concreteName, concreteFrame, LX2X(lx), LY2Y(ly)),
		m: m,
	}
	m.g.SpriteMgr.AddSprite(m.board[ly][lx])
}

func (m *myGame) createTarget(lx, ly int) {
	m.targets = append(m.targets, &sprite{
		SpriteBase: cgame.NewSpriteBase(
			m.g, m.winArena, targetName, targetFrame, LX2X(lx), LY2Y(ly)),
		m: m,
	})
	m.g.SpriteMgr.AddSprite(m.targets[len(m.targets)-1])
}

func (m *myGame) alpha() *sprite {
	for ly := 0; ly < m.lh; ly++ {
		for lx := 0; lx < m.lw; lx++ {
			if m.board[ly][lx] != nil && m.board[ly][lx].Name() == alphaName {
				return m.board[ly][lx]
			}
		}
	}
	panic("where is the alpha?!")
}

func (m *myGame) clearLevel() {
	m.g.SpriteMgr.DeleteAll()
	for ly := 0; ly < m.lh; ly++ {
		for lx := 0; lx < m.lw; lx++ {
			m.board[ly][lx] = nil
		}
	}
	m.targets = m.targets[:0]
	m.steps = 0
	m.lvl = 0
}

func (m *myGame) detectLevelCount() int {
	for lvl := 0; ; lvl++ {
		fn := fmt.Sprintf("resources/level%d.txt", lvl+1)
		if !ios.FileExists(fn) {
			return lvl
		}
	}
}

func (m *myGame) loadLevel(lvl int) {
	m.clearLevel()
	fn := fmt.Sprintf("resources/level%d.txt", lvl+1)
	b, err := ioutil.ReadFile(path.Join(cutil.GetCurFileDir(), fn))
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(b), "\n")
	if len(lines) <= m.lh {
		panic(fmt.Sprintf("%s is corrupted: lines=%d", fn, len(lines)))
	}
	for ly := 0; ly < m.lh; ly++ {
		rline := []rune(lines[ly])
		if len(rline) != m.lw {
			panic(fmt.Sprintf("%s is corrupted: len(line[%d])=%d", fn, ly+1, len(rline)))
		}
		for lx := 0; lx < m.lw; lx++ {
			switch rline[lx] {
			case lvlKeySpace:
			case lvlKeyAlpha:
				m.createAlpha(lx, ly)
			case lvlKeyBrick:
				m.createBrick(lx, ly)
			case lvlKeyConcrete:
				m.createConcrete(lx, ly)
			case lvlKeyTarget:
				m.createTarget(lx, ly)
			default:
				panic(fmt.Sprintf(
					"%s is corrupted: '%c' at col=%d,line=%d", fn, rline[lx], lx+1, ly+1))
			}
		}
	}
	// TODO more sanity check:
	// - there is 1only1 A
	// - # of targets == # of bricks
	m.lvl = lvl
}

func (m *myGame) checkLevelClear() bool {
	for _, t := range m.targets {
		if m.board[t.ly()][t.lx()] == nil || m.board[t.ly()][t.lx()].Name() != brickName {
			return false
		}
	}
	return true
}

var (
	alphaName    = "alpha"
	alphaAtrr    = cwin.Attr{Fg: cterm.ColorLightYellow}
	alpahFrameUp = cgame.FrameFromString(`
M      M
\\    //
 -(◔◔)-
`, alphaAtrr)
	alpahFrameRight = cgame.FrameFromString(`
 (◔◔)__E
  ||‾‾E
 /  \
`, alphaAtrr)
	alpahFrameDown = cgame.FrameFromString(`
 -(◔◔)-
//    \\
W      W
`, alphaAtrr)
	alpahFrameLeft = cgame.FrameFromString(`
E__(◔◔)
 E‾‾||
   /  \
`, alphaAtrr)
	alphaFrames = map[cwin.Dir]cgame.Frame{
		cwin.DirUp:    alpahFrameUp,
		cwin.DirRight: alpahFrameRight,
		cwin.DirDown:  alpahFrameDown,
		cwin.DirLeft:  alpahFrameLeft,
	}

	brickName  = "brick"
	brickAttr  = cwin.Attr{Fg: cterm.ColorWhite, Bg: cterm.ColorRed}
	brickFrame = cgame.FrameFromStringEx(`
⎡      ⎤
⎢      ⎥
⎣      ⎦`, brickAttr, false)

	concreteName  = "concrete"
	concreteAttr  = cwin.Attr{Fg: cterm.ColorWhite, Bg: cterm.ColorDarkGray}
	concreteFrame = cgame.FrameFromStringEx(`
⎡      ⎤
⎢      ⎥
⎣      ⎦`, concreteAttr, false)

	targetName  = "target"
	targetAttr  = cwin.Attr{Fg: cterm.ColorLightGreen, Bg: cterm.ColorBlue}
	targetFrame = cgame.FrameFromStringEx(`
⎡ \  / ⎤
⎢ (><) ⎥
⎣ /  \ ⎦`, targetAttr, false)
)

type sprite struct {
	*cgame.SpriteBase
	m *myGame
}

func (s *sprite) lx() int {
	return X2LX(s.Rect().X)
}

func (s *sprite) ly() int {
	return Y2LY(s.Rect().Y)
}

func (s *sprite) push(dir cwin.Dir) {
	if s.Name() != alphaName {
		panic(fmt.Sprintf("debug: cannot push non-alpha sprite '%s'", s.Name()))
	}
	s.Update(cgame.UpdateArg{F: alphaFrames[dir]})
	dlx, dly := cwin.DirOffSetXY[dir].X, cwin.DirOffSetXY[dir].Y
	curLX, curLY := s.lx(), s.ly()
	newLX := curLX + dlx
	newLY := curLY + dly
	if newLX < 0 || newLX >= s.m.lw || newLY < 0 || newLY >= s.m.lh {
		return
	}
	pushee := s.m.board[newLY][newLX]
	if pushee != nil && pushee.Name() == brickName {
		if !pushee.Update(cgame.UpdateArg{
			DXY: &cwin.Point{X: LX2X(dlx), Y: LY2Y(dly)},
			IBC: cgame.InBoundsCheckFullyVisible,
			CD:  cgame.CollisionDetectionOn}) {
			return
		}
		pushee.brickFrameUpdate()
		pushee.SendToTop()
		s.m.board[newLY+dly][newLX+dlx] = pushee
		s.m.board[newLY][newLX] = nil
	}
	if s.Update(cgame.UpdateArg{DXY: &cwin.Point{X: LX2X(dlx), Y: LY2Y(dly)}}) {
		s.m.g.SoundMgr.PlayMP3(sfxClick, 0, 1)
		s.SendToTop()
		s.m.board[newLY][newLX] = s
		s.m.board[curLY][curLX] = nil
		s.m.steps++
	}
}

func (s *sprite) brickFrameUpdate() {
	overTarget := false
	for _, t := range s.m.targets {
		if t.lx() == s.lx() && t.ly() == s.ly() {
			overTarget = true
		}
	}
	if overTarget {
		s.Update(cgame.UpdateArg{F: cgame.SetAttrInFrame(cgame.CopyFrame(targetFrame), brickAttr)})
	} else {
		s.Update(cgame.UpdateArg{F: brickFrame})
	}
}

func LX2X(lx int) int {
	return lx * (xscale + 1)
}

func LY2Y(ly int) int {
	return ly * (yscale + 1)
}

func X2LX(x int) int {
	return x / (xscale + 1)
}

func Y2LY(y int) int {
	return y / (yscale + 1)
}
