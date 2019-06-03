package server

import (
	"net/http"

	uber "github.com/sunrongya/ratelimit"
)

type UberRatelimit struct {
	limiter uber.Limiter
}

// Middleware is a struct that has a ServeHTTP method
func NewUberRatelimit(rate int, opts ...uber.Option) *UberRatelimit {
	return &UberRatelimit{
		limiter: uber.New(rate, opts...),
	}
}

// The middleware handler
func (u *UberRatelimit) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	u.limiter.Take()
	next(w, req)
}
