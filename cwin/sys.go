package cwin

import (
	"fmt"
	"strings"

	"github.com/jf-tech/console/cterm"
	"github.com/jf-tech/go-corelib/maths"
)

type Sys struct {
	winReg map[int64]Win

	sysWin  Win
	inFocus Win

	stopEvent chan struct{}
	evChan    chan cterm.Event

	scrBuf, offScrBuf []Chx
	totalChxRendered  int64
}

func Init(provider cterm.Provider) (*Sys, error) {
	cterm.SetProvider(provider)
	if err := cterm.Init(); err != nil {
		return nil, err
	}
	w, h := cterm.Size()
	n := w * h
	sys := &Sys{winReg: map[int64]Win{}}
	sysWin := NewWinBase(sys, nil, WinCfg{R: Rect{0, 0, w, h}, Name: "_root", NoBorder: true})
	sys.RegWin(sysWin)
	sys.sysWin = sysWin
	sys.scrBuf = make([]Chx, n)
	sys.offScrBuf = make([]Chx, n)
	sys.startEventListening()
	return sys, nil
}

func (s *Sys) Run(fallbackHandler EventHandler) {
	RunEventLoop(s,
		func(ev cterm.Event) EventResponse {
			resp := EventNotHandled
			if s.inFocus != nil && s.inFocus.Cfg().EventHandler != nil {
				resp = s.inFocus.Cfg().EventHandler(ev)
			}
			if resp == EventNotHandled {
				return fallbackHandler(ev)
			}
			return resp
		})
}

func (s *Sys) RegWin(w Win) {
	s.winReg[w.UID()] = w
}

func (s *Sys) SysWin() Win {
	return s.sysWin
}

func (s *Sys) TryFindWin(uid int64) (Win, bool) {
	w, ok := s.winReg[uid]
	return w, ok
}

func (s *Sys) FindWin(uid int64) Win {
	if w, ok := s.TryFindWin(uid); ok {
		return w
	}
	panic(fmt.Sprintf("unable to find Win with UID=%d", uid))
}

func (s *Sys) RemoveWin(w Win) {
	// the sysWin is non-removable
	if w.UID() == s.sysWin.UID() {
		return
	}
	parent := w.Parent()
	prev := w.Prev()
	next := w.Next()
	if prev != nil {
		prev.setNext(next)
	}
	if next != nil {
		next.setPrev(prev)
	}
	if parent != nil {
		if parent.ChildFirst().UID() == w.UID() {
			parent.setChildFirst(next)
		}
		if parent.ChildLast().UID() == w.UID() {
			parent.setChildLast(prev)
		}
	}
	if s.inFocus != nil && s.inFocus.UID() == w.UID() {
		s.inFocus = nil
	}
	delete(s.winReg, w.UID())
}

func (s *Sys) SetFocus(w Win) {
	if s.inFocus != nil && s.inFocus.UID() == w.UID() {
		return
	}
	w.SendToTop(true)
	s.inFocus = s.FindWin(w.UID())
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
	RunEventLoop(s, func(ev cterm.Event) EventResponse {
		if ev.Type == cterm.EventKey && (f == nil || f(ev.Key, ev.Ch)) {
			return EventLoopStop
		}
		return EventHandled
	})
}

func (s *Sys) CenterBanner(parent Win, title, format string, a ...interface{}) Win {
	if parent == nil {
		parent = s.sysWin
	}
	maxLineLen := 0
	msg := fmt.Sprintf(format, a...)
	lines := strings.Split(msg, "\n")
	for y := 0; y < maths.MinInt(parent.ClientRect().H-2, len(lines)); y++ {
		rline := []rune(lines[y])
		maxLineLen = maths.MaxInt(maxLineLen, len(rline))
	}
	width := maths.MinInt(maxLineLen+4, parent.ClientRect().W)  // 2 border lines + 2 padding spaces
	height := maths.MinInt(len(lines)+2, parent.ClientRect().H) // 2 border lines
	w := s.CreateWin(parent, WinCfg{
		R: Rect{
			X: (parent.ClientRect().W - width) / 2,
			Y: (parent.ClientRect().H - height) / 2,
			W: width,
			H: height},
		Name:       title,
		BorderAttr: ChAttr{Fg: cterm.ColorDefault, Bg: cterm.ColorBlue},
		ClientAttr: ChAttr{Fg: cterm.ColorDefault, Bg: cterm.ColorBlue},
	})
	w.SetTitleAligned(AlignCenter, title)
	w.SetText(msg)
	return w
}

// MessageBoxEx displays a blue-background message window centered in the client area of parent
// window, and synchronously waiting for a set of user specified keys, and return if any of the
// keys are pressed.
func (s *Sys) MessageBoxEx(
	parent Win, keys []cterm.Event, title, format string, a ...interface{}) cterm.Event {

	w := s.CenterBanner(parent, title, format, a...)
	s.SetFocus(w)
	s.Update()
	var ret cterm.Event
	RunEventLoop(s,
		func(ev cterm.Event) EventResponse {
			if FindKey(keys, ev) {
				ret = ev
				return EventLoopStop
			}
			return EventHandled
		})
	s.RemoveWin(w)
	s.Update()
	return ret
}

// MessageBox is mostly similar to MessageBoxEx but only with 2 expected keys: Enter or ESC
// It returns true if Enter is pressed or false if ESC is pressed.
func (s *Sys) MessageBox(parent Win, title, format string, a ...interface{}) bool {
	e := s.MessageBoxEx(parent, Keys(cterm.KeyEnter, cterm.KeyEsc), title, format, a...)
	return e.Key == cterm.KeyEnter
}

func (s *Sys) Update() {
	s.doUpdate(true)
	cterm.Flush()
}

func (s *Sys) Refresh() {
	s.doUpdate(false)
	cterm.Sync()
}

func (s *Sys) TotalChxRendered() int64 {
	return s.totalChxRendered
}

func (s *Sys) Close() {
	s.stopEventListening()
	cterm.Close()
}

func (s *Sys) doUpdateOffScrBuf(parentSysX, parentSysY int, w Win, sysRect Rect) {
	// First update the 'w' window content into sysWin's off-screen buffer
	var overlapped bool
	sysRect, overlapped = sysRect.Overlap(
		Rect{parentSysX + w.Cfg().R.X, parentSysY + w.Cfg().R.Y, w.Cfg().R.W, w.Cfg().R.H})
	if !overlapped {
		return
	}
	for y := 0; y < sysRect.H; y++ {
		for x := 0; x < sysRect.W; x++ {
			chx := w.Get(sysRect.X-parentSysX-w.Cfg().R.X+x, sysRect.Y-parentSysY-w.Cfg().R.Y+y)
			if chx == chxTransparent {
				continue
			}
			sysOffScrBufIdx := s.sysWin.(*WinBase).bufIdx(sysRect.X+x, sysRect.Y+y)
			s.offScrBuf[sysOffScrBufIdx] = chx
		}
	}
	// Then update the 'w' window's child windows into sysWin's off-screen buffer
	clientSysRect := Rect{
		parentSysX + w.Cfg().R.X + w.ClientRect().X,
		parentSysY + w.Cfg().R.Y + w.ClientRect().Y,
		w.ClientRect().W,
		w.ClientRect().H}
	sysRect, overlapped = sysRect.Overlap(clientSysRect)
	if !overlapped {
		return
	}
	for child := w.ChildFirst(); child != nil; child = child.Next() {
		s.doUpdateOffScrBuf(clientSysRect.X, clientSysRect.Y, child, sysRect)
	}
}

func (s *Sys) doUpdate(differential bool) {
	s.doUpdateOffScrBuf(
		0, 0, s.sysWin, Rect{0, 0, s.sysWin.Cfg().R.W, s.sysWin.Cfg().R.H})
	for y := 0; y < s.sysWin.Cfg().R.H; y++ {
		for x := 0; x < s.sysWin.Cfg().R.W; x++ {
			idx := s.sysWin.(*WinBase).bufIdx(x, y)
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
