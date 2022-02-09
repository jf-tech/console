package cgame

import (
	"bytes"
	"io"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/jf-tech/console/cwin"
)

type SoundID int64

const (
	InvalidSoundID = SoundID(-1)

	targetSampleRate  = beep.SampleRate(44100)
	resamplingQuality = 4
)

var inited = false

type SoundManager struct {
	ctrls  sync.Map // [SoundID]*beep.Ctrl
	files  sync.Map // [filepath]SoundID
	paused bool

	avoidSameClipConcurrentPlaying bool
}

func newSoundManager() *SoundManager {
	return &SoundManager{}
}

func (sm *SoundManager) Init() error {
	if inited {
		return nil
	}
	err := speaker.Init(targetSampleRate, targetSampleRate.N(time.Second/10))
	if err != nil {
		return err
	}
	inited = true
	return nil
}

func (sm *SoundManager) Close() {
	speaker.Clear()
	sm.ctrls = sync.Map{}
	sm.files = sync.Map{}
}

func (sm *SoundManager) AvoidSameClipConcurrentPlaying() {
	sm.avoidSameClipConcurrentPlaying = true
}

func (sm *SoundManager) PauseAll() {
	if sm.paused {
		return
	}
	speaker.Lock()
	sm.ctrls.Range(func(_, ctrl interface{}) bool {
		ctrl.(*beep.Ctrl).Paused = true
		return true
	})
	speaker.Unlock()
	sm.paused = true
}

func (sm *SoundManager) ResumeAll() {
	if !sm.paused {
		return
	}
	speaker.Lock()
	sm.ctrls.Range(func(_, ctrl interface{}) bool {
		ctrl.(*beep.Ctrl).Paused = false
		return true
	})
	speaker.Unlock()
	sm.paused = false
}

func (sm *SoundManager) PlayMP3(mp3FilePath string, vol float64, loop int) (SoundID, error) {
	if sm.avoidSameClipConcurrentPlaying {
		if id, found := sm.files.Load(mp3FilePath); found {
			return id.(SoundID), nil
		}
	}
	b, err := getResFile(mp3FilePath)
	if err != nil {
		return InvalidSoundID, err
	}
	decoded, format, err := mp3.Decode(&rsc{rs: bytes.NewReader(b)})
	if err != nil {
		return InvalidSoundID, err
	}
	return sm.play(mp3FilePath, decoded, format, vol, loop)
}

func (sm *SoundManager) Stop(id SoundID) {
	if c, ok := sm.ctrls.Load(id); ok {
		sm.ctrls.Delete(id)
		speaker.Lock()
		c.(*beep.Ctrl).Paused = true
		c.(*beep.Ctrl).Streamer = nil
		speaker.Unlock()
	}
}

func (sm *SoundManager) play(filepath string,
	s beep.StreamSeekCloser, format beep.Format, vol float64, loop int) (SoundID, error) {
	ctrl := &beep.Ctrl{Streamer: beep.Loop(loop, s)}
	resampled := beep.Resample(resamplingQuality, format.SampleRate, targetSampleRate, ctrl)
	volumed := &effects.Volume{
		Streamer: resampled,
		Base:     2,
		Volume:   vol,
	}
	id := SoundID(cwin.GenUID())
	sm.ctrls.Store(id, ctrl)
	if sm.avoidSameClipConcurrentPlaying {
		sm.files.Store(filepath, id)
	}
	speaker.Play(beep.Seq(volumed, beep.Callback(func() {
		sm.ctrls.Delete(id)
		if sm.avoidSameClipConcurrentPlaying {
			sm.files.Delete(filepath)
		}
	})))
	return id, nil
}

type rsc struct {
	rs io.ReadSeeker
}

func (r *rsc) Read(p []byte) (n int, err error) {
	return r.rs.Read(p)
}

func (r *rsc) Seek(offset int64, whence int) (int64, error) {
	return r.rs.Seek(offset, whence)
}

func (r *rsc) Close() error {
	return nil
}
