package cwin

import (
	"fmt"
	"strings"

	"github.com/jf-tech/go-corelib/maths"
	"github.com/nsf/termbox-go"
)

type Sys struct {
	sysWin            *Win
	scrBuf, offScrBuf []Chx
}

func (s *Sys) GetSysWin() *Win {
	return s.sysWin
}

func (s *Sys) CreateWin(parent *Win, cfg WinCfg) *Win {
	if parent == nil {
		parent = s.sysWin
	}
	w := newWin(parent, cfg)
	if parent.childn == nil {
		parent.child1 = w
		parent.childn = w
	} else {
		w.prev = parent.childn
		parent.childn.next = w
		parent.childn = w
	}
	return w
}

func (s *Sys) RemoveWin(w *Win) {
	if w == s.sysWin {
		return
	}
	parent := w.parent
	prev := w.prev
	next := w.next
	if prev != nil {
		prev.next = next
	}
	if next != nil {
		next.prev = prev
	}
	if parent != nil {
		if parent.child1 == w {
			parent.child1 = next
		}
		if parent.childn == w {
			parent.childn = prev
		}
	}
}

func (s *Sys) CenterBanner(parent *Win, title, format string, a ...interface{}) *Win {
	if parent == nil {
		parent = s.sysWin
	}
	maxLineLen := 0
	msg := fmt.Sprintf(format, a...)
	lines := strings.Split(msg, "\n")
	for y := 0; y < maths.MinInt(parent.clientR.H-2, len(lines)); y++ {
		rline := []rune(lines[y])
		maxLineLen = maths.MaxInt(maxLineLen, len(rline))
	}
	width := maths.MinInt(maxLineLen+4, parent.clientR.W)  // 2 border lines + 2 padding spaces
	height := maths.MinInt(len(lines)+2, parent.clientR.H) // 2 border lines
	w := s.CreateWin(parent, WinCfg{
		R: Rect{
			X: (parent.clientR.W - width) / 2,
			Y: (parent.clientR.H - height) / 2,
			W: width,
			H: height},
		Name:       title,
		BorderAttr: ChAttr{Fg: termbox.ColorDefault, Bg: termbox.ColorBlue},
		ClientAttr: ChAttr{Fg: termbox.ColorDefault, Bg: termbox.ColorBlue},
	})
	w.SetTitle(title, AlignCenter)
	w.SetText(msg)
	return w
}

// MessageBoxEx displays a blue-background message window centered in the client area of parent
// window, and synchronously waiting for a set of user specified keys, and return if any of the
// keys are pressed.
func (s *Sys) MessageBoxEx(
	parent *Win, keys []termbox.Event, title, format string, a ...interface{}) termbox.Event {

	w := s.CenterBanner(parent, title, format, a...)
	s.Update()

	ret := termbox.Event{Type: termbox.EventKey}
	SyncExpectKey(func(k termbox.Key, ch rune) bool {
		for _, e := range keys {
			if ch != 0 {
				if e.Ch == ch {
					ret.Ch = ch
					return true
				}
				continue
			}
			if e.Key == k {
				ret.Key = k
				return true
			}
		}
		return false
	})
	s.RemoveWin(w)
	s.Update()
	return ret
}

// MessageBox is mostly similar to MessageBoxEx but only with 2 expected keys: Enter or ESC
// It returns true if Enter is pressed or false if ESC is pressed.
func (s *Sys) MessageBox(parent *Win, title, format string, a ...interface{}) bool {
	e := s.MessageBoxEx(parent,
		[]termbox.Event{{Key: termbox.KeyEnter}, {Key: termbox.KeyEsc}},
		title, format, a...)
	return e.Key == termbox.KeyEnter
}

func (s *Sys) doUpdateOffScrBuf(parentSysX, parentSysY int, w *Win, sysRect Rect) {
	if w.hidden {
		return
	}
	// First update the 'w' window content into sysWin's off-screen buffer
	var overlapped bool
	overlapped, sysRect = sysRect.Overlap(
		Rect{parentSysX + w.cfg.R.X, parentSysY + w.cfg.R.Y, w.cfg.R.W, w.cfg.R.H})
	if !overlapped {
		return
	}
	for y := 0; y < sysRect.H; y++ {
		for x := 0; x < sysRect.W; x++ {
			winBufIdx := w.bufIdx(
				sysRect.X-parentSysX-w.cfg.R.X+x, sysRect.Y-parentSysY-w.cfg.R.Y+y)
			if w.buf[winBufIdx] == chxTransparent {
				continue
			}
			sysOffScrBufIdx := s.sysWin.bufIdx(sysRect.X+x, sysRect.Y+y)
			s.offScrBuf[sysOffScrBufIdx] = w.buf[winBufIdx]
		}
	}
	// Then update the 'w' window's child windows into sysWin's off-screen buffer
	clientSysRect := Rect{
		parentSysX + w.cfg.R.X + w.clientR.X,
		parentSysY + w.cfg.R.Y + w.clientR.Y,
		w.clientR.W,
		w.clientR.H}
	overlapped, sysRect = sysRect.Overlap(clientSysRect)
	if !overlapped {
		return
	}
	for child := w.child1; child != nil; child = child.next {
		s.doUpdateOffScrBuf(clientSysRect.X, clientSysRect.Y, child, sysRect)
	}
}

func (s *Sys) doUpdate(differential bool) {
	s.doUpdateOffScrBuf(
		0, 0, s.sysWin, Rect{0, 0, s.sysWin.cfg.R.W, s.sysWin.cfg.R.H})
	for y := 0; y < s.sysWin.cfg.R.H; y++ {
		for x := 0; x < s.sysWin.cfg.R.W; x++ {
			idx := s.sysWin.bufIdx(x, y)
			if !differential || s.scrBuf[idx] != s.offScrBuf[idx] {
				s.scrBuf[idx] = s.offScrBuf[idx]
				termbox.SetCell(x, y,
					s.scrBuf[idx].Ch, s.scrBuf[idx].Attr.Fg, s.scrBuf[idx].Attr.Bg)
			}
		}
	}
}

func (s *Sys) Update() {
	s.doUpdate(false)
	termbox.Flush()
}

func (s *Sys) DumpTree() string {
	return s.sysWin.DumpTree(0)
}

func (s *Sys) Close() {
	termbox.Close()
}

func Init() (*Sys, error) {
	if err := termbox.Init(); err != nil {
		return nil, err
	}
	w, h := termbox.Size()
	s := &Sys{sysWin: newWin(nil, WinCfg{R: Rect{0, 0, w, h}, Name: "_root", NoBorder: true})}
	n := w * h
	s.scrBuf = make([]Chx, n)
	s.offScrBuf = make([]Chx, n)
	return s, nil
}
