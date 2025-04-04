package middlewares

import (
	"net/http"
	"strings"
)

type CorsMiddleware struct {
	allowHeaders string
	corsOrigin   string
}

func NewCorsMiddleware(corsOrigin string, headers ...string) *CorsMiddleware {
	return &CorsMiddleware{allowHeaders: strings.Join(headers, ", "), corsOrigin: corsOrigin}
}

func (c *CorsMiddleware) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", c.corsOrigin)
		w.Header().Set("Access-Control-Allow-Methods", r.Method)
		w.Header().Set("Access-Control-Allow-Headers", c.allowHeaders)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			w.Write(nil)
		}
		next.ServeHTTP(w, r)
	})
}
