package tokenbucket

import (
	"time"
)

type Limiter struct {
	// Why use float64 type?
	// Time.Seconds() return float64,
	// so use float64 type is easy to calculate with Time.Seconds().
	rate   float64
	burst  float64
	tokens float64
	// Last time for consuming tokens.
	last time.Time
}

// rate:  number of events per second.
// burst: peak of events that can be consumed.
func NewLimiter(rate int, burst int) *Limiter {
	return &Limiter{
		rate:   float64(rate),
		burst:  float64(burst),
		tokens: float64(burst),
		last:   time.Now(),
	}
}

func (l *Limiter) Allow(now time.Time) (isAllowed bool, sleep time.Duration) {
	elapsed := now.Sub(l.last)

	l.tokens = l.tokens + l.rate*elapsed.Seconds()
	l.last = now

	if l.tokens > l.burst {
		l.tokens = l.burst
	}

	l.tokens = l.tokens - 1

	if l.tokens < 0 {
		return false, time.Duration(-l.tokens / l.rate * float64(time.Second))
	} else {
		return true, 0
	}
}

func (l *Limiter) Wait() {
	isAllowed, sleep := l.Allow(time.Now())
	if !isAllowed {
		time.Sleep(sleep)
	}
}
