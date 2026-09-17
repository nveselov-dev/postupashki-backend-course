package mutex

import (
	"primitives/internal/futex"
	"runtime"
	"sync/atomic"
)

const (
	free = iota
	held
	contended
)

type Mutex struct {
	state uint32
}

func (m *Mutex) Lock() {
	if m.TryLock() {
		return
	}

	for {
		if atomic.LoadUint32(&m.state) == contended || atomic.CompareAndSwapUint32(&m.state, held, contended) {

			for range 100 {
				if atomic.LoadUint32(&m.state) == free {
					break
				}
				runtime.Gosched()
			}

			futex.Wait(&m.state, contended)
		}
		runtime.Gosched()
		if atomic.CompareAndSwapUint32(&m.state, free, contended) {
			return
		}
	}
}

func (m *Mutex) TryLock() bool {
	return atomic.CompareAndSwapUint32(&m.state, free, held)
}

func (m *Mutex) Unlock() {
	old := atomic.SwapUint32(&m.state, free)
	if old == free {
		panic("unlock of unlocked lock")
	}
	if old == contended {
		futex.Wake(&m.state)
	}
}
