package ccomp

import (
	"fmt"
	"unsafe"

	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/console/cwin"
	"github.com/jf-tech/go-corelib/maths"
)

const (
	ListBoxNoneSelected = -1
)

type ListBoxOnSelect func(idx int, selected string)

type ListBoxCfg struct {
	cwin.WinCfg
	Items            []string
	SelectedAttr     cwin.Attr
	EnterKeyToSelect bool // if true, only Enter/Return key will cause OnSelect to fire.
	OnSelect         ListBoxOnSelect
}

type ListBox struct {
	*cwin.WinBase
	cfg          ListBoxCfg
	items        []string
	firstVisible int
	selected     int
}

func (lb *ListBox) SetItems(items []string) {
	lb.items = items
	lb.firstVisible = 0
	lb.selected = ListBoxNoneSelected
	lb.SetSelected(0)
}

func (lb *ListBox) String() string {
	return fmt.Sprintf("listbox['%s'|0x%X|%s]",
		lb.Cfg().Name, uintptr(unsafe.Pointer(lb)), lb.Rect())
}

func (lb *ListBox) SetSelected(selected int) {
	if selected == lb.selected || selected < 0 || selected >= len(lb.items) {
		return
	}
	lb.FillClient(
		lb.ClientRect().ToOrigin(), cwin.Chx{Ch: cwin.RuneSpace, Attr: lb.Cfg().ClientAttr})
	lb.selected = selected
	if lb.selected < lb.firstVisible {
		lb.firstVisible = lb.selected
	} else if lb.selected >= lb.firstVisible+lb.ClientRect().H {
		lb.firstVisible = lb.selected - lb.ClientRect().H + 1
	}
	for cy := 0; cy < maths.MinInt(len(lb.items)-lb.firstVisible, lb.ClientRect().H); cy++ {
		attr := lb.cfg.ClientAttr
		if cy+lb.firstVisible == lb.selected {
			attr = lb.cfg.SelectedAttr
		}
		lb.SetLine(cy, attr, lb.items[cy+lb.firstVisible])
	}
	if !lb.cfg.EnterKeyToSelect && lb.cfg.OnSelect != nil {
		lb.cfg.OnSelect(lb.selected, lb.items[lb.selected])
	}
}

func (lb *ListBox) GetSelected() (int, string) {
	if lb.selected >= 0 && lb.selected < len(lb.items) {
		return lb.selected, lb.items[lb.selected]
	}
	return ListBoxNoneSelected, ""
}

func (lb *ListBox) moveUp() {
	if lb.selected <= 0 {
		return
	}
	lb.SetSelected(lb.selected - 1)
}

func (lb *ListBox) moveDown() {
	if lb.selected >= len(lb.items)-1 {
		return
	}
	lb.SetSelected(lb.selected + 1)
}

func newListBox(sys *cwin.Sys, parent cwin.Win, cfg ListBoxCfg) *ListBox {
	if cfg.SelectedAttr == cwin.TransparentChx().Attr {
		cfg.SelectedAttr = cwin.Attr{Bg: cterm.ColorBlue}
	}
	lb := &ListBox{WinBase: cwin.NewWinBase(sys, parent, cfg.WinCfg), cfg: cfg}
	lb.SetItems(cfg.Items)
	lb.SetEventHandler(func(ev cterm.Event) cwin.EventResponse {
		if ev.Type != cterm.EventKey ||
			!cwin.FindKey(cwin.Keys(cterm.KeyArrowUp, cterm.KeyArrowDown, cterm.KeyEnter), ev) {
			return cwin.EventNotHandled
		}
		if ev.Key == cterm.KeyArrowUp {
			lb.moveUp()
		} else if ev.Key == cterm.KeyArrowDown {
			lb.moveDown()
		} else {
			// it's cterm.KeyEnter
			if !lb.cfg.EnterKeyToSelect || lb.cfg.OnSelect == nil {
				return cwin.EventNotHandled
			}
			lb.cfg.OnSelect(lb.selected, lb.items[lb.selected])
		}
		return cwin.EventHandled
	})
	return lb
}

func CreateListBox(sys *cwin.Sys, parent cwin.Win, cfg ListBoxCfg) *ListBox {
	return sys.CreateWinEx(parent, func() cwin.Win {
		return newListBox(sys, parent, cfg)
	}).(*ListBox)
}
