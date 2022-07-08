retrier
=======

[![Build Status](https://travis-ci.org/donnol/go-resiliency.svg?branch=master)](https://travis-ci.org/donnol/go-resiliency)
[![GoDoc](https://godoc.org/github.com/donnol/go-resiliency/retrier?status.svg)](https://godoc.org/github.com/donnol/go-resiliency/retrier)
[![Code of Conduct](https://img.shields.io/badge/code%20of%20conduct-active-blue.svg)](https://donnol.github.io/conduct.html)

The retriable resiliency pattern for golang.

Creating a retrier takes two parameters:
- the times to back-off between retries (and implicitly the number of times to
  retry)
- the classifier that determines which errors to retry

```go
r := retrier.New(retrier.ConstantBackoff(3, 100*time.Millisecond), nil)

err := r.Run(func() error {
	// do some work
	return nil
})

if err != nil {
	// handle the case where the work failed three times
}
```
