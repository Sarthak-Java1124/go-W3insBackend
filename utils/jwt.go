package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string
	jwt.RegisteredClaims
}

var jwtKey = []byte("myNameIsSarthakHarsh")

func GenerateJwt(userID string) (string, error) {
	expiryTime := time.Now().Add(time.Hour * 24 * 10)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(jwtKey)
}

func VerifyJWT(tokenstring string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenstring, claims, func(t *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil

}
