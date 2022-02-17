package cwin

import (
	"fmt"
	"strings"

	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/go-corelib/maths"
)

type Sys struct {
	sysWin            *Win
	scrBuf, offScrBuf []Chx

	stopEvent chan struct{}
	evChan    chan cterm.Event

	totalChxRendered int64
}

func Init(provider cterm.Provider) (*Sys, error) {
	cterm.SetProvider(provider)
	if err := cterm.Init(); err != nil {
		return nil, err
	}
	w, h := cterm.Size()
	s := &Sys{sysWin: NewWin(nil, WinCfg{R: Rect{0, 0, w, h}, Name: "_root", NoBorder: true})}
	n := w * h
	s.scrBuf = make([]Chx, n)
	s.offScrBuf = make([]Chx, n)
	s.startEventListening()
	return s, nil
}

func (s *Sys) GetSysWin() *Win {
	return s.sysWin
}

func (s *Sys) CreateWin(parent *Win, cfg WinCfg) *Win {
	if parent == nil {
		parent = s.sysWin
	}
	w := NewWin(parent, cfg)
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

// This is a non-blocking call
func (s *Sys) TryGetEvent() cterm.Event {
	if s.evChan == nil {
		panic("startEventListening not called")
	}
	select {
	case ev := <-s.evChan:
		return ev
	default:
		return cterm.Event{Type: cterm.EventNone}
	}
}

// This is a blocking call
func (s *Sys) GetEvent() cterm.Event {
	for {
		ev := s.TryGetEvent()
		if ev.Type != cterm.EventNone {
			return ev
		}
	}
}

// if f == nil, SyncExpectKey waits for any single key and then returns
// if f != nil, SyncExpectKey repeatedly waits for a key & has it processed by f, if f returns false
func (s *Sys) SyncExpectKey(f func(cterm.Key, rune) bool) {
	for {
		ev := s.GetEvent()
		if ev.Type == cterm.EventKey && (f == nil || f(ev.Key, ev.Ch)) {
			break
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
		BorderAttr: ChAttr{Fg: cterm.ColorDefault, Bg: cterm.ColorBlue},
		ClientAttr: ChAttr{Fg: cterm.ColorDefault, Bg: cterm.ColorBlue},
	})
	w.SetTitle(title, AlignCenter)
	w.SetText(msg)
	return w
}

// MessageBoxEx displays a blue-background message window centered in the client area of parent
// window, and synchronously waiting for a set of user specified keys, and return if any of the
// keys are pressed.
func (s *Sys) MessageBoxEx(
	parent *Win, keys []cterm.Event, title, format string, a ...interface{}) cterm.Event {

	w := s.CenterBanner(parent, title, format, a...)
	s.Update()

	ret := cterm.Event{Type: cterm.EventKey}
	s.SyncExpectKey(func(k cterm.Key, ch rune) bool {
		for _, ev := range keys {
			if ch != 0 {
				if ev.Ch == ch {
					ret.Ch = ch
					return true
				}
				continue
			}
			if ev.Key == k {
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
	e := s.MessageBoxEx(parent, Keys(cterm.KeyEnter, cterm.KeyEsc), title, format, a...)
	return e.Key == cterm.KeyEnter
}

func (s *Sys) Update() {
	s.doUpdate(true)
	cterm.Flush()
}

func (s *Sys) TotalChxRendered() int64 {
	return s.totalChxRendered
}

func (s *Sys) DumpTree() string {
	return s.sysWin.DumpTree(0)
}

func (s *Sys) Close() {
	s.stopEventListening()
	cterm.Close()
}

func (s *Sys) doUpdateOffScrBuf(parentSysX, parentSysY int, w *Win, sysRect Rect) {
	if w.hidden {
		return
	}
	// First update the 'w' window content into sysWin's off-screen buffer
	var overlapped bool
	sysRect, overlapped = sysRect.Overlap(
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
	sysRect, overlapped = sysRect.Overlap(clientSysRect)
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
				cterm.SetCell(x, y,
					s.scrBuf[idx].Ch, s.scrBuf[idx].Attr.Fg, s.scrBuf[idx].Attr.Bg)
				s.totalChxRendered++
			}
		}
	}
}

func (s *Sys) startEventListening() {
	if s.stopEvent != nil {
		panic("startEventListening called twice")
	}
	s.stopEvent = make(chan struct{})
	s.evChan = make(chan cterm.Event, 100)

	go func() {
	loop:
		for {
			select {
			case <-s.stopEvent:
				break loop
			default:
				s.evChan <- cterm.PollEvent()
			}
		}
	}()
}

func (s *Sys) stopEventListening() {
	if s.stopEvent == nil {
		return
	}
	close(s.stopEvent)
	s.stopEvent = nil
	// importantly need to call cterm.Interrupt() before closing the evChan because
	// cterm.Interrupt() synchronously waits for cterm.PollEvent finishes so there
	// might be one last event coming through into the evChan. If we close it before
	// calling cterm.Interrupt(), we might get a panic.
	cterm.Interrupt()
	close(s.evChan)
	s.evChan = nil
}
