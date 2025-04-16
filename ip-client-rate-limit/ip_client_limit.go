package ipclientratelimit

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Client struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
}

type ClientStore struct {
	Clients map[string]*Client
	mu      sync.Mutex
}

const (
	requestsPerSecond = 1
	burstLimit        = 3
	clientTTL         = 5 * time.Minute
)

var store = &ClientStore{
	Clients: make(map[string]*Client),
}

func init() {
	go store.cleanupExpiredClients()
}

func (cs *ClientStore) getLimiter(ip string) *rate.Limiter {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	client, exists := cs.Clients[ip]
	if !exists {
		limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), burstLimit)
		cs.Clients[ip] = &Client{
			Limiter:  limiter,
			LastSeen: time.Now(),
		}
		return limiter
	}

	client.LastSeen = time.Now()
	return client.Limiter
}

func (cs *ClientStore) cleanupExpiredClients() {
	for {
		time.Sleep(1 * time.Minute)
		cs.mu.Lock()
		for ip, client := range cs.Clients {
			if time.Since(client.LastSeen) > clientTTL {
				delete(cs.Clients, ip)
			}
		}
		cs.mu.Unlock()
	}
}

func RateLimiter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		limiter := store.getLimiter(ip)

		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Try again later.",
			})
			return
		}
		ctx.Next()
	}
}
