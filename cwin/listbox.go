package cwin

import (
	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/go-corelib/maths"
)

const (
	ListBoxNoneSelected = -1
)

type ListBoxOnSelect func(idx int, selected string)

type ListBoxCfg struct {
	WinCfg
	Items        []string
	SelectedAttr ChAttr
	Align        Align
	OnSelect     ListBoxOnSelect
}

type ListBox struct {
	*Win
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

func (lb *ListBox) SetSelected(selected int) {
	if selected == lb.selected || selected < 0 || selected >= len(lb.items) {
		return
	}
	lb.FillClient(lb.clientR.ToOrigin(), Chx{Ch: RuneSpace, Attr: lb.cfg.ClientAttr})
	lb.selected = selected
	if lb.selected < lb.firstVisible {
		lb.firstVisible = lb.selected
	} else if lb.selected >= lb.firstVisible+lb.clientR.H {
		lb.firstVisible = lb.selected - lb.clientR.H + 1
	}
	for cy := 0; cy < maths.MinInt(len(lb.items)-lb.firstVisible, lb.clientR.H); cy++ {
		attr := lb.cfg.ClientAttr
		if cy+lb.firstVisible == lb.selected {
			attr = lb.cfg.SelectedAttr
		}
		lb.setTextLine(cy, lb.items[cy+lb.firstVisible], lb.cfg.Align, attr)
	}
	if lb.cfg.OnSelect != nil {
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

func newListBox(sys *Sys, parent *Win, cfg ListBoxCfg) *ListBox {
	if cfg.SelectedAttr == TransparentChx().Attr {
		cfg.SelectedAttr = ChAttr{Bg: cterm.ColorBlue}
	}
	lb := &ListBox{Win: newWin(sys, parent, cfg.WinCfg), cfg: cfg}
	lb.SetItems(cfg.Items)
	lb.SetEventHandler(func(ev cterm.Event) EventResponse {
		if ev.Type != cterm.EventKey ||
			!FindKey(Keys(cterm.KeyArrowUp, cterm.KeyArrowDown), ev) {
			return EventNotHandled
		}
		if ev.Key == cterm.KeyArrowUp {
			lb.moveUp()
		} else {
			lb.moveDown()
		}
		return EventHandled
	})
	return lb
}

func (s *Sys) CreateListBox(parent *Win, cfg ListBoxCfg) *ListBox {
	if parent == nil {
		parent = s.sysWin
	}
	lb := newListBox(s, parent, cfg)
	// TODO code dup as in Sys.CreateWin
	if parent.childn == nil {
		parent.child1 = lb.Win
		parent.childn = lb.Win
	} else {
		lb.prev = parent.childn
		parent.childn.next = lb.Win
		parent.childn = lb.Win
	}
	return lb
}
