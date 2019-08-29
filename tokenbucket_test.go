package tokenbucket

import (
	"testing"
	"time"

	"github.com/lifenod/assert"
	"github.com/lifenod/assert/errorassert"
)

func TestLimiter_Allow(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		limiter *Limiter
		now     time.Time

		expectedLimiter *Limiter
		expectedAllowed bool
		expectedSleep   time.Duration
	}{
		{
			name: "has one token",
			limiter: &Limiter{
				rate:   1,
				burst:  1,
				tokens: 1,
				last:   now,
			},
			now: now,

			expectedLimiter: &Limiter{
				rate:   1,
				burst:  1,
				tokens: 0,
				last:   now,
			},
			expectedAllowed: true,
			expectedSleep:   0,
		},

		{
			name: "no token and enough time passed",
			limiter: &Limiter{
				rate:   1,
				burst:  1,
				tokens: 0,
				last:   now,
			},
			now: now.Add(time.Second),

			expectedLimiter: &Limiter{
				rate:   1,
				burst:  1,
				tokens: 0,
				last:   now.Add(time.Second),
			},
			expectedAllowed: true,
			expectedSleep:   0,
		},

		{
			name: "no token and some time passed, rate = 1",
			limiter: &Limiter{
				rate:   1,
				burst:  1,
				tokens: 0,
				last:   now,
			},
			now: now.Add(time.Millisecond * 200),

			expectedLimiter: &Limiter{
				rate:   1,
				burst:  1,
				tokens: -0.8,
				last:   now.Add(time.Millisecond * 200),
			},
			expectedAllowed: false,
			expectedSleep:   time.Millisecond * 800,
		},

		{
			name: "no token and some time passed, rate = 2",
			limiter: &Limiter{
				rate:   2,
				burst:  1,
				tokens: 0,
				last:   now,
			},
			now: now.Add(time.Millisecond * 200),

			expectedLimiter: &Limiter{
				rate:   2,
				burst:  1,
				tokens: -0.6,
				last:   now.Add(time.Millisecond * 200),
			},
			expectedAllowed: false,
			expectedSleep:   time.Millisecond * 300,
		},

		{
			name: "rate reach burst when enough time passed",
			limiter: &Limiter{
				rate:   100,
				burst:  1000,
				tokens: 0,
				last:   now,
			},
			now: now.Add(time.Second * 10),

			expectedLimiter: &Limiter{
				rate:   100,
				burst:  1000,
				tokens: 999,
				last:   now.Add(time.Second * 10),
			},
			expectedAllowed: true,
			expectedSleep:   0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, sleep := test.limiter.Allow(test.now)
			errorassert.Equal(t, test.expectedLimiter, test.limiter)
			errorassert.Equal(t, test.expectedAllowed, b)
			errorassert.Equal(t, test.expectedSleep, sleep)
		})
	}
}

func TestLimiter_Wait__Slept(t *testing.T) {
	lim := NewLimiter(10, 1)
	isWaitDone := false

	go func() {
		lim.Wait()
		lim.Wait()
		isWaitDone = true
	}()

	time.Sleep(time.Millisecond * 50)

	assert.Equal(t, isWaitDone, false)
}

func TestLimiter_Wait__NoSleep(t *testing.T) {
	lim := NewLimiter(40, 1)
	isWaitDone := false

	go func() {
		lim.Wait()
		lim.Wait()
		isWaitDone = true
	}()

	time.Sleep(time.Millisecond * 50)

	assert.Equal(t, isWaitDone, true)
}
