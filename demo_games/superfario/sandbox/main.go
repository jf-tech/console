package main

import (
	"path"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cgame/assets"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cutil"
	"github.com/jf-tech/console/cwin"
	"github.com/jf-tech/console/cwin/ccomp"
)

func main() {
	(&sandbox{}).main()
}

type sandbox struct {
	g               *cgame.Game
	winSandboxFrame cwin.Win
	winSandbox      cwin.Win
	winListBox      *ccomp.ListBox
	winInstr        cwin.Win
}

func (s *sandbox) main() {
	var err error
	s.g, err = cgame.Init(cterm.TCell)
	if err != nil {
		panic(err)
	}
	defer s.g.Close()

	s.setup()

	s.g.Run(assets.GameOverKeys, nil, func(ev cterm.Event) cwin.EventResponse {
		return cwin.EventNotHandled
	})
}

var (
	sandboxW      = 104
	sandboxH      = 39
	sandboxFrameW = sandboxW + 2
	sandboxFrameH = sandboxH + 2
	listboxW      = 25
	listboxH      = 30
	instrW        = listboxW
	instrH        = sandboxFrameH - listboxH
	totalW        = sandboxFrameW + listboxW

	lbIdxFarioRightWalk = 0
	lbStrFarioRightWalk = "Fario Walk →"

	fullpath = func(relpath string) string {
		return path.Join(cutil.GetCurFileDir(), relpath)
	}
	farioRightWalk1Path = fullpath("../resources/fario_right_walk_1.txt")
	farioRightWalk2Path = fullpath("../resources/fario_right_walk_2.txt")
	farioRightWalk3Path = fullpath("../resources/fario_right_walk_3.txt")
	farioRightWalk4Path = fullpath("../resources/fario_right_walk_4.txt")
	farioRightWalkPaths = []string{
		farioRightWalk1Path,
		farioRightWalk2Path,
		farioRightWalk3Path,
		farioRightWalk4Path,
	}
)

func (s *sandbox) setup() {
	sysR := s.g.WinSys.SysWin().Rect()
	s.winSandboxFrame = s.g.WinSys.CreateWin(nil, cwin.WinCfg{
		R: cwin.Rect{
			X: (sysR.W - totalW) / 2,
			Y: (sysR.H - sandboxFrameH) / 2,
			W: sandboxFrameW,
			H: sandboxFrameH},
		Name: "Sandbox",
	})
	s.winSandbox = s.g.WinSys.CreateWin(s.winSandboxFrame, cwin.WinCfg{
		R:        cwin.Rect{X: 0, Y: 0, W: sandboxW, H: sandboxH},
		NoBorder: true,
	})
	s.winListBox = ccomp.CreateListBox(s.g.WinSys, nil, ccomp.ListBoxCfg{
		WinCfg: cwin.WinCfg{
			R: cwin.Rect{
				X: s.winSandboxFrame.Rect().X + s.winSandboxFrame.Rect().W,
				Y: s.winSandboxFrame.Rect().Y,
				W: listboxW,
				H: listboxH,
			},
			Name: "Feature Selection",
		},
		Items: []string{
			lbStrFarioRightWalk,
			"Feature 2",
		},
		EnterKeyToSelect: true,
		OnSelect: func(idx int, selected string) {
			s.clearSandbox()
			switch selected {
			case lbStrFarioRightWalk:
				s.doFarioRightWalk()
			}
		},
	})
	s.winInstr = s.g.WinSys.CreateWin(nil, cwin.WinCfg{
		R: cwin.Rect{
			X: s.winListBox.Rect().X,
			Y: s.winListBox.Rect().Y + s.winListBox.Rect().H,
			W: instrW,
			H: instrH,
		},
		Name: "Help",
	})
	s.g.WinSys.SetFocus(s.winListBox)
	s.winListBox.SetSelected(lbIdxFarioRightWalk)
	s.g.Resume()
}

func (s *sandbox) clearSandbox() {
	s.g.SpriteMgr.DeleteAll()
}

func (s *sandbox) doFarioRightWalk() {
	var fs cgame.Frames
	for _, filepath := range farioRightWalkPaths {
		f, err := cgame.MultiColorFrameFromFile(filepath, 6)
		if err != nil {
			panic(err)
		}
		fs = append(fs, f)
	}
	r := cgame.FrameRect(fs[0])
	sr := s.winSandbox.Rect()
	fario := cgame.NewSpriteBase(s.g, s.winSandbox, "fario", fs[0], (sr.W-r.W)/2, (sr.H-r.H)/2)
	fario.AddAnimator(cgame.NewAnimatorFrame(fario, cgame.AnimatorFrameCfg{
		Frames: cgame.NewSimpleFrameProvider(fs, 200*time.Millisecond, true),
	}))
	s.g.SpriteMgr.AddSprite(fario)
}
