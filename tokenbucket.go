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

// rate:          number of events per second.
// burst:         peak number of events that can be consumed.
// initialTokens: initial number of events that can be consumed.
func NewLimiter(rate int, burst int, initialTokens int) *Limiter {
	return &Limiter{
		rate:   float64(rate),
		burst:  float64(burst),
		tokens: float64(initialTokens),
		last:   time.Now(),
	}
}

// This function is used for manual control sleep flow.
func (l *Limiter) Allow(now time.Time) (isAllowed bool, sleep time.Duration) {
	elapsed := now.Sub(l.last)

	// "rate * elapsed.Seconds()" mean newly obtain tokens in the past elapsed time.
	l.tokens = l.tokens + l.rate*elapsed.Seconds()
	l.last = now

	if l.tokens > l.burst {
		l.tokens = l.burst
	}

	// Consume one token.
	l.tokens = l.tokens - 1

	if l.tokens < 0 {
		// "-l.tokens / l.rate" mean how many seconds can obtain these tokens.
		return false, time.Duration(-l.tokens / l.rate * float64(time.Second))
	} else {
		return true, 0
	}
}

// Wait until consume a token.
func (l *Limiter) Wait() {
	isAllowed, sleep := l.Allow(time.Now())
	if !isAllowed {
		time.Sleep(sleep)
	}
}
