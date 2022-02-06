package cgame

import "sync"

type PairInt struct {
	A, B int
}

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
	for i := 0; i < n-1; i++ {
		f.elems[i] = f.elems[i+1]
	}
	f.elems[n-1] = nil
	f.elems = f.elems[:n-1]
	return ret, true
}

func NewThreadSafeFIFO(cap int) *ThreadSafeFIFO {
	return &ThreadSafeFIFO{
		elems: make([]interface{}, 0, cap),
	}
}
