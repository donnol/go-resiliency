package deadline

import (
	"errors"
	"log"
	"testing"
	"time"
)

func takesFiveMillis(stopper <-chan struct{}) error {
	time.Sleep(5 * time.Millisecond)
	return nil
}

func takesTwentyMillis(stopper <-chan struct{}) error {
	time.Sleep(20 * time.Millisecond)
	return nil
}

func returnsError(stopper <-chan struct{}) error {
	return errors.New("foo")
}

func batch(stopper <-chan struct{}) error {
	// 一边执行任务，一边检查stopper
	for {
		select {
		case <-stopper:
			return errors.New("timeout")
		default:
			// 正常逻辑
			log.Printf("batch\n")

			// 用睡眠来模拟真实业务的执行；
			// 显然，真实的业务执行时间是不可控的；
			time.Sleep(4 * time.Millisecond) // run 3 times under 10ms timeout
			// time.Sleep(5 * time.Millisecond) // run 2 times under 10ms timeout
		}
	}
}

func TestDeadline(t *testing.T) {
	dl := New(10 * time.Millisecond)

	if err := dl.Run(takesFiveMillis); err != nil {
		t.Error(err)
	}

	if err := dl.Run(takesTwentyMillis); err != ErrTimedOut {
		t.Error(err)
	}

	if err := dl.Run(returnsError); err.Error() != "foo" {
		t.Error(err)
	}

	if err := dl.Run(batch); err != ErrTimedOut {
		t.Fatal(err)
	}

	done := make(chan struct{})
	err := dl.Run(func(stopper <-chan struct{}) error {
		<-stopper
		close(done)
		return nil
	})
	if err != ErrTimedOut {
		t.Error(err)
	}
	<-done
}

func ExampleDeadline() {
	dl := New(1 * time.Second)

	err := dl.Run(func(stopper <-chan struct{}) error {
		// do something possibly slow
		// check stopper function and give up if timed out
		return nil
	})

	switch err {
	case ErrTimedOut:
		// execution took too long, oops
	default:
		// some other error
	}
}
