package server

import (
	"net/http"
)

type ParseFormMiddleware struct {
}

// Middleware is a struct that has a ServeHTTP method
func NewParseForm() *ParseFormMiddleware {
	return &ParseFormMiddleware{}
}

// The middleware handler
func (p *ParseFormMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	req.ParseForm()
	next(w, req)
}
