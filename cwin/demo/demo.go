package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
)

func main() {
	provider := cterm.TCell
	sys, err := cwin.Init(provider)
	if err != nil {
		panic(err)
	}
	defer sys.Close()

	sysR := sys.GetSysWin().ClientRect()
	demoR := cwin.Rect{X: 0, Y: 0, W: sysR.W * 3 / 4, H: sysR.H * 3 / 4}
	demoR.X = (sysR.W - demoR.W) / 2
	demoR.Y = (sysR.H-demoR.H)/2 - 5
	demoWin := sys.CreateWin(nil, cwin.WinCfg{
		R: demoR,
	})
	var sb strings.Builder
	for i := 0; i < 60; i++ {
		for j := 0; j <= i; j++ {
			sb.WriteString(fmt.Sprintf("%x", j%16))
		}
		sb.WriteRune('\n')
	}
	demoWin.SetText(sb.String())

	demoTitlePrefix := fmt.Sprintf("Demo [%s] (%dx%d)", provider, sysR.W, sysR.H)
	demoWin.SetTitle(fmt.Sprintf("%s - press any key for next", demoTitlePrefix), cwin.AlignLeft)

	fgColorWin := sys.CreateWin(nil, cwin.WinCfg{
		R: cwin.Rect{
			X: demoR.X,
			Y: demoR.Y + demoR.H,
			W: demoR.W,
			H: 3,
		},
		Name: "Foreground Colors",
	})
	for x, color := 0, cterm.ColorBlack; color <= cterm.ColorLightGray; x, color = x+2, color+1 {
		fgColorWin.PutClient(x, 0, cwin.Chx{Ch: '█', Attr: cwin.ChAttr{Fg: color}})
		fgColorWin.PutClient(x+1, 0, cwin.Chx{Ch: '█', Attr: cwin.ChAttr{Fg: color}})
	}
	bgColorWin := sys.CreateWin(nil, cwin.WinCfg{
		R: cwin.Rect{
			X: fgColorWin.Rect().X,
			Y: fgColorWin.Rect().Y + fgColorWin.Rect().H,
			W: fgColorWin.Rect().W,
			H: 3,
		},
		Name: "Background Colors",
	})
	for x, color := 0, cterm.ColorBlack; color <= cterm.ColorLightGray; x, color = x+2, color+1 {
		bgColorWin.PutClient(x, 0, cwin.Chx{Ch: ' ', Attr: cwin.ChAttr{Bg: color}})
		bgColorWin.PutClient(x+1, 0, cwin.Chx{Ch: ' ', Attr: cwin.ChAttr{Bg: color}})
	}

	sys.Update()
	sys.SyncExpectKey(nil)
	demoWin.SetTitle(fmt.Sprintf("%s - MessageBox", demoTitlePrefix), cwin.AlignLeft)

	ret := sys.MessageBox(demoWin,
		"MessageBox",
		`This is a default MessageBox.
It is a modal dialog box.

It can be dismissed by pressing Enter/Return, or ESC.
It returns true if Enter/Return is pressed; false if ESC is pressed.

Current time is: %s`, time.Now().Format(time.RFC3339))
	demoWin.SetText("MessageBox return value: %t", ret)
	demoWin.SetTitle(fmt.Sprintf("%s - Press any key to exit.", demoTitlePrefix), cwin.AlignLeft)
	sys.Update()
	sys.SyncExpectKey(nil)
}
