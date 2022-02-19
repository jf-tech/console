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

	sysR := sys.SysWin().Rect()
	W3_4 := sysR.W * 4 / 5
	H3_4 := sysR.H * 4 / 5

	fgColorWinW, fgColorWinH := W3_4, 3
	bgColorWinW, bgColorWinH := W3_4, 3
	listboxW, listboxH := 20, H3_4-fgColorWinH-bgColorWinH
	demoW, demoH := W3_4-listboxW, listboxH

	demoR := cwin.Rect{X: (sysR.W - W3_4) / 2, Y: (sysR.H - H3_4) / 2, W: demoW, H: demoH}
	demoWin := sys.CreateWin(nil, cwin.WinCfg{R: demoR})
	var sb strings.Builder
	for i := 0; i < 60; i++ {
		for j := 0; j <= i; j++ {
			sb.WriteString(fmt.Sprintf("%x", j%16))
		}
		sb.WriteRune('\n')
	}
	demoWin.SetTextAligned(cwin.AlignRight, sb.String())
	demoTitlePrefix := fmt.Sprintf("Demo [%s] (%dx%d)", provider, sysR.W, sysR.H)

	listboxR := cwin.Rect{X: demoR.X + demoR.W, Y: demoR.Y, W: listboxW, H: listboxH}
	listbox := sys.CreateListBox(nil, cwin.ListBoxCfg{
		WinCfg: cwin.WinCfg{R: listboxR, Name: "ListBox"},
		Items: func() []string {
			var items []string
			for i := 0; i < 100; i++ {
				items = append(items, fmt.Sprintf("Item %d", i))
			}
			return items
		}(),
		OnSelect: func(idx int, selected string) {
			demoWin.SetTitle(fmt.Sprintf(
				"%s - ListBox (%c,%c:change selection) selected: [%d]='%s'. Any other key for next",
				demoTitlePrefix,
				cwin.DirRunes[cwin.DirUp], cwin.DirRunes[cwin.DirDown],
				idx, selected))
		}})
	sys.SetFocus(listbox)

	fgColorWin := sys.CreateWin(nil, cwin.WinCfg{
		R: cwin.Rect{
			X: demoR.X,
			Y: demoR.Y + demoR.H,
			W: fgColorWinW,
			H: fgColorWinH,
		},
		Name: "Foreground Colors",
	})
	for x, color := 0, cterm.ColorBlack; color <= cterm.ColorLightGray; x, color = x+2, color+1 {
		fgColorWin.PutClient(x, 0, cwin.Chx{Ch: '█', Attr: cwin.ChAttr{Fg: color}})
		fgColorWin.PutClient(x+1, 0, cwin.Chx{Ch: '█', Attr: cwin.ChAttr{Fg: color}})
	}
	bgColorWin := sys.CreateWin(nil, cwin.WinCfg{
		R: cwin.Rect{
			X: demoR.X,
			Y: demoR.Y + demoR.H + fgColorWinH,
			W: bgColorWinW,
			H: bgColorWinH,
		},
		Name: "Background Colors",
	})
	for x, color := 0, cterm.ColorBlack; color <= cterm.ColorLightGray; x, color = x+2, color+1 {
		bgColorWin.PutClient(x, 0, cwin.Chx{Ch: ' ', Attr: cwin.ChAttr{Bg: color}})
		bgColorWin.PutClient(x+1, 0, cwin.Chx{Ch: ' ', Attr: cwin.ChAttr{Bg: color}})
	}
	sys.Update()
	sys.Run(func(ev cterm.Event) cwin.EventResponse {
		if ev.Type == cterm.EventKey {
			return cwin.EventLoopStop
		}
		return cwin.EventHandled
	})
	demoWin.SetTitle(fmt.Sprintf("%s - MessageBox", demoTitlePrefix))

	ret := sys.MessageBox(demoWin,
		"MessageBox",
		`This is a MessageBox.
It is a modal dialog box.

It can be dismissed by pressing Enter/Return, or ESC.
It returns true if Enter/Return is pressed; false if ESC is pressed.

Current time is: %s`, time.Now().Format(time.RFC3339))
	demoWin.SetText("MessageBox return value: %t", ret)
	demoWin.SetTitle(fmt.Sprintf("%s - Press any key to exit.", demoTitlePrefix))
	sys.Update()
	sys.SyncExpectKey(nil)
}
