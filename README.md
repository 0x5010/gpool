This is a goroutine pool in the [Go](http:golang.org) for easier handling and cancellation.

[![Build Status](https://travis-ci.org/0x5010/gpool.png?branch=master)](https://travis-ci.org/0x5010/gpool)

Installation
-----------

	go get github.com/0x5010/gpool

Requirements
-----------

* Need at least `go1.7` or newer.

Usage
-----------

Create and run a gpool:
```go
wp := gpool.NewWorkerPool(maxWorkers, jobCacheQueueLen)
wp.AddJob(fn1)
...
wp.AddJob(fnn)
wp.Wait()
```
cancel:
```go
wp.Stop()
...
```


