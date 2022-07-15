package taskrunner

import (
	"context"
	"strconv"
	"testing"
	"time"
)

func TestSampleTask(t *testing.T) {
	ch := make(chan int, 5000)
	fn := func(ctx context.Context, mtr *Task) {
		for {
			num := <-ch
			mtr.AddMeta(strconv.Itoa(num), nil)
			tc := time.After(time.Second * 10)
			select {
			case <-ctx.Done():
				return
			case <-tc:
				//fmt.Println("hello World", ctx.Err())
			}
		}
	}
	tm := NewTaskManager()
	for i := 0; i < 20; i++ {
		tm.GO(fn)
	}
	for i := 0; i < 100; i++ {
		ch <- i
	}
	time.Sleep(time.Second)
	tm.CancelTaskFromMetaKey("1")
	time.Sleep(time.Second * 4)
}
