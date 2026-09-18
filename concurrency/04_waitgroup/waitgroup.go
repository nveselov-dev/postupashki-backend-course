package waitgroup

import (
	"runtime"
	"sync/atomic"
)

type WaitGroup struct {
	count uint32
}

func (wg *WaitGroup) Add(delta int) {
	for {
		c := atomic.LoadUint32(&wg.count)
		newC := int32(c) + int32(delta)

		if newC < 0 {
			panic("negative counter")
		}

		if atomic.CompareAndSwapUint32(&wg.count, c, uint32(newC)) {
			return
		}
	}
}

func (wg *WaitGroup) Done() {
	wg.Add(-1)
}

func (wg *WaitGroup) Wait() {
	for atomic.LoadUint32(&wg.count) != 0 {
		runtime.Gosched()
	}
}
