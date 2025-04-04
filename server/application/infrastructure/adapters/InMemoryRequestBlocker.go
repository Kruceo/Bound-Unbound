package adapters

import (
	"context"
	"sync"
	"time"
)

type blockedRequester struct {
	LastTry   time.Time
	LimitTime time.Time
}

type InMemoryRequestBlocker struct {
	storage map[string]*blockedRequester
	mu      sync.Mutex
}

func NewInMemoryBlocker() *InMemoryRequestBlocker {
	b := &InMemoryRequestBlocker{
		storage: make(map[string]*blockedRequester),
	}
	go b.RunCleanupService(context.Background())
	return b
}

func (b *InMemoryRequestBlocker) IsBlocked(ip string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	req, exists := b.storage[ip]
	if !exists {
		return false
	}
	now := time.Now()
	result := now.Before(req.LimitTime)
	if result {
		if now.UnixMilli()-b.storage[ip].LimitTime.UnixMilli() > 20 {
			b.storage[ip].LimitTime = b.storage[ip].LimitTime.Add(3 * time.Second)
			return true
		}
	}
	return false
}

func (b *InMemoryRequestBlocker) MarkAttempt(ip string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	b.storage[ip] = &blockedRequester{
		LastTry:   now,
		LimitTime: now.Add(2 * time.Second),
	}
}

func (b *InMemoryRequestBlocker) Cleanup() {
	b.mu.Lock()
	now := time.Now()
	for k, v := range b.storage {
		if now.Sub(v.LastTry) > time.Minute {
			delete(b.storage, k)
		}
	}
	b.mu.Unlock()
}

func (b *InMemoryRequestBlocker) RunCleanupService(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			b.Cleanup()
		}
	}
}
