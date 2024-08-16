package ratelimiter

import "time"

type rateLimiter struct {
	rate     int
	buckSize int
	tokens   int
	ticker   *time.Ticker
	allow    chan bool
	done     chan bool
}

func (rl *rateLimiter) allowed() bool {
	return <-rl.allow
}

func (rl *rateLimiter) stop() {
	close(rl.done)
}

func newRateLimiter(rate, buckSize int) *rateLimiter {
	rl := &rateLimiter{
		rate:     rate,
		buckSize: buckSize,
		tokens:   buckSize,
		ticker:   time.NewTicker(time.Second / time.Duration(rate)),
		allow:    make(chan bool),
		done:     make(chan bool),
	}

	go func() {
		for {
			select {
			case <-rl.done:
				return
			case rl.allow <- rl.tokens > 0:
				if rl.tokens > 0 {
					rl.tokens--
				}
			case <-rl.ticker.C:
				if rl.tokens < rl.buckSize {
					rl.tokens++
				}
			}
		}
	}()

	return rl
}

type userRateLimiter struct {
	rateLimiters map[string]*rateLimiter
	rate         int
	buckSize     int
}

func (url *userRateLimiter) Allow(username string) bool {
	rl, exist := url.rateLimiters[username]
	if !exist {
		rl = newRateLimiter(url.rate, url.buckSize)
		url.rateLimiters[username] = rl
	}

	return rl.allowed()
}

func (url *userRateLimiter) StopAll() {
	for _, rl := range url.rateLimiters {
		rl.stop()
	}
}

func NewUserRateLimiter(rate, buckSize int) *userRateLimiter {
	return &userRateLimiter{
		rateLimiters: make(map[string]*rateLimiter),
		rate:         rate,
		buckSize:     buckSize,
	}
}
