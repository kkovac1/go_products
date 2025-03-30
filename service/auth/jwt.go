package auth

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kkovac1/products/config"
)

func GenerateJWT(secret []byte, userId int) (string, error) {
	// Generate JWT token here
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    strconv.Itoa(userId),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
