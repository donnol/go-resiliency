// Package deadline implements the deadline (also known as "timeout") resiliency pattern for Go.
package deadline

import (
	"errors"
	"time"
)

// ErrTimedOut is the error returned from Run when the deadline expires.
var ErrTimedOut = errors.New("timed out waiting for function to finish")

// Deadline implements the deadline/timeout resiliency pattern.
type Deadline struct {
	timeout time.Duration
}

// New constructs a new Deadline with the given timeout.
func New(timeout time.Duration) *Deadline {
	return &Deadline{
		timeout: timeout,
	}
}

// Run runs the given function, passing it a stopper channel. If the deadline passes before
// the function finishes executing, Run returns ErrTimeOut to the caller and closes the stopper
// channel so that the work function can attempt to exit gracefully. It does not (and cannot)
// simply kill the running function, so if it doesn't respect the stopper channel then it may
// keep running after the deadline passes. If the function finishes before the deadline, then
// the return value of the function is returned from Run.
func (d *Deadline) Run(work func(<-chan struct{}) error) error {
	result := make(chan error)
	stopper := make(chan struct{})

	// 两个select阻塞住，直到结果返回或超时

	go func() {
		// 如果work执行的时间远超timeout，不也没办法让work在超时的时间结束吗；
		// 所以要将stopper管道传进去，work函数实现时就要根据它来判断是不是要及时停止
		//
		// 所以，如果work函数的实现里，对于stopper的检查不够及时，也会运行得比超时时间更长
		value := work(stopper)

		// work一直在执行就不返回，那还怎么进入select呢？
		// 所以，work里对于stopper的检查是非常必要的，虽然可能会稍慢一点，但这个goroutine最终还是会停止的
		select {
		case result <- value:
		case <-stopper:
		}
	}()

	select {
	case ret := <-result:
		return ret
	case <-time.After(d.timeout):
		close(stopper)
		return ErrTimedOut
	}
}
