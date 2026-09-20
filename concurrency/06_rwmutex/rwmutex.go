package rwmutex

import (
	"primitives/internal/futex"
	"sync/atomic"
)

const (
	writer = 1 << 31
)

type RWMutex struct {
	state uint32
}

func (rw *RWMutex) RLock() {
	for {
		c := atomic.LoadUint32(&rw.state)
		if c < writer-1 {
			if atomic.CompareAndSwapUint32(&rw.state, c, c+1) {
				return
			}
		}

		futex.Wait(&rw.state, c)
	}
}

func (rw *RWMutex) RUnlock() {
	for {
		c := atomic.LoadUint32(&rw.state)
		if c == 0 {
			panic("runtime error: unlock of unlocked RWMutex")
		}
		if atomic.CompareAndSwapUint32(&rw.state, c, c-1) {
			if c-1 == 0 {
				futex.Wake(&rw.state)
			}
			return
		}
	}
}

func (rw *RWMutex) Lock() {
	for {
		c := atomic.LoadUint32(&rw.state)
		if c == 0 {
			if atomic.CompareAndSwapUint32(&rw.state, 0, writer) {
				return
			}
		} else {
			futex.Wait(&rw.state, c)
		}
	}
}

func (rw *RWMutex) Unlock() {
	old := atomic.SwapUint32(&rw.state, 0)
	if old != writer {
		panic("unlock of unlocked RWMutex")
	}
	futex.WakeAll(&rw.state)
}
