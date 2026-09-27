package retry

import (
	"errors"
	"testing"
	"time"
)

func TestDoSuccess(t *testing.T) {
	v, err := Do(NewOptions(Tries(3)), func(try int) (int, error) { return 42, nil })
	if err != nil || v != 42 {
		t.Fatalf("got %v,%v want 42,nil", v, err)
	}
}

func TestDoRetriesThenSuccess(t *testing.T) {
	calls := 0
	v, err := Do(NewOptions(Tries(3), Delay(time.Millisecond)), func(try int) (int, error) {
		calls++
		if try < 3 {
			return 0, &ContinueRetry{Err: errors.New("tmp")}
		}
		return 7, nil
	})
	if err != nil || v != 7 || calls != 3 {
		t.Fatalf("got %v,%v calls=%d", v, err, calls)
	}
}

func TestDoExhausts(t *testing.T) {
	sentinel := errors.New("boom")
	_, err := Do(NewOptions(Tries(2), Delay(time.Millisecond)), func(try int) (int, error) {
		return 0, &ContinueRetry{Err: sentinel}
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("want sentinel, got %v", err)
	}
}

func TestTriesMinus1Infinite(t *testing.T) {
	calls := 0
	v, err := Do(NewOptions(Tries(-1), Delay(time.Millisecond)), func(try int) (int, error) {
		calls++
		if try < 3 {
			return 0, &ContinueRetry{Err: errors.New("tmp")}
		}
		return 9, nil
	})
	if err != nil || v != 9 || calls != 3 {
		t.Fatalf("got %v,%v calls=%d", v, err, calls)
	}
}

func TestGetDelayFactor(t *testing.T) {
	o := NewOptions(Delay(100*time.Millisecond), MaxDelay(time.Second))
	o.Factor = 2
	if d := getDelay(o, 1); d != 100*time.Millisecond {
		t.Fatalf("try1 got %v", d)
	}
	if d := getDelay(o, 2); d != 200*time.Millisecond {
		t.Fatalf("try2 got %v", d)
	}
	if d := getDelay(o, 10); d != time.Second {
		t.Fatalf("cap got %v", d)
	}
}
