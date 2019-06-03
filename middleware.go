package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"sku-manage/jwt"
	"sku-manage/mixin"
)

type JWT interface {
	PublicJWT() *jwt.JWT
}

type ServerJWT struct {
	_pubJWT *jwt.JWT
}

func PubJWT(serviceName string) *ServerJWT {
	return &ServerJWT{
		_pubJWT: jwt.NewJwt(serviceName, "1234567890abcdef", time.Duration(31536000)*time.Second),
	}
}

func (this *ServerJWT) PublicJWT() *jwt.JWT {
	return this._pubJWT
}

type TokenMiddleware struct {
	*mixin.ResponseMixin
	_jwt JWT
}

// Middleware is a struct that has a ServeHTTP method
func NewTokenMiddleware(jwt JWT) *TokenMiddleware {
	return &TokenMiddleware{
		ResponseMixin: mixin.NewResponseMixin(),
		_jwt:          jwt,
	}
}

// The middleware handler
func (this *TokenMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	token := req.Header.Get("X-Inspect-Token")
	if token == "" {
		logrus.Error("[TokenMiddleware.ServeHTTP] X-Inspect-Token is nil")
		this.ResponseUnauthorized(w)
		return
	}

	id, subject, issuedAt, ok := this._jwt.PublicJWT().Decode(token)
	if !ok || id != mux.Vars(req)["user_id"] {
		logrus.Errorf("[TokenMiddleware.ServeHTTP] token check error => ok: %v, id: %s, subject: %s %d", ok, id, subject, issuedAt)
		this.ResponseUnauthorized(w)
		return
	}

	next(w, req)
}
