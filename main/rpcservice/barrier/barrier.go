package barrier

import "sync"

type Barrier struct {
	c      int
	n      int
	m      sync.Mutex
	before chan int
	after  chan int
}

func New(n int) *Barrier {
	b := Barrier{
		n:      n,
		before: make(chan int, n),
		after:  make(chan int, n),
	}
	return &b
}

func (b *Barrier) Wait_on_barrier() {
	b.m.Lock()
	b.c += 1
	if b.c == b.n {
		// open gate
		for i := 0; i < b.n; i++ {
			b.before <- 1
		}
	}
	b.m.Unlock()
	<-b.before
}
