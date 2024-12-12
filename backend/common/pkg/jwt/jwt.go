package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidTokenOrClaims = errors.New("invalid token or claims")
)

type JWT interface {
	GenerateToken(claims map[string]interface{}) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GetClaims(token *jwt.Token) (map[string]interface{}, error)
}

type jwtHandler struct {
	config JWTConfig
}

func New(config JWTConfig) JWT {
	return &jwtHandler{
		config: config,
	}
}

// GenerateToken generates a new JWT token with customizable claims
func (j *jwtHandler) GenerateToken(claims map[string]interface{}) (string, error) {
	jwtClaims := jwt.MapClaims(claims)

	jwtClaims["exp"] = time.Now().Add(j.config.TokenDuration).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString(j.config.SecretKey)
}

// ValidateToken validates the given token string
func (j *jwtHandler) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return j.config.SecretKey, nil
	})
}

// GetClaims extracts the claims from a valid token
func (j *jwtHandler) GetClaims(token *jwt.Token) (map[string]interface{}, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidTokenOrClaims
	}
	return claims, nil
}
