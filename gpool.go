package gpool

import (
	"context"
	"sync"
)

// Job 任务
type Job struct {
	fn func()
}

// Worker 任务消费者
type Worker struct {
	jobCacheQueue chan Job
	wait          bool
	wg            *sync.WaitGroup
}

// Run 任务协程
func (w *Worker) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case job := <-w.jobCacheQueue:
				job.fn()
				if w.wait {
					w.wg.Done()
				}
			}
		}
	}()
}

// WorkerPool 协程池
type WorkerPool struct {
	maxWorkers    int
	workers       []*Worker
	jobCacheQueue chan Job
	wg            *sync.WaitGroup
	wait          bool
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewWorkerPool 初始化协程池
func NewWorkerPool(maxWorkers, jobCacheQueueLen int) *WorkerPool {
	jobCacheQueue := make(chan Job, jobCacheQueueLen)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	wp := &WorkerPool{
		maxWorkers:    maxWorkers,
		jobCacheQueue: jobCacheQueue,
		wait:          true,
		wg:            &wg,
		cancel:        cancel,
	}

	wp.Start(ctx)

	return wp
}

// Start 协程池运行
func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < wp.maxWorkers; i++ {
		worker := &Worker{
			jobCacheQueue: wp.jobCacheQueue,
			wait:          wp.wait,
			wg:            wp.wg,
		}
		worker.Run(ctx)
	}
}

// Stop 强制终止
func (wp *WorkerPool) Stop() {
	wp.cancel()
	for _ = range wp.jobCacheQueue {
		wp.wg.Done()
	}
}

// Wait 等待全部任务运行完
func (wp *WorkerPool) Wait() {
	if wp.wait {
		wp.wg.Wait()
	}
	wp.cancel()
}

// AddJob 添加任务
func (wp *WorkerPool) AddJob(fn func()) {
	job := Job{
		fn: fn,
	}
	if wp.wait {
		wp.wg.Add(1)
	}
	wp.jobCacheQueue <- job
}
