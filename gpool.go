package gpool

import (
	"context"
	"sync"
)

// Job 任务
type Job struct {
	fn func(ctx context.Context)
}

// Worker 任务消费者
type Worker struct {
	jobCacheQueue chan Job
	wait          bool
	wg            *sync.WaitGroup
}

// Run 任务协程
func (w *Worker) run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case job := <-w.jobCacheQueue:
				job.fn(ctx)
				if w.wait {
					w.wg.Done()
				}
			}
		}
	}()
}

// GPool 协程池
type GPool struct {
	maxWorkers    int
	workers       []*Worker
	jobCacheQueue chan Job
	wg            *sync.WaitGroup
	wait          bool
	ctx           context.Context
	cancel        context.CancelFunc
}

// New 初始化协程池
func New(maxWorkers, jobCacheQueueLen int, wait bool) *GPool {
	jobCacheQueue := make(chan Job, jobCacheQueueLen)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	gp := &GPool{
		maxWorkers:    maxWorkers,
		jobCacheQueue: jobCacheQueue,
		wait:          wait,
		wg:            &wg,
		cancel:        cancel,
	}

	gp.Start(ctx)

	return gp
}

// Start 协程池运行
func (gp *GPool) Start(ctx context.Context) {
	for i := 0; i < gp.maxWorkers; i++ {
		worker := &Worker{
			jobCacheQueue: gp.jobCacheQueue,
			wait:          gp.wait,
			wg:            gp.wg,
		}
		worker.run(ctx)
	}
}

// Stop 强制终止
func (gp *GPool) Stop() {
	gp.cancel()
	for _ = range gp.jobCacheQueue {
		gp.wg.Done()
	}
}

// Wait 等待全部任务运行完
func (gp *GPool) Wait() {
	if gp.wait {
		gp.wg.Wait()
		gp.cancel()
	}
}

// AddJob 添加任务
func (gp *GPool) AddJob(fn func(ctx context.Context)) {
	job := Job{
		fn: fn,
	}
	if gp.wait {
		gp.wg.Add(1)
	}
	gp.jobCacheQueue <- job
}
