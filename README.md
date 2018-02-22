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
var fn1, fn2 func() // the function which you want to  execute, anonymous functions form closures is better
var fn3, fn4 func(ctx context.Context) // with context, will canceled when pool stop
var limit, jobCount int   // the number of goroutine and job
var wait bool                          // whether blocking

gp := gpool.New(limit, jobCount, wait)
gp.AddJob(fn1)
gp.AddJob(fn2)
gp.AddJobWithCtx(fn3)
gp.AddJobWithCtx(fn4)

if wait {
	gp.Wait()
}
```

termination:

```go
gp.Stop()
...
```
