package semaphore

import (
	"primitives/internal/futex"
	"sync/atomic"
)

type Semaphore struct {
	permits uint32
}

func New(n int) *Semaphore {
	return &Semaphore{permits: uint32(n)}
}

func (s *Semaphore) Acquire() {
	for {
		cur := atomic.LoadUint32(&s.permits)
		if cur > 0 {
			if atomic.CompareAndSwapUint32(&s.permits, cur, cur-1) {
				return
			}
			continue
		}
		futex.Wait(&s.permits, cur)
	}
}

func (s *Semaphore) TryAcquire() bool {
	for {
		c := atomic.LoadUint32(&s.permits)
		if c == 0 {
			return false
		}
		if atomic.CompareAndSwapUint32(&s.permits, c, c-1) {
			return true
		}
	}
}

func (s *Semaphore) Release() {
	for {
		cur := atomic.LoadUint32(&s.permits)
		newC := cur + 1
		if newC < cur {
			panic("semaphore overflow")
		}
		if atomic.CompareAndSwapUint32(&s.permits, cur, newC) {
			futex.Wake(&s.permits)
			return
		}
	}
}

func (s *Semaphore) Available() int {
	return int(atomic.LoadUint32(&s.permits))
}
