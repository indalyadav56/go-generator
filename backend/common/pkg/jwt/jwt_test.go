package jwt

// import (
// 	"testing"
// 	"time"

// 	"github.com/golang-jwt/jwt/v4"
// 	"github.com/stretchr/testify/assert"
// )

// const secretKey = "my_secret_key"

// func TestGenerateToken(t *testing.T) {
// 	jwtHandler := New(secretKey)

// 	claims := map[string]interface{}{
// 		"user_id": "12345",
// 		"role":    "admin",
// 	}

// 	tokenString, err := jwtHandler.GenerateToken(claims, time.Hour*1)
// 	assert.NoError(t, err, "Expected no error when generating token")
// 	assert.NotEmpty(t, tokenString, "Generated token should not be empty")
// }

// func TestValidateToken(t *testing.T) {
// 	jwtHandler := New(secretKey)

// 	claims := map[string]interface{}{
// 		"user_id": "12345",
// 	}
// 	tokenString, err := jwtHandler.GenerateToken(claims, time.Hour*1)
// 	assert.NoError(t, err)

// 	token, err := jwtHandler.ValidateToken(tokenString)
// 	assert.NoError(t, err, "Expected no error when validating token")
// 	assert.True(t, token.Valid, "Token should be valid")
// }

// func TestGetClaims(t *testing.T) {
// 	jwtHandler := New(secretKey)

// 	claims := map[string]interface{}{
// 		"user_id": "12345",
// 		"role":    "admin",
// 	}
// 	tokenString, err := jwtHandler.GenerateToken(claims, time.Hour*1)
// 	assert.NoError(t, err)

// 	token, err := jwtHandler.ValidateToken(tokenString)
// 	assert.NoError(t, err)

// 	extractedClaims, err := jwtHandler.GetClaims(token)
// 	assert.NoError(t, err, "Expected no error when extracting claims")
// 	assert.Equal(t, "12345", extractedClaims["user_id"], "Expected user_id claim to match")
// 	assert.Equal(t, "admin", extractedClaims["role"], "Expected role claim to match")
// }

// func TestInvalidToken(t *testing.T) {
// 	jwtHandler := New(secretKey)

// 	invalidTokenString := "invalid_token_string"
// 	_, err := jwtHandler.ValidateToken(invalidTokenString)
// 	assert.Error(t, err, "Expected an error when validating an invalid token")
// }

// func TestInvalidSigningMethod(t *testing.T) {
// 	jwtHandler := New(secretKey)

// 	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
// 		"user_id": "12345",
// 	})
// 	tokenString, _ := token.SignedString([]byte(secretKey))

// 	_, err := jwtHandler.ValidateToken(tokenString)
// 	assert.Equal(t, ErrInvalidSigningMethod, err, "Expected ErrInvalidSigningMethod error")
// }
