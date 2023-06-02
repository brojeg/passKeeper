package models

import (
	"context"
	"net/http"
	"time"

	server "passKeeper/internal/models/server"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

var (
	ContextUserKey = contextKey("user")
)

// var jwtPassword string
var settings JWTSettings

func InitJWTPassword(pass string, expTime int) {
	settings.jwtPassword = pass
	settings.expirationTime = expTime
}

type JWTSettings struct {
	jwtPassword    string
	expirationTime int
}

type Token struct {
	UserID uint
	jwt.StandardClaims
}

func GetUserFromContext(ctx context.Context) (uint, bool) {
	caller, ok := ctx.Value(ContextUserKey).(uint)
	return caller, ok
}

func GenerateToken(id uint) string {
	expirationTime := time.Now().Add(time.Duration(settings.expirationTime) * time.Minute)
	tk := &Token{UserID: id, StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()}}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte(settings.jwtPassword))
	if err != nil {
		panic(err)
	}

	return tokenString
}

func EncryptPassword(pass string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func IsPasswordsEqual(existing, new string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(existing), []byte(new))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return false
	}
	return true
}

func ValidateToken(r *http.Request) server.Response {
	tokenHeader := r.Header.Get("Authorization")
	expirationTime := time.Now().Add(time.Duration(settings.expirationTime) * time.Minute)
	tk := &Token{StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()}}
	token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(settings.jwtPassword), nil
	})
	if err != nil {
		return server.Response{Message: "Malformed authentication token", ServerCode: 401}
	}
	if !token.Valid {
		return server.Response{Message: "Token is not valid.", ServerCode: 400}
	}
	return server.Response{Message: tk.UserID, ServerCode: 200}
}
