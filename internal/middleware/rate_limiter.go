package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	// различные параметры
	ips   map[string]*rate.Limiter // карта
	mu    sync.Mutex               // мутекс для карты
	r     rate.Limit               // объект лимитера
	burst int                      // поток обновления
}

func NewRateLimiter(r rate.Limit, burst int) *RateLimiter {
	return &RateLimiter{
		ips:   make(map[string]*rate.Limiter),
		mu:    sync.Mutex{},
		r:     r,
		burst: burst,
	}
}

func (r *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		limiter := r.getRateLimiter(ip)
		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}
		ctx.Next()
	}
}

func (r *RateLimiter) getRateLimiter(ip string) *rate.Limiter {
	r.mu.Lock()
	defer r.mu.Unlock()
	limiter, ok := r.ips[ip]
	if !ok {
		rateLimiter := r.newRateLimiter()
		r.ips[ip] = rateLimiter
		return rateLimiter
	}
	return limiter
}

func (r *RateLimiter) newRateLimiter() *rate.Limiter {
	return rate.NewLimiter(r.r, r.burst)
}

/*package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu         sync.Mutex
	buckets    map[string]int
	refillRate time.Duration
	capacity   int
}

func NewRateLimiter(refillRate time.Duration, capacity int) *RateLimiter {
	return &RateLimiter{
		mu:         sync.Mutex{},
		buckets:    make(map[string]int),
		refillRate: refillRate,
		capacity:   capacity,
	}
}

func (r *RateLimiter) Start() chan<- struct{} {
	stop := make(chan struct{})
	ticker := time.NewTicker(r.refillRate)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r.mu.Lock()
				for ip := range r.buckets {
					r.buckets[ip] = r.capacity
				}
				r.mu.Unlock()
			case <-stop:
				return
			}
		}
	}()
	return stop
}

func (r *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		if !r.isAllow(ip) {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too manu requests",
			})
			return
		}
		ctx.Next()
	}
}

func (r *RateLimiter) isAllow(ip string) bool {
	defer r.mu.Unlock()
	r.mu.Lock()
	tokens, ok := r.buckets[ip]
	if !ok {
		r.buckets[ip] = r.capacity
		return true
	}
	if tokens < 1 {
		return false
	}
	tokens--
	r.buckets[ip] = tokens
	return true
}*/
