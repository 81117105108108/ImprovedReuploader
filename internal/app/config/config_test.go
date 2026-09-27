package config

import "testing"

func TestGetIntFallback(t *testing.T) {
	if got := parseIntOr("", 7); got != 7 {
		t.Fatalf("empty %d", got)
	}
	if got := parseIntOr("abc", 7); got != 7 {
		t.Fatalf("bad %d", got)
	}
	if got := parseIntOr("42", 7); got != 42 {
		t.Fatalf("good %d", got)
	}
	if got := parseIntOr("-3", 7); got != -3 {
		t.Fatalf("neg %d", got)
	}
}

func TestRateDefaults(t *testing.T) {
	if StartsPerMinute() <= 0 || MaxConcurrent() <= 0 || MaxParallelChunks() <= 0 || UploadTries() <= 0 {
		t.Fatal("rate defaults must be positive")
	}
}
