package retry

import (
	"time"
)

func getDelay(o *retryOptions, tries int) time.Duration {
	factor := o.Factor
	if factor <= 0 {
		factor = 1.0
	}
	delay := float64(o.Delay)
	for i := 1; i < tries; i++ {
		delay *= factor
		if o.MaxDelay > 0 && time.Duration(delay) > o.MaxDelay {
			return o.MaxDelay
		}
	}

	if o.MaxDelay == 0 {
		return time.Duration(delay)
	}

	if time.Duration(delay) > o.MaxDelay {
		return o.MaxDelay
	}
	return time.Duration(delay)
}

func Do[T any](options *retryOptions, callback func(try int) (T, error)) (T, error) {
	var tries int

	for {
		tries++

		res, err := callback(tries)
		if err == nil {
			return res, nil
		}

		switch err := err.(type) {
		case *ExitRetry:
			return res, err.Err
		case *ContinueRetry:
			if !canRetry(options, tries) {
				return res, err.Err
			}

			time.Sleep(getDelay(options, tries))
		default:
			return res, err
		}
	}
}
