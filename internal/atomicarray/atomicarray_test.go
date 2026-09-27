package atomicarray

import (
	"sync"
	"testing"
)

func TestUpdateEmptySafe(t *testing.T) {
	empty := []int{}
	a := New(&empty)
	a.Update(func(arr []int) []int {
		if len(arr) != 0 {
			t.Error("want empty")
		}
		return nil
	})
	if len(a.Load()) != 0 {
		t.Fatal("mutated empty")
	}
}

func TestUpdateCopyOnWrite(t *testing.T) {
	base := []int{1, 2, 3}
	a := New(&base)
	a.Update(func(arr []int) []int {
		out := append([]int(nil), arr...)
		out[0] = 9
		return out
	})
	if got := a.Load(); got[0] != 9 || base[0] != 1 {
		t.Fatalf("cow failed: got %v base %v", got, base)
	}
}

func TestConcurrentLoadUpdate(t *testing.T) {
	base := []int{1, 2, 3}
	a := New(&base)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _ = a.Load() }()
		wg.Add(1)
		go func(v int) { defer wg.Done(); a.Store([]int{v}) }(i)
	}
	wg.Wait()
}
