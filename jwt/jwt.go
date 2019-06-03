package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const MagicString string = "goweb-restful-api-3344"

type JWT struct {
	issuer    string
	secret    []byte
	expiresAt time.Duration
}

func NewJwt(issuer string, secret string, expiresAt time.Duration) *JWT {
	return &JWT{
		issuer:    issuer,
		secret:    []byte(secret),
		expiresAt: expiresAt,
	}
}

type GowebClaims struct {
	Magic string `json:"magic"`
	jwt.StandardClaims
}

func (this *JWT) Encode(id, subject string) (string, error) {
	claims := GowebClaims{
		MagicString,
		jwt.StandardClaims{
			Id:        id,
			ExpiresAt: time.Now().Add(this.expiresAt).Unix(),
			Issuer:    this.issuer,
			Subject:   subject,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(this.secret)
}

func (this *JWT) Decode(token string) (id, subject string, issuedAt int64, ok bool) {
	if token, err := jwt.ParseWithClaims(token, &GowebClaims{}, func(token *jwt.Token) (interface{}, error) {
		return this.secret, nil
	}); err == nil {
		if claims, ok := token.Claims.(*GowebClaims); ok && token.Valid && claims.Magic == MagicString {
			return claims.Id, claims.Subject, claims.IssuedAt, true
		}
	}

	return "", "", 0, false
}
