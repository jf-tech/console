package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
)

func main() {
	cterm.SetProvider(cterm.TCell)
	sys, err := cwin.Init()
	if err != nil {
		panic(err)
	}
	defer sys.Close()

	sysR := sys.GetSysWin().ClientRect()
	demoR := cwin.Rect{X: 0, Y: 0, W: sysR.W * 3 / 4, H: sysR.H * 3 / 4}
	demoR.X = (sysR.W - demoR.W) / 2
	demoR.Y = (sysR.H - demoR.H) / 2
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
	demoWin.SetTitle(
		fmt.Sprintf("Demo (%dx%d) - press any key for next", sysR.W, sysR.H), cwin.AlignLeft)
	sys.Update()
	cwin.SyncExpectKey(nil)
	demoWin.SetTitle(
		fmt.Sprintf("Demo (%dx%d) - MessageBox", sysR.W, sysR.H), cwin.AlignLeft)

	sys.MessageBox(demoWin,
		"MessageBox",
		`This is a default MessageBox.
It is a modal dialog box.

It can be dismissed by pressing Enter/Return, or ESC.
It returns true if Enter/Return is pressed; false if ESC is pressed.

Current time is: %s`, time.Now().Format(time.RFC3339))
}
