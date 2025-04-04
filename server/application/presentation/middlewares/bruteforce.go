package middlewares

import (
	"fmt"
	"net/http"
	"server2/application/infrastructure"
	"server2/application/infrastructure/adapters"
	"server2/application/presentation"
)

type BruteForceMiddleware struct {
	fastErrorResponses presentation.FastErrorResponses
	blocker            infrastructure.RequestBlocker
}

func NewBruteForceMiddleware() *BruteForceMiddleware {
	b := &BruteForceMiddleware{
		fastErrorResponses: presentation.NewFastErrorResponses(),
		blocker:            adapters.NewInMemoryBlocker(),
	}
	return b
}

func (b *BruteForceMiddleware) BruteForceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if b.blocker.IsBlocked(r.RemoteAddr) {
			b.fastErrorResponses.Execute(w, r, "AUTH_BLOCKED", http.StatusUnauthorized)
			fmt.Println("blocked", r.RemoteAddr)
			return
		}

		b.blocker.MarkAttempt(r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
