package pools

import "sync"

// Pool goroutine pool
type Pool interface {
	// Wait 等待令牌
	Wait()
	// Done 归还令牌
	Done()
	// Num 当前发放的令牌书
	Num() int
	// Size 总令牌数
	Size() int

	// WaitAll 同步等待令牌全部归还
	WaitAll()
	// AsyncWaitAll 异步等待令牌全部归还
	AsyncWaitAll() <-chan struct{}
}

type pool struct {
	pool chan struct{}

	wg sync.WaitGroup
}

func (p *pool) Init(size int) {
	if size >= 0 {
		p.pool = make(chan struct{}, size)
	}
}

func (p *pool) Wait() {
	if p.pool != nil {
		p.wg.Add(1)
		p.pool <- struct{}{}
	}
}

func (p *pool) Done() {
	if p.pool != nil {
		<-p.pool
		p.wg.Done()
	}
}

func (p *pool) Num() int {
	if p.pool != nil {
		return len(p.pool)
	}
	return 0
}

func (p *pool) Size() int {
	if p.pool != nil {
		return cap(p.pool)
	}
	return 0
}

func (p *pool) WaitAll() {
	p.wg.Wait()
}

func (p *pool) AsyncWaitAll() <-chan struct{} {
	sig := make(chan struct{})
	go func() {
		p.WaitAll()
		sig <- struct{}{}
	}()
	return sig
}
