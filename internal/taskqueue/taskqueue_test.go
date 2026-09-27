package taskqueue

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestSmoothQueueCompletes(t *testing.T) {
	q := NewSmoothQueue[int](4, 600)
	var sum atomic.Int64
	done := make([]chan TaskResult[int], 5)
	for i := 0; i < 5; i++ {
		i := i
		done[i] = q.QueueTask(func() (int, error) { sum.Add(int64(i)); return i, nil })
	}
	for i, ch := range done {
		r := <-ch
		if r.Error != nil || r.Result != i {
			t.Fatalf("task %d got %+v", i, r)
		}
	}
	if sum.Load() != 10 {
		t.Fatalf("sum %d", sum.Load())
	}
}

func TestUniformPacerGap(t *testing.T) {
	p := NewUniformPacer(50 * time.Millisecond)
	p.Wait()
	start := time.Now()
	p.Wait()
	if elapsed := time.Since(start); elapsed < 40*time.Millisecond {
		t.Fatalf("gap too short: %v", elapsed)
	}
}

func TestFixedWindowSerializes(t *testing.T) {
	q := New[int](200*time.Millisecond, 2)
	start := time.Now()
	chs := []chan TaskResult[int]{
		q.QueueTask(func() (int, error) { return 1, nil }),
		q.QueueTask(func() (int, error) { return 2, nil }),
		q.QueueTask(func() (int, error) { return 3, nil }),
	}
	for _, ch := range chs {
		<-ch
	}
	// 3 tasks with limit 2 per 200ms must span at least one window
	if elapsed := time.Since(start); elapsed < 150*time.Millisecond {
		t.Fatalf("expected serialization, took %v", elapsed)
	}
}
