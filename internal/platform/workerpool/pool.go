package workerpool

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"
	"time"
)

type Job func(ctx context.Context) error

type task struct {
	ctx context.Context
	job Job
}

type Pool struct {
	name    string
	tasks   chan task
	timeout time.Duration
	wg      sync.WaitGroup

	mu     sync.RWMutex
	closed bool
}

func New(name string, workers, queueSize int, timeout time.Duration) *Pool {
	p := &Pool{
		name:    name,
		tasks:   make(chan task, queueSize),
		timeout: timeout,
	}

	for range workers {
		p.wg.Go(p.worker)
	}

	return p
}

// Submit 不會阻塞。佇列已滿或已關閉時回傳 false，由呼叫端決定怎麼處理
func (p *Pool) Submit(ctx context.Context, j Job) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.closed {
		return false // 關機後才送進來的任務：回傳 false，而不是對已關閉的 channel 送值而 panic
	}

	select {
	case p.tasks <- task{ctx: context.WithoutCancel(ctx), job: j}:
		return true
	default:
		return false
	}
}

func (p *Pool) Shutdown(ctx context.Context) error {
	p.mu.Lock()
	if !p.closed {
		p.closed = true
		close(p.tasks)
	}
	p.mu.Unlock()

	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("shutdown pool %s: %w", p.name, ctx.Err())
	}
}

func (p *Pool) worker() {
	for t := range p.tasks {
		p.run(t)
	}
}

func (p *Pool) run(t task) {
	ctx, cancel := context.WithTimeout(t.ctx, p.timeout)
	defer cancel()

	defer func() {
		if rec := recover(); rec != nil {
			slog.ErrorContext(ctx, "worker panic",
				"pool", p.name,
				"panic", rec,
				"stack", string(debug.Stack()),
			)
		}
	}()

	// 背景任務的錯誤邊界：統一在這裡 log 一次
	if err := t.job(ctx); err != nil {
		slog.ErrorContext(ctx, "background job failed",
			"pool", p.name,
			"error", err,
		)
	}

}
