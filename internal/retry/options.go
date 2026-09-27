package retry

import "time"

type retryOptions struct {
	Tries    int
	Delay    time.Duration
	MaxDelay time.Duration
	Factor   float64
	BackOff  time.Duration // compat: sets Factor (e.g. BackOff(2*time.Second) => Factor=2)
}

func NewOptions(options ...func(*retryOptions)) *retryOptions {
	o := &retryOptions{
		Tries:    -1,
		Delay:    time.Second,
		MaxDelay: 0,
		Factor:   1.0,
		BackOff:  1,
	}

	for _, option := range options {
		option(o)
	}
	if o.Factor <= 0 {
		o.Factor = 1.0
	}

	return o
}

func Tries(tries int) func(*retryOptions) {
	return func(o *retryOptions) {
		o.Tries = tries
	}
}

func Delay(delay time.Duration) func(*retryOptions) {
	return func(o *retryOptions) {
		o.Delay = delay
	}
}

func MaxDelay(maxDelay time.Duration) func(*retryOptions) {
	return func(o *retryOptions) {
		o.MaxDelay = maxDelay
	}
}

func BackOff(backOff time.Duration) func(*retryOptions) {
	return func(o *retryOptions) {
		o.BackOff = backOff
		if backOff > 1 {
			o.Factor = float64(backOff) / float64(time.Second)
			if o.Factor < 1 {
				o.Factor = 1
			}
		}
	}
}

func canRetry(o *retryOptions, tries int) bool {
	if o.Tries < 0 {
		return true
	}
	return tries < o.Tries
}
