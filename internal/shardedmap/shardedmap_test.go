package shardedmap

import (
	"sync"
	"testing"
)

func TestGetOrCreateSameInstance(t *testing.T) {
	m := New[int]()
	a := m.GetOrCreateShard("k")
	b := m.GetOrCreateShard("k")
	if a != b {
		t.Fatal("different shards")
	}
	a.Set("x", 1)
	if v, ok := b.Get("x"); !ok || v != 1 {
		t.Fatal("not shared")
	}
}

func TestConcurrentGetOrCreate(t *testing.T) {
	m := New[int]()
	var wg sync.WaitGroup
	results := make([]*Shard[int], 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) { defer wg.Done(); results[i] = m.GetOrCreateShard("same") }(i)
	}
	wg.Wait()
	for i := 1; i < 10; i++ {
		if results[i] != results[0] {
			t.Fatal("race created duplicates")
		}
	}
}
