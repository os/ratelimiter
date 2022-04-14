package rate

import (
	"testing"
	"time"
)

func TestFixedWindowLimiter_Allow(t *testing.T) {
	for _, tc := range []struct {
		Description       string
		Limit             int
		CallCount         int
		ExpectedRemaining int
		ExpectedError     error
	}{
		{
			Description:       "Less than limit number of calls returns no error with remaining balance",
			Limit:             10,
			CallCount:         9,
			ExpectedRemaining: 1,
			ExpectedError:     nil,
		},
		{
			Description:       "Limit number of calls returns no error with 0 remaining balance",
			Limit:             10,
			CallCount:         10,
			ExpectedRemaining: 0,
			ExpectedError:     nil,
		},
		{
			Description:       "More than limit number of calls returns ErrRateLimitExceeded error with 0 remaining balance",
			Limit:             10,
			CallCount:         11,
			ExpectedRemaining: 0,
			ExpectedError:     ErrRateLimitExceeded,
		},
	} {
		t.Run(tc.Description, func(t *testing.T) {
			var err error
			var remaining int
			var response Response

			limiter := NewFixedWindowLimiter(tc.Limit, time.Second*10, NewMemoryStore(time.Second*10))
			for i := 0; i < tc.CallCount; i++ {
				response, err = limiter.Allow("default")
				if err == nil {
					remaining = response.Remaining
				}
			}

			if err != tc.ExpectedError {
				t.Errorf("Expected error: %s, got: %s", tc.ExpectedError, err)
			}
			if remaining != tc.ExpectedRemaining {
				t.Errorf("Expected remaining: %d, got: %d", tc.ExpectedRemaining, remaining)
			}
		})
	}
}
