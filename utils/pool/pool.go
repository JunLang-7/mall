package pool

import (
	"github.com/panjf2000/ants/v2"
)

type Pool struct {
	pool *ants.Pool
}

func NewPoolWithSize(size int) *Pool {
	pool, _ := ants.NewPool(size)
	return &Pool{
		pool: pool,
	}
}

func (p *Pool) RunGo(f func()) {
	_ = p.pool.Submit(f)
}

func (p *Pool) Wait() {
	p.pool.Release()
}
