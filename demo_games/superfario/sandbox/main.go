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
	groundY         int
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

	farioName = "fario"

	lbStrFarioRightWalk = "Fario Walk →"
	lbStrFarioRightJump = "Fario Jump ↗"

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

	farioRightJump1Path = fullpath("../resources/fario_right_jump_1.txt")
	farioRightJump2Path = fullpath("../resources/fario_right_jump_2.txt")
	farioRightJump3Path = fullpath("../resources/fario_right_jump_3.txt")
	farioRightJumpPaths = []string{
		farioRightJump1Path,
		farioRightJump2Path,
		farioRightJump3Path,
	}

	groundName  = "ground"
	groundFrame = cgame.FrameFromStringEx(" ", cwin.Attr{Bg: cterm.ColorDarkGray}, false)
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
			lbStrFarioRightJump,
		},
		EnterKeyToSelect: true,
		OnSelect: func(idx int, selected string) {
			s.clearFario()
			switch selected {
			case lbStrFarioRightWalk:
				s.doFarioRightWalk()
			case lbStrFarioRightJump:
				s.doFarioRightJump()
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
	s.winListBox.SetSelected(0)
	s.g.Resume()

	for y := sandboxH * 3 / 4; y < sandboxH; y++ {
		for x := 0; x < sandboxW; x++ {
			s.g.SpriteMgr.AddSprite(
				cgame.NewSpriteBase(s.g, s.winSandbox, groundName, groundFrame, x, y))
		}
	}
	s.groundY = sandboxH * 3 / 4

	s.g.SpriteMgr.CollidableRegistry().Register(farioName, groundName)
}

func (s *sandbox) clearFario() {
	sprites := s.g.SpriteMgr.Sprites()
	for _, sprite := range sprites {
		if sprite.Name() == farioName {
			s.g.SpriteMgr.DeleteSprite(sprite)
		}
	}
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
	fario := cgame.NewSpriteBase(
		s.g, s.winSandbox, farioName, fs[0], (sr.W-r.W)/2, s.groundY-r.H)
	fario.AddAnimator(cgame.NewAnimatorFrame(fario, cgame.AnimatorFrameCfg{
		Frames: cgame.NewSimpleFrameProvider(fs, 200*time.Millisecond, true),
	}))
	s.g.SpriteMgr.AddSprite(fario)
}

func (s *sandbox) doFarioRightJump() {
	var fs cgame.Frames
	for _, filepath := range farioRightJumpPaths {
		f, err := cgame.MultiColorFrameFromFile(filepath, 6)
		if err != nil {
			panic(err)
		}
		fs = append(fs, f)
	}
	r := cgame.FrameRect(fs[0])
	sr := s.winSandbox.Rect()
	fario := cgame.NewSpriteBase(s.g, s.winSandbox, farioName, fs[0], (sr.W-r.W)/2, s.groundY-r.H)
	fario.AddAnimator(cgame.NewAnimatorFrame(fario, cgame.AnimatorFrameCfg{
		Frames: cgame.NewSimpleFrameProvider(fs, 200*time.Millisecond, false),
		AnimatorCfgCommon: cgame.AnimatorCfgCommon{
			KeepAliveWhenFinished: true,
		},
	}))
	fario.AddAnimator(cgame.NewAnimatorWaypoint(fario, cgame.AnimatorWaypointCfg{
		Waypoints: cgame.NewWaypointProviderAcceleration(cgame.WaypointProviderAccelerationCfg{
			Clock:      s.g.MasterClock,
			InitXSpeed: 10,
			InitYSpeed: -20,
			AccX:       0,
			AccY:       20,
			DeltaT:     time.Millisecond,
		}),
		AnimatorCfgCommon: cgame.AnimatorCfgCommon{
			KeepAliveWhenFinished: true,
			AfterFinish: func() {
				fario.Update(cgame.UpdateArg{
					F:   fs[0],
					IBC: cgame.InBoundsCheckNone,
					CD:  cgame.CollisionDetectionOff,
				})
			},
		},
	}))
	s.g.SpriteMgr.AddSprite(fario)
}
