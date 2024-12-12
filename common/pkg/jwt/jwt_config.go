package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTConfig struct {
	SecretKey        []byte
	SigningAlgorithm jwt.SigningMethod
	TokenDuration    time.Duration
}
