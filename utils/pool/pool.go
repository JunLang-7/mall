package pool

import (
	"sync"

	"github.com/panjf2000/ants/v2"
)

type Pool struct {
	pool *ants.Pool
	wg   *sync.WaitGroup
}

func NewPoolWithSize(size int) *Pool {
	pool, _ := ants.NewPool(size)
	return &Pool{
		pool: pool,
		wg:   &sync.WaitGroup{},
	}
}

func (p *Pool) RunGo(f func()) {
	p.wg.Add(1)
	_ = p.pool.Submit(func() {
		defer p.wg.Done()
		f()
	})
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Release() {
	p.pool.Release()
}
