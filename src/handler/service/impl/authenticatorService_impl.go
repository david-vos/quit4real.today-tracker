package impl

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yohcop/openid-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"quit4real.today/config"
	"quit4real.today/src/model"
	"time"
)

type AuthServiceImpl struct {
	OpenID openid.OpenID
}

// NewAuthServiceImpl creates a new instance of AuthServiceImpl.
func NewAuthServiceImpl(openId openid.OpenID) *AuthServiceImpl {
	return &AuthServiceImpl{
		OpenID: openId,
	}
}

func (service *AuthServiceImpl) GetOpenId() openid.OpenID {
	return service.OpenID
}
func (service *AuthServiceImpl) HashPassword(password string) ([]byte, error) {
	// This already does some salting so there is no need to do it later again.
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (service *AuthServiceImpl) CheckPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func (service *AuthServiceImpl) GenerateJWT(user model.User) (string, error) {
	claims := jwt.MapClaims{
		"username":  user.Name,
		"steamName": user.SteamUserName,
		"steamID":   user.SteamID,
		"exp":       time.Now().Add(time.Hour * 1).Unix(), // 1-hour expiration
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JwtSecret())
}

func (service *AuthServiceImpl) GetFieldFromJWT(tokenString string, field string) (string, error) {
	claims := &jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JwtSecret(), nil
	})
	if err != nil {
		return "", fmt.Errorf("could not parse token")
	}
	return (*claims)[field].(string), nil
}

func (service *AuthServiceImpl) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.JwtSecret(), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
