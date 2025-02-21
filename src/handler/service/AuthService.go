package service

import (
	"github.com/yohcop/openid-go"
	"net/http"
	"quit4real.today/src/model"
)

type AuthService interface {
	GetOpenId() openid.OpenID
	HashPassword(password string) ([]byte, error)
	CheckPassword(hashedPassword, password string) bool
	GenerateJWT(user model.User) (string, error)
	GetFieldFromJWT(tokenString string, field string) (string, error)
	AuthMiddleware(next http.HandlerFunc) http.HandlerFunc
}
