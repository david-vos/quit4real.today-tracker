package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/url"
	"quit4real.today/config"
	"time"
)

type AuthService struct {
}

func (service *AuthService) HashPassword(password string) ([]byte, error) {
	// This already does some salting so there is no need to do it later again.
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (service *AuthService) CheckPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func (service *AuthService) GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(), // 1-hour expiration
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JwtSecret())
}

func (service *AuthService) GetFieldFromJWT(tokenString string, field string) (string, error) {
	claims := &jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JwtSecret(), nil
	})
	if err != nil {
		return "", fmt.Errorf("could not parse token")
	}
	return (*claims)[field].(string), nil
}

func (service *AuthService) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
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

func (service *AuthService) SteamLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Define the callback URL where Steam will redirect after login
	callbackURL := config.BackendUrl() + "/api/auth/steam/callback"

	// Construct the OpenID authentication URL manually
	params := url.Values{}
	params.Set("openid.ns", "http://specs.openid.net/auth/2.0")
	params.Set("openid.mode", "checkid_setup")
	params.Set("openid.return_to", callbackURL)
	params.Set("openid.realm", config.FrontendUrl())
	params.Set("openid.identity", "http://specs.openid.net/auth/2.0/identifier_select")
	params.Set("openid.claimed_id", "http://specs.openid.net/auth/2.0/identifier_select")

	// Construct the final redirect URL
	redirectURL := "https://steamcommunity.com/openid/login?" + params.Encode()

	// Redirect the user to Steam
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
