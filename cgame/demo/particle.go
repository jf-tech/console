package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/jf-tech/console/cgame"
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
)

func main() {
	g, err := cgame.Init(cterm.TCell)
	if err != nil {
		panic(err)
	}
	defer g.Close()
	sysWinR := g.WinSys.GetSysWin().Rect()

	debugWinW := 60
	debugWinR := cwin.Rect{X: sysWinR.W - debugWinW, Y: 0, W: debugWinW, H: sysWinR.H}
	debugWin := g.WinSys.CreateWin(nil, cwin.WinCfg{
		R:    debugWinR,
		Name: "Debug"})

	demoWinR := cwin.Rect{X: 0, Y: 0, W: sysWinR.W - debugWinR.W, H: sysWinR.H}
	demoWin := g.WinSys.CreateWin(nil, cwin.WinCfg{
		R:    demoWinR,
		Name: "Demo - ↑ to add particle; ↓ to remove; space to pause. Any other key to exit."})

	g.WinSys.Update() // nothing shows onto screen unless Update() is called.
	g.Resume()        // game (master clock) is always paused right after init.

	doDemo(g, demoWin, debugWin)
}

type spriteParticle struct {
	*cgame.SpriteBase
	dx, dy         int
	speed          cgame.CharPerSec
	animator       cgame.Animator
	collisionCount func()
	hitBoundsCount func()
}

var (
	particleImg         = "◯"
	particleName        = "particle"
	particleFrameNoAttr = cgame.FrameFromString(particleImg, cwin.ChAttr{})
)

// implements cgame.WaypointProvider
func (sp *spriteParticle) Next() (cgame.Waypoint, bool) {
	// just random number that's big enough to make sure the particle
	// will keep on going until either hits something or hits the boundary
	dist := 1000
	dx := sp.dx * dist
	dy := sp.dy * dist
	wp := cgame.Waypoint{
		DX: dx,
		DY: dy,
		T: time.Duration(
			math.Sqrt(float64(dx*dx+dy*dy)) / float64(sp.speed) * float64(time.Second)),
	}
	return wp, true
}

// implements cgame.InBoundsCheckResponse
func (sp *spriteParticle) InBoundsCheckNotify(result cgame.InBoundsCheckResult) cgame.InBoundsCheckResponseType {
	switch result {
	case cgame.InBoundsCheckResultN, cgame.InBoundsCheckResultS:
		sp.dy = -sp.dy
		sp.resetAnimator()
		sp.hitBoundsCount()
		return cgame.InBoundsCheckResponseAbandon
	case cgame.InBoundsCheckResultE, cgame.InBoundsCheckResultW:
		sp.dx = -sp.dx
		sp.resetAnimator()
		sp.hitBoundsCount()
		return cgame.InBoundsCheckResponseAbandon
	}
	return cgame.InBoundsCheckResponseJustDoIt
}

// implements cgame.CollisionResponse
func (sp *spriteParticle) CollisionNotify(
	initiator bool, collidedWith []cgame.Sprite) cgame.CollisionResponseType {
	if !initiator {
		// let the collision initiating particle to handle speed exchange.
		// this non collision initiator's response doesn't matter, always ignored.
		return cgame.CollisionResponseAbandon
	}
	dx, dy := sp.dx, sp.dy
	speed := sp.speed
	other := collidedWith[0].(*spriteParticle)
	sp.dx, sp.dy = other.dx, other.dy
	sp.speed = other.speed
	other.dx, other.dy = dx, dy
	other.speed = speed
	sp.resetAnimator()
	other.resetAnimator()
	sp.collisionCount()
	return cgame.CollisionResponseAbandon
}

func (sp *spriteParticle) resetAnimator() {
	if sp.animator != nil {
		sp.DeleteAnimator(sp.animator)
	}
	aw := cgame.NewAnimatorWaypoint(sp.SpriteBase, cgame.AnimatorWaypointCfg{
		Waypoints: sp,
		AnimatorCfgCommon: cgame.AnimatorCfgCommon{
			InBoundsCheckType:      cgame.InBoundsCheckFullyVisible,
			CollisionDetectionType: cgame.CollisionDetectionOn,
			KeepAliveWhenFinished:  true,
		},
	})
	sp.animator = aw
	sp.AddAnimator(aw)
}

func genParticleColor() cterm.Attribute {
	min := int(cterm.ColorRed)
	max := int(cterm.ColorLightGray)
	return cterm.Attribute(rand.Int()%(max-min+1) + min)
}

func genParticleDXY() (int, int) {
	for {
		dx := rand.Int()%21 - 10 // [-10, 10]
		dy := rand.Int()%21 - 10 // [-10, 10]
		if dx != 0 || dy != 0 {
			return dx, dy
		}
	}
}

func genParticleSpeed() cgame.CharPerSec {
	return cgame.CharPerSec(rand.Int()%36 + 5) // [5,40]
}

func doDemo(g *cgame.Game, demoWin, debugWin *cwin.Win) {
	g.SpriteMgr.CollidableRegistry().Register(particleName, particleName)
	r := demoWin.ClientRect().ToOrigin()
	collision := int64(0)
	hitBounds := int64(0)
	var ids []int64
	createParticle := func(x, y, dx, dy int, color cterm.Attribute, speed cgame.CharPerSec) bool {
		attr := cwin.ChAttr{Fg: color}
		s := &spriteParticle{
			SpriteBase: cgame.NewSpriteBase(g, demoWin, particleName,
				cgame.SetAttrInFrame(cgame.CopyFrame(particleFrameNoAttr), attr), x, y),
			dx:    dx,
			dy:    dy,
			speed: speed,
			collisionCount: func() {
				collision++
			},
			hitBoundsCount: func() {
				hitBounds++
			}}
		// A little trick to do a collision test before the sprite is even added to the system
		// by Update with no dx/dy and no new frame, simmply with a bounds check and collision
		// check. If the update fails, then there is a collision (or the sprite is out of bounds)
		// circle though)
		f := s.Frame()
		if !s.Update(cgame.UpdateArg{
			F:   f,
			IBC: cgame.InBoundsCheckFullyVisible,
			CD:  cgame.CollisionDetectionOn}) {
			s.Destroy() // do remember to destroy the sprite as its cwin is already created.
			return false
		}
		s.resetAnimator()
		g.SpriteMgr.AsyncCreateSprite(s)
		ids = append(ids, s.UID())
		return true
	}

	// classic Newton's cradle
	createParticle(r.W/2-30, r.H/2, 1, 0, cterm.ColorLightYellow, 40)
	createParticle(r.W/2, r.H/2, 0, 0, cterm.ColorLightCyan, 0)
	createParticle(r.W/2+1, r.H/2, 0, 0, cterm.ColorLightGreen, 0)
	createParticle(r.W/2+2, r.H/2, 0, 0, cterm.ColorLightBlue, 0)
	createParticle(r.W/2+3, r.H/2, 0, 0, cterm.ColorWhite, 0)

	dc := cgame.NewStopwatch(g.MasterClock)

	showDebugInfo := func() {
		var sb strings.Builder
		sb.WriteString(fmt.Sprint("Stats:\n"))
		sb.WriteString(fmt.Sprintf("- Time: %s\n", g.MasterClock.Now().Round(time.Millisecond)))
		sb.WriteString(fmt.Sprintf("- FPS: %.0f\n", g.FPS()))
		sb.WriteString(fmt.Sprintf("- Mem: %s\n", cwin.ByteSizeStr(g.HeapUsageInBytes())))
		sb.WriteString(fmt.Sprintf("- Pixels: %s\n", cwin.ByteSizeStr(g.WinSys.TotalChxRendered())))
		sb.WriteString(fmt.Sprintf("- Loop time: %s\n", dc.Total()))
		sb.WriteString(fmt.Sprintf("- Particle #: %d\n", len(ids)))
		sb.WriteString(fmt.Sprintf("- Collisions: %d\n", collision))
		sb.WriteString(fmt.Sprintf("- Boundary Hits: %d\n", hitBounds))
		sb.WriteString(fmt.Sprint("\n"))
		sb.WriteString(fmt.Sprint("Particles:\n"))
		for _, s := range g.SpriteMgr.Sprites() {
			sp := s.(*spriteParticle)
			sb.WriteString(fmt.Sprintf(
				"- id(%2d): x/y=%3d/%3d, dx/dy=%3d/%3d, speed=%2.1f\n",
				sp.UID(), sp.Rect().X, sp.Rect().Y, sp.dx, sp.dy, sp.speed))
		}
		debugWin.SetText(sb.String())
		dc.Reset()
	}

	g.Run(cwin.Keys(cterm.KeyEsc, 'q'), cwin.Keys(' '), func(ev cterm.Event) bool {
		showDebugInfo()
		dc.Start()
		defer dc.Stop()
		if ev.Type == cterm.EventKey {
			if ev.Key == cterm.KeyArrowUp {
				for {
					x, y := rand.Int()%r.W, rand.Int()%r.H
					dx, dy := genParticleDXY()
					color := genParticleColor()
					speed := genParticleSpeed()
					// It's possible the randonly generated particle collides with others
					// upon creation. So keep trying until it's not.
					if createParticle(x, y, dx, dy, color, speed) {
						break
					}
				}
				return false
			}
			if ev.Key == cterm.KeyArrowDown {
				if len(ids) > 0 {
					idx := rand.Int() % len(ids)
					id := ids[idx]
					s := g.SpriteMgr.FindByUID(id)
					g.SpriteMgr.AsyncDeleteSprite(s)
					copy(ids[idx:], ids[idx+1:])
					ids = ids[:len(ids)-1]
				}
				return false
			}
			return true
		}
		return false
	})
}
