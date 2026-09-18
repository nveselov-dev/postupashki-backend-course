package once

import (
	"runtime"
	"sync/atomic"
)

type Once struct {
	state uint32
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.state) == 2 {
		return
	}

	if atomic.CompareAndSwapUint32(&o.state, 0, 1) {
		func() {
			defer func() {
				recover()
			}()
			f()
		}()
		atomic.StoreUint32(&o.state, 2)
		return
	}

	for atomic.LoadUint32(&o.state) != 2 {
		runtime.Gosched()
	}
}

func (o *Once) Done() bool {
	return atomic.LoadUint32(&o.state) == 2
}
