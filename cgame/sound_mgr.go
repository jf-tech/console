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
	"github.com/jf-tech/console/cutil"
	"github.com/jf-tech/console/cwin"
)

type SoundID int64

const (
	InvalidSoundID = SoundID(-1)

	targetSampleRate  = beep.SampleRate(44100)
	resamplingQuality = 4
)

var inited = false

type clipEntry struct {
	id       SoundID
	ctrl     *beep.Ctrl
	filepath string
}

type SoundManager struct {
	clips  sync.Map // [SoundID]*clipEntry
	files  sync.Map // [filepath]*clipEntry
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
	sm.clips = sync.Map{}
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
	defer speaker.Unlock()
	sm.clips.Range(func(_, clip interface{}) bool {
		clip.(*clipEntry).ctrl.Paused = true
		return true
	})
	sm.paused = true
}

func (sm *SoundManager) ResumeAll() {
	if !sm.paused {
		return
	}
	speaker.Lock()
	defer speaker.Unlock()
	sm.clips.Range(func(_, clip interface{}) bool {
		clip.(*clipEntry).ctrl.Paused = false
		return true
	})
	sm.paused = false
}

func (sm *SoundManager) Stop(id SoundID) {
	if clip, found := sm.clips.Load(id); found {
		speaker.Lock()
		defer speaker.Unlock()
		clip.(*clipEntry).ctrl.Paused = true
		clip.(*clipEntry).ctrl.Streamer = nil
		sm.clips.Delete(id)
		sm.files.Delete(id)
	}
}

func (sm *SoundManager) PlayMP3(mp3FilePath string, vol float64, loop int) (SoundID, error) {
	if sm.avoidSameClipConcurrentPlaying {
		if clip, found := sm.files.Load(mp3FilePath); found {
			sm.Stop(clip.(*clipEntry).id)
		}
	}
	b, err := cutil.LoadCachedFile(mp3FilePath)
	if err != nil {
		return InvalidSoundID, err
	}
	decoded, format, err := mp3.Decode(&rsc{rs: bytes.NewReader(b)})
	if err != nil {
		return InvalidSoundID, err
	}
	return sm.play(mp3FilePath, decoded, format, vol, loop)
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
	clip := &clipEntry{
		id:       SoundID(cwin.GenUID()),
		ctrl:     ctrl,
		filepath: filepath,
	}
	sm.clips.Store(clip.id, clip)
	if sm.avoidSameClipConcurrentPlaying {
		sm.files.Store(filepath, clip)
	}
	speaker.Play(beep.Seq(volumed, beep.Callback(func() {
		sm.clips.Delete(clip.id)
		sm.files.Delete(filepath)
	})))
	return clip.id, nil
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
