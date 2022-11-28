package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GetToken(id int64, expAt time.Time, signKey string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "user",
		ID:        fmt.Sprintf("%v", id),
		ExpiresAt: jwt.NewNumericDate(expAt),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signKey))
}

func GetRefreshToken(id int64, expAt time.Time, signKey string) (string, error) {
	claims := jwt.RegisteredClaims{
		ID:        fmt.Sprintf("%v", id),
		ExpiresAt: jwt.NewNumericDate(expAt),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signKey))
}

func ParseToken(token, signKey string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		} else if method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("signing method invalid")
		}

		return []byte(signKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		return nil, fmt.Errorf("unexpected error")
	}

	return claims, nil
}
