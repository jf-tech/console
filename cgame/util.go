package cgame

import (
	"path"
	"runtime"
	"sync"
	"time"
)

type ThreadSafeFIFO struct {
	sync.Mutex
	elems []interface{}
}

func (f *ThreadSafeFIFO) Push(e interface{}) {
	f.Lock()
	defer f.Unlock()
	f.elems = append(f.elems, e)
}

func (f *ThreadSafeFIFO) TryPop() (interface{}, bool) {
	f.Lock()
	defer f.Unlock()
	n := len(f.elems)
	if n <= 0 {
		return nil, false
	}
	ret := f.elems[0]
	copy(f.elems[0:], f.elems[1:])
	f.elems = f.elems[:n-1]
	return ret, true
}

func NewThreadSafeFIFO(cap int) *ThreadSafeFIFO {
	return &ThreadSafeFIFO{
		elems: make([]interface{}, 0, cap),
	}
}

func GetCurFileDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

type DurationCounter struct {
	clock     *Clock
	total     time.Duration
	startedOn time.Duration
}

func (d *DurationCounter) Start() {
	if d.started() {
		panic("cannot Start() a started counter")
	}
	d.startedOn = d.clock.Now()
}

func (d *DurationCounter) Stop() {
	if !d.started() {
		panic("cannot Stop() a stopped counter")
	}
	d.total += d.clock.Now() - d.startedOn
	d.startedOn = -1
}

func (d *DurationCounter) Reset() {
	if d.started() {
		panic("cannot Reset() a started counter")
	}
	d.total = 0
	d.startedOn = -1
}

func (d *DurationCounter) Total() time.Duration {
	if d.started() {
		panic("cannot get Total() on a started counter")
	}
	return d.total
}

func (d *DurationCounter) started() bool {
	return d.startedOn >= 0
}

func NewDurationCounter(clock *Clock) *DurationCounter {
	return &DurationCounter{clock: clock, startedOn: -1}
}
