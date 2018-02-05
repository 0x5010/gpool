/*
 *
 * Created by 0x5010 on 2018/01/12.
 * gpool
 * https://github.com/0x5010/gpool
 *
 * Copyright 2017 0x5010.
 * Licensed under the MIT license.
 *
 */
package gpool

import (
	"context"
	"sync"
)

// GPool 协程池
type GPool struct {
	limit  int
	queue  chan func(ctx context.Context)
	wg     *sync.WaitGroup
	wait   bool
	cancel context.CancelFunc
}

// New 初始化协程池
func New(limit, jobCount int, wait bool) *GPool {
	if limit > jobCount {
		limit = jobCount
	}
	jQueue := make(chan func(ctx context.Context), jobCount)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	gp := &GPool{
		limit:  limit,
		queue:  jQueue,
		wait:   wait,
		wg:     &wg,
		cancel: cancel,
	}

	gp.Start(ctx)

	return gp
}

// AddJob 添加任务
func (gp *GPool) AddJob(fn func(ctx context.Context)) {
	if gp.wait {
		gp.wg.Add(1)
	}
	gp.queue <- fn
}

// Start 协程池运行
func (gp *GPool) Start(ctx context.Context) {
	for i := 0; i < gp.limit; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return

				case fn := <-gp.queue:
					fn(ctx)
					if gp.wait {
						gp.wg.Done()
					}
				}
			}
		}()
	}
}

// Wait 等待全部任务运行完
func (gp *GPool) Wait() {
	if gp.wait {
		gp.wg.Wait()
		gp.cancel()
	}
}

// Stop 强制终止
func (gp *GPool) Stop() {
	gp.cancel()
	for range gp.queue {
		gp.wg.Done()
	}
}
