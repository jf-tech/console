package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cutil"
	"github.com/jf-tech/console/cwin"
	"github.com/jf-tech/console/cwin/ccomp"
)

var gravityString = []string{
	"0 char/s^2",
	"10 char/s^2",
	"20 char/s^2",
	"40 char/s^2",
	"100 char/s^2",
	"150 char/s^2",
}
var gravityVals = []cgame.CharPerSecSec{
	0,
	10,
	20,
	40,
	100,
	150,
}
var curGravityIdx = 2

var bounce = true

func main() {
	g, err := cgame.Init(cterm.TCell)
	if err != nil {
		panic(err)
	}
	defer g.Close()
	sysWinR := g.WinSys.SysWin().Rect()

	gWinW := 41
	gWinH := 10
	gWin := ccomp.CreateListBox(g.WinSys, nil, ccomp.ListBoxCfg{
		WinCfg: cwin.WinCfg{
			R:    cwin.Rect{X: sysWinR.W - gWinW, Y: 0, W: gWinW, H: gWinH},
			Name: "Choose a gravity value",
		},
		Items: gravityString,
		OnSelect: func(idx int, selected string) {
			curGravityIdx = idx
			g.SpriteMgr.DeleteAll()
		},
	})
	curGravityIdx = 3
	gWin.SetSelected(curGravityIdx)
	g.WinSys.SetFocus(gWin)

	statsWinW := gWinW
	statsWinH := sysWinR.H - gWinH
	statsWinR := cwin.Rect{X: sysWinR.W - statsWinW, Y: gWinH, W: statsWinW, H: statsWinH}
	statsWin := g.WinSys.CreateWin(nil, cwin.WinCfg{R: statsWinR, Name: "Stats"})

	demoWinR := cwin.Rect{X: 0, Y: 0, W: sysWinR.W - statsWinR.W, H: sysWinR.H}
	demoWin := g.WinSys.CreateWin(nil, cwin.WinCfg{
		R:    demoWinR,
		Name: "Demo - ↑↓ to change gravity; 'b' to turn bounce on/off; space to pause/resume. Any other key to exit."})

	g.WinSys.Update() // nothing shows onto screen unless Update() is called.
	g.Resume()        // game (master clock) is always paused right after init.

	doDemo(g, demoWin, statsWin)
}

type spriteParticle struct {
	*cgame.SpriteBase
	id       int64
	demoWinR cwin.Rect
}

func (s *spriteParticle) createGravityAnimator(vx, vy cgame.CharPerSec) {
	s.AddAnimator(cgame.NewAnimatorWaypoint(s, cgame.AnimatorWaypointCfg{
		Waypoints: cgame.NewWaypointProviderAcceleration(cgame.WaypointProviderAccelerationCfg{
			Clock:      s.Game().MasterClock,
			InitXSpeed: vx,
			InitYSpeed: vy,
			AccX:       0,
			AccY:       gravityVals[curGravityIdx],
			DeltaT:     time.Millisecond,
		}),
		AnimatorCfgCommon: cgame.AnimatorCfgCommon{
			// We don't want the particle to be destroy
			// the monent it flies out of bound - we only want to kill it when its x coord
			// outside the demoWin X range; for y, we do hope to see those going up will
			// eventually (and hopefully :) coming down. Thus, we need to turn off the
			// automatic bounds check.
			InBoundsCheckType: cgame.InBoundsCheckNone,
			// Turn off collision check, although strictly speaking this is not necessary
			// since we didn't register anything in the CollidableRegistry so nothing
			// will collide with each other anyway.
			CollisionDetectionType: cgame.CollisionDetectionOff,
			AfterUpdate: func() {
				if s.IsDestroyed() {
					return
				}
				if s.demoWinR.Contain(s.Rect().X, s.Rect().Y) {
					return
				}
				wp := s.Animators()[0].(*cgame.AnimatorWaypoint).
					Cfg().Waypoints.(*cgame.WaypointProviderAcceleration)
				if wp.Cfg().AccY != 0 &&
					s.Rect().X >= 0 && s.Rect().X < s.demoWinR.W && s.Rect().Y < s.demoWinR.H {
					return
				}
				vx, vy := wp.CurSpeed()
				if bounce &&
					s.Rect().X >= 0 && s.Rect().X < s.demoWinR.W &&
					s.Rect().Y >= s.demoWinR.H && vy > 0 {
					s.DeleteAnimator(s.Animators()[0])
					s.Update(cgame.UpdateArg{
						DXY: &cwin.Point{Y: s.demoWinR.H - s.Rect().Y - 1},
						IBC: cgame.InBoundsCheckNone,
						CD:  cgame.CollisionDetectionOff,
					})
					s.createGravityAnimator(vx, -vy)
					return
				}
				s.Mgr().DeleteSprite(s)
			},
		},
	}))
}

var (
	particleImg         = "◯"
	particleName        = "particle"
	particleFrameNoAttr = cgame.FrameFromString(particleImg, cwin.Attr{})
)

func genParticleColor() cterm.Attribute {
	min := int(cterm.ColorRed)
	max := int(cterm.ColorLightGray)
	return cterm.Attribute(rand.Int()%(max-min+1) + min)
}

func genParticleXSpeed() cgame.CharPerSec {
	for {
		vx := cgame.CharPerSec(rand.Int()%61 - 30) // [-30, -2] and [2, 30]
		if vx != 0 && vx != -1 && vx != 1 {
			return vx
		}
	}
}

func genParticleYSpeed() cgame.CharPerSec {
	return cgame.CharPerSec(rand.Int()%56 - 60) // [-60,-5]
}

func doDemo(g *cgame.Game, demoWin, statsWin cwin.Win) {
	r := demoWin.ClientRect().ToOrigin()
	createParticle := func(x, y int, vx, vy cgame.CharPerSec, color cterm.Attribute) {
		attr := cwin.Attr{Fg: color}
		p := &spriteParticle{
			SpriteBase: cgame.NewSpriteBase(g, demoWin, particleName,
				cgame.SetAttrInFrame(cgame.CopyFrame(particleFrameNoAttr), attr), x, y),
			id:       cwin.GenUID(),
			demoWinR: r,
		}
		p.createGravityAnimator(vx, vy)
		g.SpriteMgr.AddSprite(p)
	}

	showStats := func() {
		var sb strings.Builder
		sb.WriteString(fmt.Sprint("Stats:\n"))
		sb.WriteString(fmt.Sprintf("- Time: %s\n", g.MasterClock.Now().Round(time.Millisecond)))
		sb.WriteString(fmt.Sprintf("- Bounce: %t\n", bounce))
		sb.WriteString(fmt.Sprintf("- Particles: %d\n", len(g.SpriteMgr.Sprites())))
		sb.WriteString(fmt.Sprintf("- FPS: %d\n", g.WinSys.FPS()))
		sb.WriteString(fmt.Sprintf("- Mem: %s\n", cwin.ByteSizeStr(g.HeapUsageInBytes())))
		sb.WriteString(fmt.Sprintf("- Pixels rendered: %s\n",
			cwin.ByteSizeStr(g.WinSys.TotalChxRendered())))
		sb.WriteString(fmt.Sprint("\n"))
		sb.WriteString(fmt.Sprint("Particles:\n"))
		for _, s := range g.SpriteMgr.Sprites() {
			sp := s.(*spriteParticle)
			wpAcc := sp.Animators()[0].(*cgame.AnimatorWaypoint).Cfg().Waypoints.(*cgame.WaypointProviderAcceleration)
			vx, vy := wpAcc.CurSpeed()
			sb.WriteString(fmt.Sprintf(
				"- [%4d]: x/y=%3d/%3d, vx/vy=%3.0f/%3.0f\n",
				sp.id, sp.Rect().X, sp.Rect().Y, vx, vy))
		}
		statsWin.SetText(sb.String())
	}

	prob := cutil.NewPeriodicProbabilityChecker("10%", 100*time.Millisecond)
	prob.Reset(g.MasterClock)
	g.Run(nil, cwin.Keys(' '), func(ev cterm.Event) cwin.EventResponse {
		showStats()
		if prob.Check() {
			x, y := r.W/2, r.H-1
			xspeed := genParticleXSpeed()
			yspeed := genParticleYSpeed()
			color := genParticleColor()
			createParticle(x, y, xspeed, yspeed, color)
		}
		if ev.Type != cterm.EventKey {
			return cwin.EventNotHandled
		}
		if ev.Ch == 'b' {
			bounce = !bounce
			g.SpriteMgr.DeleteAll()
			return cwin.EventHandled
		}
		return cwin.EventLoopStop
	})
}