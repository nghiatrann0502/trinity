package common

import (
	"fmt"
	"time"
)

func RetryWithBackoff(maxRetries int, baseDelay time.Duration, task func() error) {
	delay := baseDelay
	var err error
	for i := 0; i < maxRetries; i++ {
		err = task()
		if err == nil {
			return
		}
		fmt.Printf("Attempt %d failed: %v\n", i+1, err)
		time.Sleep(delay)
		delay *= 2
	}

	// return err
}
