package cgame

import (
	"path"
	"runtime"
	"sync"
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
