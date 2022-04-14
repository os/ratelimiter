package rate

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrRateLimitExceeded indicates that a rate limit is exceeded.
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

// Limiter provides a method to query if a given ID is allowed or not.
type Limiter interface {
	Allow(id string) (Response, error)
}

// windowStart time for the given duration and the current time.
func windowStart(duration time.Duration) time.Time {
	timestamp := time.Now().UnixNano()
	period := duration.Nanoseconds()
	start := int64(float64(timestamp)/float64(period)) * period
	return time.Unix(0, start)
}

// windowReset time for the given duration and the current time.
func windowReset(duration time.Duration) time.Time {
	start := windowStart(duration)
	return start.Add(duration)
}

// Response of a rate limiter represents the state of the current limit window
// for the requested key.
type Response struct {
	Limit     int
	Duration  time.Duration
	Remaining int
	Reset     time.Time
}

// NewResponse returns a Response instance with the given parameters.
func NewResponse(limit int, duration time.Duration, count int, reset time.Time) Response {
	delta := float64(limit - count)
	remaining := int(math.Max(delta, 0))
	return Response{
		Limit:     limit,
		Duration:  duration,
		Remaining: remaining,
		Reset:     reset,
	}
}

// FixedWindowLimiter is an implementation of fixed window counter algorithm
// and it implements the Limiter interface.
type FixedWindowLimiter struct {
	limit    int
	duration time.Duration
	store    Store
}

// NewFixedWindowLimiter returns a FixedWindowLimiter instance with the given
// parameters.
func NewFixedWindowLimiter(limit int, duration time.Duration, store Store) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		limit:    limit,
		duration: duration,
		store:    store,
	}
}

// key for the current window and id.
func (l *FixedWindowLimiter) key(id string) string {
	start := windowStart(l.duration)
	seconds := int64(l.duration.Seconds())
	parts := []string{
		id,
		strconv.FormatInt(seconds, 10),
		strconv.FormatInt(start.Unix(), 10),
	}
	return strings.Join(parts, ":")
}

// expiry time of the current window.
func (l *FixedWindowLimiter) expiry() int64 {
	return windowReset(l.duration).Unix()
}

// Allow returns a Response which represents the state of the current limit
// window for the requested id. If the rate limit is exceeded it returns an
// ErrRateLimitExceeded error.
func (l *FixedWindowLimiter) Allow(id string) (Response, error) {
	key := l.key(id)
	count, err := l.store.Increment(key, l.expiry())
	response := NewResponse(l.limit, l.duration, count, windowReset(l.duration))
	if err != nil {
		return response, err
	}
	if count > l.limit {
		return response, ErrRateLimitExceeded
	}

	return response, nil
}
