semaphore
=========

[![Build Status](https://travis-ci.org/donnol/go-resiliency.svg?branch=master)](https://travis-ci.org/donnol/go-resiliency)
[![GoDoc](https://godoc.org/github.com/donnol/go-resiliency/semaphore?status.svg)](https://godoc.org/github.com/donnol/go-resiliency/semaphore)
[![Code of Conduct](https://img.shields.io/badge/code%20of%20conduct-active-blue.svg)](https://donnol.github.io/conduct.html)

The semaphore resiliency pattern for golang.

Creating a semaphore takes two parameters:
- ticket count (how many tickets to give out at once)
- timeout (how long to wait for a ticket if none are currently available)

```go
sem := semaphore.New(3, 1*time.Second)

if err := sem.Acquire(); err != nil {
	// could not acquire semaphore
	return err
}
defer sem.Release()
```
