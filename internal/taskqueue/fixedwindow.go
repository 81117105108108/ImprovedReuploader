package taskqueue

import (
	"sync"
	"time"
)

type fixedWindow struct {
	requests int
	window   time.Duration
	limit    int
	start    time.Time
	mu       sync.Mutex
}

func newFixedWindow(window time.Duration, limit int) *fixedWindow {
	return &fixedWindow{
		window: window,
		limit:  limit,
		start:  time.Now(),
	}
}

func (w *fixedWindow) timeRemaining(t time.Time) time.Duration {
	reset := w.start.Add(w.window)
	return reset.Sub(t)
}

func (w *fixedWindow) Increment() bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now()
	if now.After(w.start.Add(w.window)) {
		w.start = now
		w.requests = 0
	}

	if w.requests >= w.limit {
		return false
	}

	w.requests++
	return true
}

func (w *fixedWindow) Decrement() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.requests > 0 {
		w.requests--
	}
}

func (w *fixedWindow) Wait() {
	if w.Increment() {
		return
	}

	w.mu.Lock()
	nextWindowTime := w.start.Add(w.window)
	w.mu.Unlock()

	if d := time.Until(nextWindowTime); d > 0 {
		time.Sleep(d)
	}

	w.mu.Lock()
	now := time.Now()
	if now.After(w.start.Add(w.window)) {
		w.start = now
		w.requests = 0
	}
	w.requests++
	w.mu.Unlock()
}
