This is a goroutine pool library in the [Go](http:golang.org) for easier handling and termination.

[![LICENSE](https://img.shields.io/badge/license-MIT-orange.svg)](LICENSE)
[![Build Status](https://travis-ci.org/0x5010/gpool.png?branch=master)](https://travis-ci.org/0x5010/gpool)
[![Go Report Card](https://goreportcard.com/badge/github.com/0x5010/gpool)](https://goreportcard.com/report/github.com/0x5010/gpool)

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
var fn1, fnn func(ctx context.Context) // the function which you want to  execute,  anonymous functions form closures is better
var maxWorkers, jobCacheQueueLen int   // the number of goroutine and cache size
var wait bool                          // whether blocking

gp := gpool.New(maxWorkers, jobCacheQueueLen, wait)
gp.AddJob(fn1)
...
gp.AddJob(fnn)

if wait {
	gp.Wait()
}
```


termination:
```go
gp.Stop()
...
```


