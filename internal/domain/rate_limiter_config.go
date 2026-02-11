package domain

import "time"

type RateLimiterConfig struct {
	Rate  time.Duration
	Burst int
}
