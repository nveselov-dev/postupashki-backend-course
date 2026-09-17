package spinlock

import (
	"runtime"
	"sync/atomic"
)

type Spinlock struct {
	locked atomic.Bool
}

func (s *Spinlock) Lock() {
	for !s.TryLock() {
		runtime.Gosched()
	}
}

func (s *Spinlock) TryLock() bool {
	return s.locked.CompareAndSwap(false, true)
}

func (s *Spinlock) Unlock() {
	if !s.locked.CompareAndSwap(true, false) {
		panic("unlock of unlocked lock")
	}
}

type TTAS struct {
	locked atomic.Bool
}

func (s *TTAS) Lock() {
	for {
		for s.locked.Load() {
			runtime.Gosched()
		}

		if s.TryLock() {
			return
		}
	}
}

func (s *TTAS) TryLock() bool {
	if !s.locked.Load() {
		return s.locked.CompareAndSwap(false, true)
	}
	return false
}

func (s *TTAS) Unlock() {
	if !s.locked.CompareAndSwap(true, false) {
		panic("unlock of unlocked lock")
	}
}
