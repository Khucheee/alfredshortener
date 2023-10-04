package app

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

const TOKEN_EXP = 24 * time.Hour
const SECRET_KEY = "supersecretkey"

func parseUserFromCookie(r *http.Request) (string, error) {
	c, _ := r.Cookie("auth")
	uid, err := ParseToken(c.Value)
	if err != nil {
		return "", err
	}
	return uid, nil
}

func MakeToken() (string, error) { //создаем токен
	u := uuid.New()
	jt := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: u.String(),
	})

	tokenString, err := jt.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func ParseToken(tokenString string) (string, error) { //парсим токен

	claims := &Claims{}
	jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	return claims.UserID, nil
}

func makeAuthCookie() (*http.Cookie, error) {
	tokenString, err := MakeToken()
	if err != nil {
		return nil, err
	}
	newCookie := &http.Cookie{
		Name:  "auth",
		Value: tokenString,
	}
	return newCookie, nil
}
